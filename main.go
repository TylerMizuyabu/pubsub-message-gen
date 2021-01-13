package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	types "pubsub-message-gen/types"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/brianvoe/gofakeit"
	"github.com/ilyakaznacheev/cleanenv"
)


func main() {
	files := getFileArgs()
	if len(files) == 0 {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(files))
	for _, file := range files {
		var f string = file
		go func(){
			defer wg.Done()

			cfg := readConfig(f)
			if cfg == nil {
				fmt.Printf("Could not read config file %s", f)
				return
			}
			client := createPubsubClient(context.Background(), cfg.ProjectID)
			if client == nil {
				fmt.Printf("Could not create a pubsub client. Double check config\n\tConfigFile: %s\n\tConfig:%s", f, cfg)
				return
			}

			p := &Publisher{config: *cfg}

			p.Start(client)
		}()
	}

	wg.Wait()
}

func readConfig(file string) *types.Config {
	var cfg types.Config
	err := cleanenv.ReadConfig(file, &cfg)
	if err != nil {
		return nil
	}
	return &cfg
}

func createPubsubClient(ctx context.Context, projectID string) *pubsub.Client {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil
	}
	return client
}

func getFileArgs() []string {
	return os.Args[1:]
}

/////

//Publisher struct
type Publisher struct {
	config types.Config
}

func (p *Publisher) Start(client *pubsub.Client) {
	ticker := time.NewTicker(time.Duration(p.config.MessageInterval) * time.Millisecond)
	messageGenerator := selectMessageGenerator(p.config.MessageType)
	for {
		select {
		case _ = <- ticker.C:
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

func selectMessageGenerator(messageType string) func() pubsub.Message {
	switch(messageType) {
	case "UsageRecord":
		return generateUsageRecordMessage
	default:
		return func() pubsub.Message {
			return pubsub.Message{}
		}
	}
}

////
var generateUsageRecordMessage func() pubsub.Message = func() pubsub.Message {
	var body types.UsageRecord
	gofakeit.Struct(&body)
	data, err := json.Marshal(body)
	if err != nil {
		// Doing this because I don't want to handle an error case
		fmt.Println("Failed to generate usage record")
		data = make([]byte, 0)
	}
	return pubsub.Message{
		Data: data,
	}
}