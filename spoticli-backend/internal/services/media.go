package services

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"slices"
)

type MediaService struct {
}

var mediaService *MediaService

func GetMusiceService() *MediaService {
	if mediaService == nil {
		mediaService = &MediaService{}
		println("MusicService Instantiated")
	}
	return mediaService
}

// ReadID3v2Header takes a []byte and returns a
// []byte without an ID3 header
func ReadID3v2Header(b []byte) []byte {
	idString := "ID3"
	identifier := b[:3]
	if !bytes.Equal(identifier, []byte(idString)) {
		panic("No ID3v2 indicator found")
	}
	major := int(b[3])
	revision := int(b[4])
	fmt.Printf("Major Version %d \nRevision %d\n", major, revision)
	flags := b[5]
	fmt.Printf("flags byte %08b\n", flags)
	// the following size bytes are sync safe so 7-bits
	//  s1 := int(b[6])
	//  s2 := int(b[7])
	s3 := int(b[8])
	s4 := int(b[9])
	size := s3*int(math.Pow(2, 7)) + s4 + 10
	fmt.Printf("ID3v2 tag size is %d bytes\n", size)
	return b[size+1:]
}

// PartitionMp3Frames takes an entire
// mp3 file andreturns a slice of frames
func PartitionMp3Frames(b []byte) [][]byte {
	if len(b) == 0 {
		return [][]byte{}
	}
	syncStart, _ := hex.DecodeString("FF")
	// For Layer I files us this formula:
	//
	//	FrameLengthInBytes = (12 * BitRate / SampleRate + Padding) * 4
	//
	// For Layer II & III files use this formula:
	//
	//	FrameLengthInBytes = 144 * BitRate / SampleRate + Padding
	NextFrameStart := func(b []byte) int {
		ffIndex := bytes.LastIndex(b, syncStart) // ff is 1111 1111, so now check left and right byte for 111
		BIT_COUNT_TO_BE_FOUND := 3
		leftByte := b[ffIndex-1]
		rightByte := b[ffIndex+1]
		fmt.Printf("%08b, %08b, %08b", leftByte, b[ffIndex], rightByte)

		BITS_FOUND := 0
		// for i := range BIT_COUNT_TO_BE_FOUND {
		// 	if leftByte&(1<<i) == 1 {
		// 		BITS_FOUND += 1
		// 		if BITS_FOUND == BIT_COUNT_TO_BE_FOUND {
		// 			return ffIndex
		// 		}
		// 	} else {
		// 		break
		// 	}
		// }
		BITS_FOUND = 0
		for i := range BIT_COUNT_TO_BE_FOUND {
			if (rightByte>>(8-i-1))&1 == 1 {
				BITS_FOUND += 1
				if BITS_FOUND == BIT_COUNT_TO_BE_FOUND {
					return ffIndex
				}
			} else {
				break
			}
		}
		fmt.Println("No possible sync seen")
		return -1
	}
	var frames [][]byte
	endIndex := len(b)
	startIndex := NextFrameStart(b) // as asserted by above, this is a sync, so we'll add the length of sync togo to  header of following frame
	for startIndex != -1 {
		clip := b[startIndex:endIndex]
		frames = append(frames, clip)
		b = b[:startIndex]
		startIndex = NextFrameStart(b)
		if startIndex == -1 {
			break
		}
		endIndex = len(b)
	}
	slices.Reverse(frames)
	return frames
}
