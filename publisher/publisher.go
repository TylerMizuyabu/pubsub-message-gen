package publisher

import (
	"context"
	"fmt"
	"pubsub-message-gen/publisher/messages"
	"pubsub-message-gen/types"
	"time"

	"cloud.google.com/go/pubsub"
)

// Publisher struct
type Publisher struct {
	config types.Config
}

// NewPublisher creates a new Publisher struct given a config
func NewPublisher(config *types.Config) *Publisher {
	return &Publisher{
		config: *config,
	}
}

// Start function begins the process of publishing messages using the provided config
func (p *Publisher) Start(client *pubsub.Client) {
	ticker := time.NewTicker(time.Duration(p.config.MessageInterval) * time.Millisecond)
	messageGenerator := messages.SelectMessageGenerator(p.config.MessageType)
	for {
		select {
		case _ = <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.config.MessageInterval)*time.Millisecond-100)
			defer cancel()
			message := messageGenerator()
			result := client.Topic(p.config.Topic).Publish(ctx, &message)
			messageID, err := result.Get(ctx)
			if err != nil {
				fmt.Printf("Error occurred attempting to publish message to topic %s", err.Error())
				continue
			}
			fmt.Printf("Successfully sent message. %s", messageID)
			cancel()
		}
	}
}