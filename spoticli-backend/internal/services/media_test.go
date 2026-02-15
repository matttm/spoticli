package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type frameCollectionTestCase struct {
	expectedFrames int
	file           []byte
	purpose        string
}

func TestMediaService_PartitionMp3Frames_Success(t *testing.T) {
	var table = []frameCollectionTestCase{}
	for _, v := range table {
		frames := PartitionMp3Frames(v.file)
		assert.Equal(t, v.expectedFrames, len(frames), v.purpose)
	}
}

func TestMediaService_PartitionMp3Frames_Panic(t *testing.T) {
	var badTable = []frameCollectionTestCase{
		{
			expectedFrames: -1,
			file:           []byte{0x49, 0x44, 0x43, 0xFF, 0xE0, 0x24, 0x15, 0xA2, 0x0F, 0xFE, 0xB3, 0x2D, 0xFF, 0xE0},
			purpose:        "Invalid frame header should result in zero frames",
		},
	}
	for _, v := range badTable {
		frames := PartitionMp3Frames(v.file)
		assert.Equal(t, 0, len(frames), v.purpose)
	}
}

func TestPartitionMp3Frames_EmptyAndPadding(t *testing.T) {
	frames := PartitionMp3Frames([]byte{})
	assert.Equal(t, 0, len(frames), "empty input should yield zero frames")

	padding := make([]byte, 100)
	frames = PartitionMp3Frames(padding)
	assert.Equal(t, 0, len(frames), "padding-only input should yield zero frames")
}

func TestPartitionMp3Frames_ID3v1Tag(t *testing.T) {
	b := make([]byte, 128)
	b[0] = 'T'
	b[1] = 'A'
	b[2] = 'G'
	frames := PartitionMp3Frames(b)
	assert.Equal(t, 0, len(frames), "ID3v1 tag at start should end partitioning with zero frames")
}
