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
			purpose:  "should not decode frame header when missing sync header",
			frame:    []byte{0xff, 0xfb, 0x92, 0x0c},
			expected: 418,
		},
		{
			purpose:  "should decode frame header",
			frame:    []byte{0xff, 0xfb, 0xb2, 0x60},
			expected: 627,
		},
		{
			purpose:  "should decode an ID3v1 tag",
			frame:    []byte{0x54, 0x41, 0x47, 0x54, 0x68, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expected: 0,
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

type readID3v2HeaderTestCase struct {
	purpose  string
	input    []byte
	expected []byte
	panics   bool
}

func TestReadID3v2Header(t *testing.T) {
	var table = []readID3v2HeaderTestCase{
		{
			purpose: "should read ID3v2 header with minimal size",
			input: []byte{
				'I', 'D', '3', // ID3 identifier
				0x03, 0x00, // version 3.0
				0x00,                   // flags
				0x00, 0x00, 0x00, 0x00, // size (10 bytes total - just header)
				'r', 'e', 's', 't', // remaining data
			},
			expected: []byte{'r', 'e', 's', 't'},
			panics:   false,
		},
		{
			purpose: "should read ID3v2 header with size calculation",
			input: []byte{
				'I', 'D', '3', // ID3 identifier
				0x04, 0x00, // version 4.0
				0x00,                   // flags
				0x00, 0x00, 0x00, 0x0A, // size (10 + 10 = 20 bytes)
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // padding
				'd', 'a', 't', 'a', // remaining data
			},
			expected: []byte{'d', 'a', 't', 'a'},
			panics:   false,
		},
		{
			purpose: "should panic when declared size is larger than available data",
			input: []byte{
				'I', 'D', '3', // ID3 identifier
				0x03, 0x00, // version 3.0
				0x80,                   // flags
				0x00, 0x00, 0x01, 0x00, // size (128 + 10 = 138 bytes)
			},
			expected: nil,
			panics:   true,
		},
		{
			purpose: "should panic when ID3 identifier is missing",
			input: []byte{
				'I', 'D', '2', // wrong identifier
				0x03, 0x00,
				0x00,
				0x00, 0x00, 0x00, 0x00,
			},
			expected: nil,
			panics:   true,
		},
		{
			purpose: "should panic with invalid identifier",
			input: []byte{
				'X', 'Y', 'Z', // wrong identifier
				0x03, 0x00,
				0x00,
				0x00, 0x00, 0x00, 0x00,
			},
			expected: nil,
			panics:   true,
		},
	}

	for _, v := range table {
		t.Run(v.purpose, func(t *testing.T) {
			if v.panics {
				result, err := ReadID3v2Header(v.input)
				assert.Error(t, err, v.purpose)
				assert.Nil(t, result, v.purpose)
			} else {
				result, err := ReadID3v2Header(v.input)
				assert.NoError(t, err, v.purpose)
				assert.Equal(t, v.expected, result, v.purpose)
			}
		})
	}
}

func TestReadID3v2Header_TruncatedAndEmpty_Panics(t *testing.T) {
	// very small input will cause slices in the function to panic (index OOB)
	assert.Panics(t, func() { _, _ = ReadID3v2Header([]byte{}) })
	assert.Panics(t, func() { _, _ = ReadID3v2Header([]byte{'I', 'D'}) })
}
