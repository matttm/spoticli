package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type frameLengthTestCase struct {
	purpose  string
	expected int
	frame    []byte
}

func TestDecoderService_DetermineFrameLength(t *testing.T) {
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
