// Package services
package services

import (
	"github.com/coder/flog"

	"github.com/matttm/spoticli/spoticli-backend/internal/utilities"
)

type CacheService struct {
	Redis         map[string][][]byte // map of id to song, which is cut in segments
	configService ConfigServiceApi
}

var cacheService *CacheService

func getCacheService() *CacheService {
	if cacheService == nil {
		cacheService = &CacheService{}
		cacheService.Redis = make(map[string][][]byte)
		cacheService.configService = GetConfigService()
	}
	return cacheService
}

func isItemCached(key string) bool {
	_, ok := getCacheService().Redis[key]
	return ok
}

func getSegmentFromCache(key string, reqStart, reqEnd *int64) []byte {
	// TODO: REIMPLEMENT WITH BIN-SEARCH
	var sum int64 = 0
	for i, v := range cacheService.Redis[key] {
		c := int64(len(v))
		sum += c
		if *reqStart <= sum {
			// TODO: DOCUMENT SIDE EFFECT
			*reqStart = sum - c
			*reqEnd = sum
			flog.Infof("Sending frame %d", i)
			return v
		}
	}
	flog.Errorf("Unable to get cache store for key %s with start %d of %d b", key, *reqStart, sum)
	return nil
}
func filesize(key string) int64 {
	// TODO: REIMPLEMENT WITH BIN-SEARCH
	var sum int64 = 0
	for _, v := range cacheService.Redis[key] {
		c := int64(len(v))
		sum += c
	}
	return sum
}

func cacheItem(key string, frames [][]byte, reqStart, reqEnd int64, reqFrames chan []byte) error {
	flog.Infof("Caching...")
	if isItemCached(key) {
		flog.Errorf("Item is already cached")
	}
	cs := getCacheService()
	frameClusterSize := cs.configService.GetConfigValueInt64("FRAME_CLUSTER_SIZE")
	var startFrame int64 = 0
	var endFrame int64 = 0
	var curByte int64 = 0
	n := int64(len(frames))
	var songSegments [][]byte
	for startFrame < n {
		flog.Infof("start %d, end %d, cur %d, end %d", startFrame, endFrame, curByte, n)
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
	cs.Redis[key] = songSegments
	return nil
}
