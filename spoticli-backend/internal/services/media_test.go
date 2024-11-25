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
type frameLengthTestCase struct {
	purpose  string
	expected int
	frame    []byte
}

func TestMediaService_PartitionMp3Frames_Success(t *testing.T) {
	var table = []frameCollectionTestCase{
		frameCollectionTestCase{
			expectedFrames: 2,
			file:           []byte{0xFF, 0xE0, 0x24, 0x15, 0xA2, 0x0F, 0xFE, 0xB3, 0x2D},
			purpose:        "Expect file to be processed starting with 0xFFE0",
		},
		frameCollectionTestCase{
			expectedFrames: 2,
			file:           []byte{0x0F, 0xFE, 0x24, 0x15, 0xA2, 0x0F, 0xFE, 0xB3, 0x2D},
			purpose:        "Expect file to be processed starting with 0x0FFE",
		},
		frameCollectionTestCase{
			expectedFrames: 3,
			file:           []byte{0xFF, 0xE0, 0x24, 0x15, 0xA2, 0x0F, 0xFE, 0xB3, 0x2D, 0xFF, 0xE0},
			purpose:        "Expect a trailing emptyy frame to increase frame count",
		},
		frameCollectionTestCase{
			expectedFrames: 0,
			file:           []byte{},
			purpose:        "Expect empty slice to not cause error",
		},
	}
	for _, v := range table {
		frames := PartitionMp3Frames(v.file)
		assert.Equal(t, v.expectedFrames, len(frames), v.purpose)
	}
}

func TestMediaService_PartitionMp3Frames_Panic(t *testing.T) {
	var badTable = []frameCollectionTestCase{
		frameCollectionTestCase{
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
func TestMediaService_DetermineFrameLength(t *testing.T) {
	var table = []frameLengthTestCase{
		{
			purpose:  "should decode frame header",
			frame:    []byte{0xff, 0xfb, 0x50, 0x00},
			expected: 208,
		},
		{
			purpose:  "should decode frame header",
			frame:    []byte{0xff, 0xfb, 0x90, 0x64},
			expected: 417,
		},
		{
			purpose:  "should decode frame header with padding",
			frame:    []byte{0xff, 0xfb, 0x92, 0x64},
			expected: 418,
		},
	}
	for _, v := range table {
		assert.Equal(
			t,
			getCurrentFrameLength(v.frame),
			v.expected,
			v.purpose,
		)
	}
}
