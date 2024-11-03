package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	expectedFrames int
	file           []byte
}

var table []testCase = []testCase{
	testCase{
		expectedFrames: 2,
		file:           []byte{0xFF, 0xE0, 0x24, 0x15, 0xA2, 0x0F, 0xFE, 0xB3, 0x2D},
	},
	testCase{
		expectedFrames: 2,
		file:           []byte{0x0F, 0xFE, 0x24, 0x15, 0xA2, 0x0F, 0xFE, 0xB3, 0x2D},
	},
	testCase{
		expectedFrames: 3,
		file:           []byte{0xFF, 0xE0, 0x24, 0x15, 0xA2, 0x0F, 0xFE, 0xB3, 0x2D, 0xFF, 0xE0},
	},
	testCase{
		expectedFrames: 0,
		file:           []byte{},
	},
}

func TestMediaService_PartitionMp3Frames_Success(t *testing.T) {
	for _, v := range table {
		frames := PartitionMp3Frames(v.file)
		assert.Equal(t, v.expectedFrames, len(frames), "Expected")
	}
}
