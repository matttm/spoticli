package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	expectedFrames int
	file           []byte
	purpose        string
}

var table []testCase = []testCase{
	testCase{
		expectedFrames: 2,
		file:           []byte{0xFF, 0xE0, 0x24, 0x15, 0xA2, 0x0F, 0xFE, 0xB3, 0x2D},
		purpose:        "Expect file to be processed starting with 0xFFE0",
	},
	testCase{
		expectedFrames: 2,
		file:           []byte{0x0F, 0xFE, 0x24, 0x15, 0xA2, 0x0F, 0xFE, 0xB3, 0x2D},
		purpose:        "Expect file to be processed starting with 0x0FFE",
	},
	testCase{
		expectedFrames: 3,
		file:           []byte{0xFF, 0xE0, 0x24, 0x15, 0xA2, 0x0F, 0xFE, 0xB3, 0x2D, 0xFF, 0xE0},
		purpose:        "Expect a trailing emptyy frame to increase frame count",
	},
	testCase{
		expectedFrames: 0,
		file:           []byte{},
		purpose:        "Expect empty slice to not cause error",
	},
}

func TestMediaService_PartitionMp3Frames_Success(t *testing.T) {
	for _, v := range table {
		frames := PartitionMp3Frames(v.file)
		assert.Equal(t, v.expectedFrames, len(frames), v.purpose)
	}
}

var badTable = []testCase{
	testCase{
		expectedFrames: -1,
		file:           []byte{0x49, 0x44, 0x43, 0xFF, 0xE0, 0x24, 0x15, 0xA2, 0x0F, 0xFE, 0xB3, 0x2D, 0xFF, 0xE0},
		purpose:        "Expect an ID3v2 header to cause a panic",
	},
}

func TestMediaService_PartitionMp3Frames_Panic(t *testing.T) {
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
