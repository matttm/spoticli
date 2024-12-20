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
			purpose:        "Expect an ID3v2 header to cause a panic",
		},
	}
	for _, v := range badTable {
		func() {
			defer func() {
				if recover() == nil {
					t.Error("The counter function did not panic")
				}
			}()
			_ = PartitionMp3Frames(v.file)
		}()
	}
}
