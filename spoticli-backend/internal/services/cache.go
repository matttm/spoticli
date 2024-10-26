// Package services
package services

type CacheService struct {
	Redis map[int][][]byte // map of id to song, which is cut in segments
}

var cacheService *CacheService

func getCacheService() *CacheService {
	if cacheService == nil {
		cacheService = &CacheService{}
	}
	return cacheService
}

func isItemCached(id int) bool {
	_, ok := getCacheService().Redis[id]
	return ok
}

func getSegmentFromCache(id, segment int) []byte {
	return getCacheService().Redis[id][segment]
}

func cacheItem(id int, frames [][]byte) error {
	if isItemCached(id) {
		panic("Item is already cached")
	}
	frameClusterSize := GetConfigValue[int]("FRAME_CLUSTER_SIZE")
	cur := 0
	end := 0
	n := len(frames)
	var songSlices [][]byte
	for cur < n {
		end = min(cur+frameClusterSize, n)
		b := flatten(frames[cur:end])
		songSlices = append(songSlices, b)
		cur += frameClusterSize
	}
	getCacheService().Redis[id] = songSlices
	return nil
}
