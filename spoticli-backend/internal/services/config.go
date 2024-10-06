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
		storageLock.Lock()
		defer storageLock.Unlock()
		if configService == nil {
			configService = &ConfigService{}
			cfg, err := config.LoadDefaultConfig(context.TODO(),
				config.WithSharedCredentialsFiles(
					[]string{"test/credentials", "data/credentials"},
				),
				config.WithSharedConfigFiles(
					[]string{"test/config", "data/config"},
				),
			)
			if err != nil {
				panic(err)
			}
			configService.Config = cfg
		}
	}

	return configService
}
