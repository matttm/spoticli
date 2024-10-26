// Package services
package services

type CacheService struct {
	redis map[int]map[int64][]byte // map of id to song, which is cut in segments
}

var cacheService *CacheService

func getCacheService() *CacheService {
	if cacheService == nil {
		cacheService = &CacheService{}
	}
	return cacheService
}

func isItemCached(id int) bool {
}

func getSegmentFromCache(id, segment int) []byte {
}

func cacheItem(id int, data [][]byte) error {
}
