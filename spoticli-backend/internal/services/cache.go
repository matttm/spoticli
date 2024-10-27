// Package services
package services

import (
	"fmt"

	"github.com/matttm/spoticli/spoticli-backend/internal/utilities"
)

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

func getSegmentFromCache(key string, reqStart int64) []byte {
	// TODO: REIMPLEMENT WITH BIN-SEARCH
	fmt.Printf("merged frame count from in getter: %d\n", len(cacheService.Redis[key]))
	var sum int64 = 0
	for i, v := range cacheService.Redis[key] {
		fmt.Printf("merged frame i size: %d\n", len(cacheService.Redis[key][i]))
		fmt.Printf("start %d -- sum %d\n", reqStart, sum)
		sum += int64(len(v))
		if reqStart <= sum {
			fmt.Printf("Sending frame %d\n", i)
			return v
		}
	}
	panic("Unable to get cache store for key")
	return nil
}

func cacheItem(key string, frames [][]byte, reqStart, reqEnd int64, reqFrames chan []byte) error {
	if isItemCached(key) {
		panic("Item is already cached")
	}
	frameClusterSize := GetConfigValue[int64]("FRAME_CLUSTER_SIZE")
	var startFrame int64 = 0
	var endFrame int64 = 0
	var curByte int64 = 0
	n := int64(len(frames))
	var songSegments [][]byte
	for startFrame < n {
		endFrame = min(startFrame+frameClusterSize, n)
		b := utilities.Flatten(frames[startFrame:endFrame])
		songSegments = append(songSegments, b)
		startFrame += frameClusterSize
		curByte += int64(len(b))
		if curByte >= reqEnd {
		}
	}
	// TODO: optimize reqFrames later
	// reqFrames <- getSegmentFromCache(key, reqStart)
	getCacheService().Redis[key] = songSegments
	fmt.Printf("merged frame count: %d\n", len(cacheService.Redis[key]))
	fmt.Printf("merged frame-0 size: %d\n", len(cacheService.Redis[key][0]))
	return nil
}
