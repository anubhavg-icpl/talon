package queue

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishCompletion publishes payload directly to the agent.pentest.output
// queue against the default exchange (routing key == queue name), the
// simplest correct equivalent of executor.py's
// celery_app.send_task("get_pentest_ouput", queue="agent.pentest.output")
// -- the "ouput" typo is Celery's task name on the platform side, not
// anything we need to reproduce in a plain AMQP message.
func PublishCompletion(ctx context.Context, ch *amqp.Channel, payload CompletionPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("queue: marshal completion payload: %w", err)
	}
	return ch.PublishWithContext(ctx, "", outputQueue, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}
