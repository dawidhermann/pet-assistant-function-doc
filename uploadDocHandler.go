package handler

import (
	"context"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/dawidhermann/pet-assistant-function-doc/processor"
)

type DocUploadedEvent struct {
	Bucket string `json:"bucket,omitempty"`
	Key    string `json:"key,omitempty"`
}

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("UploadDocHandler", uploadDocHandler)
}

func uploadDocHandler(ctx context.Context, e event.Event) error {
	uploadEvent, err := processor.UnmarshalEvent(e)
	if err != nil {
		return err
	}
	return processor.HandleEvent(ctx, uploadEvent)
}
