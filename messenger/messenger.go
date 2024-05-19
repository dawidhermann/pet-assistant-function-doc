package messenger

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
)

func WriteMessage(ctx context.Context, projectId string, topic string, message []byte) error {
	client, err := pubsub.NewClient(ctx, projectId)
	fmt.Println("project", projectId)
	fmt.Println("topic", topic)
	if err != nil {
		return fmt.Errorf("failed to create pubsub client: %w", err)
	}
	defer client.Close()
	t := client.Topic(topic)
	result := t.Publish(ctx, &pubsub.Message{
		Data: message,
	})
	_, err = result.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to send a message to topic %s: %w", topic, err)
	}
	fmt.Println("Sent!")
	return nil
}
