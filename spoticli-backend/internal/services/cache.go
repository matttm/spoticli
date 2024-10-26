// Package services
package services

import "github.com/matttm/spoticli/spoticli-backend/internal/utilities"

type CacheService struct {
	Redis map[string][][]byte // map of id to song, which is cut in segments
}

var cacheService *CacheService

func getCacheService() *CacheService {
	if cacheService == nil {
		cacheService = &CacheService{}
		cacheService.Redis = make(map[string][][]byte)
	}
	return cacheService
}

func isItemCached(key string) bool {
	_, ok := getCacheService().Redis[key]
	return ok
}

func getSegmentFromCache(key string, segment int64) []byte {
	// TODO: REIMPLEMENT WITH BIN-SEARCH
	sum := 0
	for _, v := range getCacheService().Redis[key] {
		sum += len(v)
		if int(segment) > sum {
			return v
		}
	}
	return nil
}

func cacheItem(key string, frames [][]byte, reqStart, reqEnd int64, reqFrames chan []byte) error {
	if isItemCached(key) {
		panic("Item is already cached")
	}
	frameClusterSize := int64(GetConfigValue[int]("FRAME_CLUSTER_SIZE"))
	var cur int64 = 0
	var end int64 = 0
	n := int64(len(frames))
	var songSegments [][]byte
	for cur < n {
		end = min(cur+frameClusterSize, n)
		b := utilities.Flatten(frames[cur:end])
		songSegments = append(songSegments, b)
		cur += frameClusterSize
		if cur >= reqEnd {
			reqFrames <- getSegmentFromCache(key, reqStart)
		}
	}
	getCacheService().Redis[key] = songSegments
	return nil
}
