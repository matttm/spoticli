package services

/// TODO: use go:generate to make singleton

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type ConfigService struct {
	Config aws.Config
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
			configService.Config = cfg
		}
	}
	println("ConfigService Instantiated")
	return configService
}
