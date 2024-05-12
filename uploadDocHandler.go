package handler

import (
	"context"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

type DocUploadedEvent struct {
	Bucket string `json:"bucket,omitempty"`
	Key    string `json:"key,omitempty"`
}

const location = "eu"

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("UploadDocHandler", uploadDocHandler)
}

func uploadDocHandler(ctx context.Context, e event.Event) error {
	fmt.Println(e)
	return nil
}
