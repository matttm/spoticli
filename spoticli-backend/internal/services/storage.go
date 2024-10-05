package services

import "sync"

type StorageService struct {
}

var lock = &sync.Mutex{}

var instance *StorageService

func GetInstance() *StorageService {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = &StorageService{}
		}
	}

	return instance
}
