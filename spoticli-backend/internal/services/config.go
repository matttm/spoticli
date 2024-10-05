package services

/// TODO: use go:generate to make singleton

import "sync"

type ConfigService struct {
}

var configLock = &sync.Mutex{}

var configService *ConfigService

func GetConfigService() *ConfigService {
	if configService == nil {
		storageLock.Lock()
		defer storageLock.Unlock()
		if configService == nil {
			configService = &ConfigService{}
		}
	}

	return configService
}
