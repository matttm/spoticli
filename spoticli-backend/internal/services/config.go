package services

/// TODO: use go:generate to make singleton

import (
	"context"
	"os"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

//	A ConfigService manages configuratio
//
// needed for services and environment vars
type ConfigService struct {
	CloudConfig aws.Config
	Config      map[string]interface{}
}

var configLock = &sync.Mutex{}

var configService *ConfigService

func GetConfigService() *ConfigService {
	if configService == nil {
		configLock.Lock()
		defer configLock.Unlock()
		if configService == nil {
			configService = &ConfigService{}
			cfg, err := config.LoadDefaultConfig(context.TODO())
			if err != nil {
				panic(err)
			}
			configService.CloudConfig = cfg
			configService.Config = map[string]interface{}{}
			segSz, _ := strconv.Atoi(os.Getenv("STREAM_SEGMENT_SIZE"))
			configService.Config["STREAM_SEGMENT_SIZE"] = int64(segSz)
			println("ConfigService Instantiated")
		}
	}
	return configService
}

// GetConfigValue gets an config var
func GetConfigValue[T any](k string) T {
	return GetConfigService().Config[k].(T) // type assertion
}
