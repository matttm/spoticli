package services

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// / TODO: use go:generate to make singleton
type StorageService struct {
	client   *s3.Client
	psClient *s3.PresignClient
}

var storageLock = &sync.Mutex{}

var storageService *StorageService

func GetStorageService() *StorageService {
	if storageService == nil {
		storageLock.Lock()
		defer storageLock.Unlock()
		if storageService == nil {
			storageService = &StorageService{}
			storageService.client = s3.NewFromConfig(GetConfigService().Config)
			storageService.psClient = s3.NewPresignClient(storageService.client)
		}
	}

	return storageService
}

func (s *StorageService) GetPresignedUrl() (*v4.PresignedHTTPRequest, error) {
	return s.psClient.PresignGetObject(
		context.TODO(),
		&s3.GetObjectInput{
			Bucket: aws.String("spoticli-tracks"),
			Key:    aws.String("bat_country.mp3"),
		},
	)
}
