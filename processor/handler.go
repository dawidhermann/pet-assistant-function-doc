package processor

import (
	documentai "cloud.google.com/go/documentai/apiv1"
	"cloud.google.com/go/documentai/apiv1/documentaipb"
	"context"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/event"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
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
	fmt.Println(os.Environ())
	client, err := documentai.NewDocumentProcessorClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		fmt.Println(fmt.Errorf("error creating Document AI client: %w", err))
	}
	defer client.Close()

	req := &documentaipb.BatchProcessRequest{
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
		//OCR_RESULTS_BUCKET
		//Source: &documentaipb.ProcessRequest_GcsDocument{
		//	GcsDocument: &documentaipb.GcsDocument{
		//		GcsUri:   createGcsUri(uploadEvent.Bucket, uploadEvent.Name),
		//		MimeType: uploadEvent.ContentType,
		//	},
		//},
	}
	resp, err := client.BatchProcessDocuments(ctx, req)
	if err != nil {
		fmt.Println(fmt.Errorf("processDocument: %w", err))
	}

	// Handle the results.
	fmt.Println(fmt.Sprintf("gs://%s", os.Getenv("OCR_RESULTS_BUCKET")))
	fmt.Println(resp.Name())
	fmt.Println(resp.Metadata())
	fmt.Println(resp.Done())
	//document := resp.Done()
	//fmt.Printf("Document Text: %s", document.GetText())
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
