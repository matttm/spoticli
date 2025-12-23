package services

import (
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	mock_services "github.com/matttm/spoticli/spoticli-backend/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestStorageService_GetPresignedUrl_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPresignClient := mock_services.NewMockS3PresignClientApi(ctrl)
	storageService := &StorageService{
		psClient: mockPresignClient,
	}

	expectedURL := "https://s3.amazonaws.com/spoticli-tracks/test-key?presigned=true"
	key := "test-key"

	mockPresignClient.EXPECT().
		PresignGetObject(
			gomock.Any(),
			gomock.Any(),
		).
		Return(&v4.PresignedHTTPRequest{URL: expectedURL}, nil).
		Times(1)

	url, err := storageService.GetPresignedUrl(key)

	assert.NoError(t, err)
	assert.Equal(t, expectedURL, url)
}

func TestStorageService_PostPresignedUrl_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPresignClient := mock_services.NewMockS3PresignClientApi(ctrl)
	storageService := &StorageService{
		psClient: mockPresignClient,
	}

	expectedURL := "https://s3.amazonaws.com/spoticli-tracks/test-key?presigned=true"
	key := "test-key"

	mockPresignClient.EXPECT().
		PresignPutObject(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Return(&v4.PresignedHTTPRequest{URL: expectedURL}, nil).
		Times(1)

	url, err := storageService.PostPresignedUrl(key)

	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, expectedURL, *url)
}

func TestStorageService_DownloadFile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_services.NewMockS3ClientApi(ctrl)
	storageService := &StorageService{
		client: mockClient,
	}

	key := "test-key"
	// Create a valid ID3v2 header: "ID3" + version + flags + size (sync-safe)
	audioData := []byte("fake audio data")
	// Size in sync-safe format: 15 bytes = 0x0F = 0b00001111
	// Sync-safe means using only 7 bits per byte, so: 0x00 0x00 0x00 0x0F
	id3Header := []byte{
		0x49, 0x44, 0x33, // "ID3"
		0x03, 0x00, // version 2.3.0
		0x00,                   // flags
		0x00, 0x00, 0x00, 0x0F, // size: 15 bytes (sync-safe)
	}
	// ID3 tag content (15 bytes to match the size)
	id3TagContent := make([]byte, 15)
	// Actual audio data comes after the ID3 tag
	fullData := append(id3Header, id3TagContent...)
	fullData = append(fullData, audioData...)
	mockBody := io.NopCloser(strings.NewReader(string(fullData)))

	mockClient.EXPECT().
		GetObject(
			gomock.Any(),
			gomock.Any(),
		).
		Return(&s3.GetObjectOutput{
			Body: mockBody,
		}, nil).
		Times(1)

	body, err := storageService.DownloadFile(key, nil)

	assert.NoError(t, err)
	assert.NotNil(t, body)
	// The body should have the ID3 header stripped
	assert.Equal(t, audioData, body)
}

func TestStorageService_DownloadFile_WithRange_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_services.NewMockS3ClientApi(ctrl)
	storageService := &StorageService{
		client: mockClient,
	}

	key := "test-key"
	rangeStr := "bytes=0-100"
	// Create a valid ID3v2 header
	audioData := []byte("fake audio data")
	// Size in sync-safe format: 15 bytes = 0x0F
	id3Header := []byte{
		0x49, 0x44, 0x33, // "ID3"
		0x03, 0x00, // version 2.3.0
		0x00,                   // flags
		0x00, 0x00, 0x00, 0x0F, // size: 15 bytes (sync-safe)
	}
	// ID3 tag content (15 bytes to match the size)
	id3TagContent := make([]byte, 15)
	// Actual audio data comes after the ID3 tag
	fullData := append(id3Header, id3TagContent...)
	fullData = append(fullData, audioData...)
	mockBody := io.NopCloser(strings.NewReader(string(fullData)))

	mockClient.EXPECT().
		GetObject(
			gomock.Any(),
			gomock.Eq(&s3.GetObjectInput{
				Bucket: aws.String("spoticli-tracks"),
				Key:    aws.String(key),
				Range:  aws.String(rangeStr),
			}),
		).
		Return(&s3.GetObjectOutput{
			Body: mockBody,
		}, nil).
		Times(1)

	body, err := storageService.DownloadFile(key, &rangeStr)

	assert.NoError(t, err)
	assert.NotNil(t, body)
	assert.Equal(t, audioData, body)
}
