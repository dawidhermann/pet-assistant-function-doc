package processor

import (
	"time"
)

type StorageUploadEvent struct {
	Kind                    string    `json:"kind"`
	Id                      string    `json:"id"`
	SelfLink                string    `json:"selfLink"`
	Name                    string    `json:"name"`
	Bucket                  string    `json:"bucket"`
	Generation              string    `json:"generation"`
	Metageneration          string    `json:"metageneration"`
	ContentType             string    `json:"contentType"`
	TimeCreated             string    `json:"timeCreated"`
	Updated                 time.Time `json:"updated"`
	StorageClass            string    `json:"storageClass"`
	TimeStorageClassUpdated time.Time `json:"timeStorageClassUpdated"` // `2024-05-12T18:58:49.662Z`
	Size                    string    `json:"size"`
	Md5Hash                 string    `json:"md5Hash"`
	MediaLink               string    `json:"mediaLink"`
	Etag                    string    `json:"etag"`
}
