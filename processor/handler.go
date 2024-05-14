package processor

import (
	documentai "cloud.google.com/go/documentai/apiv1"
	"cloud.google.com/go/documentai/apiv1/documentaipb"
	"context"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/event"
	"google.golang.org/api/option"
	"os"
)

func HandleEvent(ctx context.Context, e event.Event) error {
	uploadEvent, err := unmarshalEvent(e)
	if err != nil {
		return err
	}
	location := os.Getenv("PROCESSOR_LOCATION")
	processorId := os.Getenv("PROCESSOR_ID")
	endpoint := fmt.Sprintf("%s-documentai.googleapis.com:443", location)
	fmt.Println(endpoint)
	client, err := documentai.NewDocumentProcessorClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		fmt.Println(fmt.Errorf("error creating Document AI client: %w", err))
	}
	defer client.Close()

	req := &documentaipb.ProcessRequest{
		Name: processorId,
		Source: &documentaipb.ProcessRequest_GcsDocument{
			GcsDocument: &documentaipb.GcsDocument{
				GcsUri:   createGcsUri(uploadEvent.Bucket, uploadEvent.Name),
				MimeType: uploadEvent.ContentType,
			},
		},
	}
	resp, err := client.ProcessDocument(ctx, req)
	if err != nil {
		fmt.Println(fmt.Errorf("processDocument: %w", err))
	}

	// Handle the results.
	document := resp.GetDocument()
	fmt.Printf("Document Text: %s", document.GetText())
	return nil
}

func unmarshalEvent(e event.Event) (StorageUploadEvent, error) {
	var uploadEvent StorageUploadEvent
	err := e.DataAs(&uploadEvent)
	if err != nil {
		return StorageUploadEvent{}, err
	}
	return uploadEvent, nil
}

func createGcsUri(bucket string, object string) string {
	return fmt.Sprintf("gs://%s/%s", bucket, object)
}
