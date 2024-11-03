package services

import (
	"testing"

	mock_services "github.com/matttm/spoticli/spoticli-backend/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const TEST_KEY = "test"

func Test_CacheService_CacheItem_Success(t *testing.T) {
	cs := getCacheService()
	data := [][]byte{
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// inject mock
	configSvc := mock_services.NewMockConfigServiceApi(ctrl)
	cs.configService = configSvc
	configSvc.EXPECT().GetConfigValueInt64("FRAME_CLUSTER_SIZE").Return(int64(2)).AnyTimes()
	var channel chan []byte
	_ = cacheItem(TEST_KEY, data, int64(4), int64(0), channel)
	res := cs.Redis[TEST_KEY]
	assert.Equal(t, 8, len(res), "Expected %d frame clusters: got %d", 8, len(res))

}
