package services

import "testing"

const TEST_KEY = "test"

func Test_CacheService_CacheItem_Success(t *testing.T) {
	cs := getCacheService()
	cs.Redis[TEST_KEY] = [][]byte{
		[]byte("frame:one"),
		[]byte("frame:two"),
		[]byte("frame:three"),
		[]byte("frame:four"),
		[]byte("frame:five"),
		[]byte("frame:six"),
		[]byte("frame:seven"),
		[]byte("frame:eight"),
		[]byte("frame:nine"),
		[]byte("frame:A"),
		[]byte("frame:B"),
		[]byte("frame:C"),
		[]byte("frame:D"),
		[]byte("frame:E"),
		[]byte("frame:F"),
	}
}
