package messages

import (
	"encoding/json"
	"fmt"
	"pubsub-message-gen/types"

	"cloud.google.com/go/pubsub"
	"github.com/brianvoe/gofakeit"
)

const (
	// UsageRecordMessage constant is the string representation for the UsageRecord message type
	UsageRecordMessage string =  "UsageRecord"
)

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

// SelectMessageGenerator function returns the desired message generator based on the message type argument
func SelectMessageGenerator(messageType string) func() pubsub.Message {
	switch(messageType) {
	case UsageRecordMessage:
		return generateUsageRecordMessage
	default:
		return func() pubsub.Message {
			return pubsub.Message{}
		}
	}
}