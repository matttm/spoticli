package services

import (
	"bytes"
	"encoding/hex"
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
		panic("No ID3 indicator found")
	}
	syncStart, err := hex.DecodeString("0FFE")
	if err != nil {
		panic(err)
	}
	index := bytes.Index(b, syncStart)
	if index < 0 {
		panic("Cannot find sync header")
	}
	return b[index:]
}

// PartitionMp3Frames takes an entire
// mp3 file andreturns a slice of frames
func PartitionMp3Frames(b []byte) [][]byte {
	if len(b) == 0 {
		return [][]byte{}
	}
	syncStartv1, _ := hex.DecodeString("0FFE")
	syncStartv2, _ := hex.DecodeString("FFE0")
	if !bytes.HasPrefix(b, syncStartv1) && !bytes.HasPrefix(b, syncStartv2) {
		panic("Invalid Input: b does not start with sync header")
	}
	NextFrameStart := func(b []byte) int {
		return max(
			bytes.LastIndex(b, syncStartv1),
			bytes.LastIndex(b, syncStartv2),
		)
	}
	var frames [][]byte
	endIndex := len(b)
	startIndex := NextFrameStart(b) // as asserted by above, this is a sync, so we'll add the length of sync togo to  header of following frame
	for startIndex != -1 {
		clip := b[startIndex:endIndex]
		frames = append(frames, clip)
		b = b[:startIndex]
		startIndex = NextFrameStart(b)
		endIndex = len(b)
	}
	slices.Reverse(frames)
	return frames
}
