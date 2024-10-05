package services

import "sync"

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
