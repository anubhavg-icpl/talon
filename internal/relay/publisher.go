package relay

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishCompletion publishes payload directly to the agent.pentest.output
// queue against the default exchange (routing key == queue name). That
// queue name is a fixed external integration contract -- whatever
// downstream platform consumes completion events listens on it by name.
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
