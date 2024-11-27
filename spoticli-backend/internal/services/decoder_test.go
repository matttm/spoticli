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
		{
			purpose:  "should not decode frame header when missing sync header",
			frame:    []byte{0xff, 0xdb, 0xb2, 0x60},
			expected: -1,
		},
		{
			purpose:  "should decode frame header",
			frame:    []byte{0xff, 0xfb, 0xb2, 0x60},
			expected: 627,
		},
		{
			purpose:  "should decode an ID3v1 tag",
			frame:    []byte{0x54, 0x41, 0x47, 0x54, 0x68},
			expected: 627,
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
