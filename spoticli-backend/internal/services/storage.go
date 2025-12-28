package services

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/coder/flog"
)

// A StorageService interacts directly with s3
type StorageService struct {
	client S3ClientApi
	// for unauthenticated users:
	psClient S3PresignClientApi
}

var storageLock = &sync.Mutex{}
var TRACKS_BUCKET_NAME = aws.String("spoticli-tracks")
var MIME_MP3 = aws.String("audio/mp3")

var storageService *StorageService

func GetStorageService() *StorageService {
	if storageService == nil {
		storageLock.Lock()
		defer storageLock.Unlock()
		if storageService == nil {
			storageService = &StorageService{}
			cfg := GetConfigService().CloudConfig
			// Use path-style addressing for LocalStack compatibility
			client := s3.NewFromConfig(cfg, func(o *s3.Options) {
				o.UsePathStyle = true
			})
			storageService.client = client
			flog.Infof("S3 client configured with endpoint: %s", *cfg.BaseEndpoint)
			storageService.psClient = s3.NewPresignClient(client)
			flog.Infof("S3 presign client initialized with endpoint: %s", *client.Options().BaseEndpoint)
			flog.Infof("StorageService Instantiated")
		}
	}
	return storageService
}

func (s *StorageService) PostPresignedUrl(key string) (*string, error) {
	input := &s3.PutObjectInput{
		Bucket: TRACKS_BUCKET_NAME,
		//  ContentType: MIME_MP3,
		Key:         aws.String(key),
		ContentType: aws.String("audio/mp3"),
	}
	req, err := s.psClient.PresignPutObject(
		context.TODO(),
		input,
		func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(60 * int64(time.Second))
		},
	)
	if err != nil {
		return nil, err
	}
	return &req.URL, nil
}

// GetPresignedUrl invokes presigned GetObject cmd
func (s *StorageService) GetPresignedUrl(key string) (string, error) {
	res, err := s.psClient.PresignGetObject(
		context.TODO(),
		&s3.GetObjectInput{
			Bucket: TRACKS_BUCKET_NAME,
			Key:    aws.String(key),
		},
	)
	if err != nil {
		flog.Errorf(err.Error())
	}
	return res.URL, nil
}

// DownloadFile invokes GetObject command with a range if provided
func (s *StorageService) DownloadFile(key string, _range *string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: TRACKS_BUCKET_NAME,
		Key:    aws.String(key),
	}
	if _range != nil {
		input.Range = aws.String(*_range)
	}
	res, err := s.client.GetObject(
		context.TODO(),
		input,
	)
	if err != nil {
		flog.Errorf(err.Error())
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		flog.Errorf(err.Error())
		return nil, err
	}
	body, err = ReadID3v2Header(body)
	if err != nil {
		flog.Errorf(err.Error())
		return nil, err
	}
	return body, nil
}

func (s *StorageService) StreamFile(key string, start, end *int64) ([]byte, int64, error) {
	requestedFrames := make(chan []byte, 1)
	if !isItemCached(key) {
		input := &s3.GetObjectInput{
			Bucket: TRACKS_BUCKET_NAME,
			Key:    aws.String(key),
		}
		res, err := s.client.GetObject(
			context.TODO(),
			input,
		)
		if err != nil {
			flog.Errorf(err.Error())
			return nil, 0, err
		}
		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			flog.Errorf(err.Error())
			return nil, 0, err
		}
		body, err = ReadID3v2Header(body)
		if err != nil {
			flog.Errorf(err.Error())
			return nil, 0, err
		}
		framesBytes := len(body)
		frames := PartitionMp3Frames(body)
		flog.Infof("Frame count: %d", len(frames))
		// end test NOTE:
		// TODO : put in a goroutine
		cacheItem(key, frames, *start, *end, requestedFrames)
		return getSegmentFromCache(key, start, end), int64(framesBytes), nil
	} else {
		return getSegmentFromCache(key, start, end), filesize(key), nil
	}
	// x := <-requestedFrames
	// return x, nil
}
