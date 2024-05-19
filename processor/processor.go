package processor

import (
	documentai "cloud.google.com/go/documentai/apiv1"
	"cloud.google.com/go/documentai/apiv1/documentaipb"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/dawidhermann/pet-assistant-function-doc/messenger"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"os"
	"strings"
)

func HandleEvent(ctx context.Context, e event.Event) error {
	uploadEvent, err := unmarshalEvent(e)
	if err != nil {
		return err
	}
	location := os.Getenv("PROCESSOR_LOCATION")
	processorId := os.Getenv("PROCESSOR_ID")
	projectId := os.Getenv("PROJECT_ID")
	statusTopic := os.Getenv("PUBSUB_STATUS_TOPIC")
	endpoint := fmt.Sprintf("%s-documentai.googleapis.com:443", location)
	client, err := documentai.NewDocumentProcessorClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("error creating Document AI client: %w", err)
	}
	defer client.Close()
	req := createBatchProcessRequest(processorId, uploadEvent)
	res, err := client.BatchProcessDocuments(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to trigger batch document processing: %w", err)
	}
	batchId, err := getBatchId(res.Name())
	if err != nil {
		return err
	}
	message := processingTriggeredEvent{
		Name:        uploadEvent.Name,
		Bucket:      uploadEvent.Bucket,
		BatchId:     batchId,
		FullBatchId: res.Name(),
	}
	messageStr, err := json.Marshal(message)
	fmt.Printf("BatchId: %s", batchId)
	fmt.Printf("Message: %s", string(messageStr))
	if err != nil {
		return fmt.Errorf("failed to convert message to json format: %w", err)
	}
	return messenger.WriteMessage(ctx, projectId, statusTopic, messageStr)
}

func unmarshalEvent(e event.Event) (StorageUploadEvent, error) {
	var uploadEvent StorageUploadEvent
	err := e.DataAs(&uploadEvent)
	if err != nil {
		return StorageUploadEvent{}, err
	}
	return uploadEvent, nil
}

func createBatchProcessRequest(processorId string, uploadEvent StorageUploadEvent) *documentaipb.BatchProcessRequest {
	return &documentaipb.BatchProcessRequest{
		Name: processorId,
		InputDocuments: &documentaipb.BatchDocumentsInputConfig{
			Source: &documentaipb.BatchDocumentsInputConfig_GcsDocuments{
				GcsDocuments: &documentaipb.GcsDocuments{
					Documents: []*documentaipb.GcsDocument{
						{
							GcsUri:   createGcsUri(uploadEvent.Bucket, uploadEvent.Name),
							MimeType: uploadEvent.ContentType,
						},
					},
				},
			},
		},
		DocumentOutputConfig: &documentaipb.DocumentOutputConfig{
			Destination: &documentaipb.DocumentOutputConfig_GcsOutputConfig_{
				GcsOutputConfig: &documentaipb.DocumentOutputConfig_GcsOutputConfig{
					GcsUri: fmt.Sprintf("gs://%s", os.Getenv("OCR_RESULTS_BUCKET")),
					FieldMask: &fieldmaskpb.FieldMask{
						Paths: []string{"text"},
					},
				},
			},
		},
		ProcessOptions: &documentaipb.ProcessOptions{
			OcrConfig: &documentaipb.OcrConfig{
				Hints: &documentaipb.OcrConfig_Hints{
					LanguageHints: []string{"pl-PL"},
				},
				DisableCharacterBoxesDetection: true,
			},
		},
	}
}

func createGcsUri(bucket string, object string) string {
	return fmt.Sprintf("gs://%s/%s", bucket, object)
}

func getBatchId(batchOperationName string) (string, error) {
	res := strings.Split(batchOperationName, "/")
	resLen := len(res)
	if resLen == 0 {
		return "", fmt.Errorf("failed to get batch operation id from: %s", batchOperationName)
	}
	return res[resLen-1], nil
}
