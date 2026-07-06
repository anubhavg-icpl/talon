// Package relay is a plain AMQP consumer/publisher: it consumes pentest
// run requests off a task queue and publishes completion events, using a
// flat JSON envelope rather than any particular task-queue framework's
// wire protocol.
package relay

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/anubhavg-icpl/pentester2/internal/config"
	"github.com/anubhavg-icpl/pentester2/internal/core"
)

const (
	taskQueue   = "execute_agent_task"
	outputQueue = "agent.pentest.output"

	// maxAutoApprovals bounds the auto-approve loop below in case the
	// orchestrator keeps re-raising the same interrupt.
	maxAutoApprovals = 50
)

// AgentInputs is the target/attacker-context payload for one pentest run.
type AgentInputs struct {
	TargetIP    string `json:"target_ip"`
	CVEID       string `json:"cve_id"`
	LHOST       string `json:"lhost"`
	LPORT       int    `json:"lport"`
	Description string `json:"description"`
}

// AgentTask is the run_id/project_id/agent_name/agent_inputs message
// consumed off taskQueue.
type AgentTask struct {
	RunID       string      `json:"run_id"`
	ProjectID   string      `json:"project_id"`
	AgentName   string      `json:"agent_name"`
	AgentInputs AgentInputs `json:"agent_inputs"`
}

// CompletionResult is the "result" object in CompletionPayload.
type CompletionResult struct {
	Summary   string `json:"summary"`
	RawOutput string `json:"raw_output"`
}

// CompletionPayload is the run-completion event published to outputQueue.
type CompletionPayload struct {
	EventType     string           `json:"event_type"`
	RunID         string           `json:"run_id"`
	ProjectID     string           `json:"project_id"`
	AgentName     string           `json:"agent_name"`
	OverallStatus string           `json:"overall_status"`
	Result        CompletionResult `json:"result"`
	Timestamp     string           `json:"timestamp"`
}

// Worker consumes execute_agent_task messages and publishes RUN_COMPLETED
// events to agent.pentest.output.
type Worker struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// NewWorker dials the broker at url and declares both queues used by the
// worker. Fails fast (no guest:guest@localhost fallback -- see
// config.LoadAMQPConfig) if url is empty.
func NewWorker(url string) (*Worker, error) {
	if url == "" {
		return nil, fmt.Errorf("queue: AMQP url is empty (set AMQP_URL)")
	}
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("queue: dial: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("queue: open channel: %w", err)
	}
	if _, err := ch.QueueDeclare(taskQueue, true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("queue: declare %s: %w", taskQueue, err)
	}
	if _, err := ch.QueueDeclare(outputQueue, true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("queue: declare %s: %w", outputQueue, err)
	}
	return &Worker{conn: conn, ch: ch}, nil
}

// Close shuts down the channel and connection.
func (w *Worker) Close() error {
	chErr := w.ch.Close()
	connErr := w.conn.Close()
	if chErr != nil {
		return chErr
	}
	return connErr
}

// Consume runs the receive loop until ctx is cancelled or the delivery
// channel closes, mirroring execute_pentest_task's body per message.
func (w *Worker) Consume(ctx context.Context, orchestrator *core.Orchestrator) error {
	deliveries, err := w.ch.Consume(taskQueue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("queue: consume %s: %w", taskQueue, err)
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case d, ok := <-deliveries:
			if !ok {
				return fmt.Errorf("queue: delivery channel closed")
			}
			w.handle(ctx, orchestrator, d)
		}
	}
}

func (w *Worker) handle(ctx context.Context, orchestrator *core.Orchestrator, d amqp.Delivery) {
	var task AgentTask
	if err := json.Unmarshal(d.Body, &task); err != nil {
		log.Printf("[PentestWorker] bad message, dropping: %v", err)
		d.Nack(false, false)
		return
	}
	log.Printf("[PentestWorker] Received execution request for run %s", task.RunID)

	result, err := runToCompletion(ctx, orchestrator, task)
	if err != nil {
		log.Printf("[PentestWorker] Execution failed for run %s: %v", task.RunID, err)
		d.Nack(false, false)
		return
	}

	if task.RunID != "" && task.ProjectID != "" {
		payload := CompletionPayload{
			EventType:     "RUN_COMPLETED",
			RunID:         task.RunID,
			ProjectID:     task.ProjectID,
			AgentName:     task.AgentName,
			OverallStatus: "completed",
			Result: CompletionResult{
				Summary:   result.FinalMessage,
				RawOutput: rawOutput(result),
			},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		if err := PublishCompletion(ctx, w.ch, payload); err != nil {
			log.Printf("[PentestWorker] Failed to report results for run %s: %v", task.RunID, err)
		} else {
			log.Printf("[PentestWorker] Reported completion for run %s", task.RunID)
		}
	}

	d.Ack(false)
}

// runToCompletion drives the orchestrator to a final result, auto-approving
// any HITL interrupt (nmap_scan) since this worker path has no human
// attached -- unlike the HTTP API path, which surfaces PendingInterrupt for
// a real decision. Deliberate simplification: revisit if a queue-based
// approval protocol is ever needed here too.
func runToCompletion(ctx context.Context, orchestrator *core.Orchestrator, task AgentTask) (core.RunResult, error) {
	input := core.RunInput{
		TargetIP:    task.AgentInputs.TargetIP,
		CVEID:       task.AgentInputs.CVEID,
		Description: task.AgentInputs.Description,
		Context: config.Context{
			LHOST: defaultString(task.AgentInputs.LHOST, "0.0.0.0"),
			LPORT: defaultInt(task.AgentInputs.LPORT, 4444),
		},
	}

	result, err := orchestrator.Run(ctx, input)
	for i := 0; err == nil && result.Interrupted && i < maxAutoApprovals; i++ {
		result, err = orchestrator.Resume(ctx, input, core.Decision{Type: "approve"})
	}
	return result, err
}

func rawOutput(result core.RunResult) string {
	b, err := json.Marshal(result)
	if err != nil {
		return fmt.Sprintf("%+v", result)
	}
	return string(b)
}

func defaultString(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

func defaultInt(v, fallback int) int {
	if v == 0 {
		return fallback
	}
	return v
}
