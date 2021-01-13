package main

import (
	"context"
	"fmt"
	"os"
	"pubsub-message-gen/publisher"
	"pubsub-message-gen/types"
	"sync"

	"cloud.google.com/go/pubsub"
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
				fmt.Printf("Could not create a pubsub client. Double check config\n\tConfigFile: %s\n\tConfig:%v", f, cfg)
				return
			}

			p := publisher.NewPublisher(cfg)

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