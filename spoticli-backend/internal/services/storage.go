package services

import "sync"

// / TODO: use go:generate to make singleton
type StorageService struct {
}

var storageLock = &sync.Mutex{}

var storageService *StorageService

func GetStorageService() *StorageService {
	if storageService == nil {
		storageLock.Lock()
		defer storageLock.Unlock()
		if storageService == nil {
			storageService = &StorageService{}
		}
	}

	return storageService
}
