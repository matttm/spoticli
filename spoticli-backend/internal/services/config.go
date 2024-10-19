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
			configService.Config["STREAM_SEGMENT_SIZE"], _ = strconv.Atoi(os.Getenv("STREAM_SEGMENT_SIZE"))
			println("ConfigService Instantiated")
		}
	}
	return configService
}
func GetConfigValue[T any](k string) T {
	return GetConfigService().Config[k].(T) // type assertion
}
