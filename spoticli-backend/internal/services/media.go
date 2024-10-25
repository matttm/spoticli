package services

import (
	"bytes"
	"encoding/hex"

	models "github.com/matttm/spoticli/spoticli-models"
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
func GetTrack(id int) (*models.Track, error) {
	t := new(models.Track)
	t.Title = "bat_country.mp3"
	t.FileSize = 7523726
	return t, nil
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
	syncStart, err := hex.DecodeString("0FFE")
	if err != nil {
		panic(err)
	}
	if !bytes.HasPrefix(b, syncStart) {
		panic("Invalid Input: b does not start with sync header")
	}
	var frames [][]byte
	syncLen := len(syncStart)
	startIndex := 0
	endIndex := bytes.Index(b[syncLen:], syncStart) // as asserted by above, this is a sync, so we'll add the length of sync togo to  header of following frame
	for endIndex != -1 {
		segment := b[startIndex:endIndex]
		frames = append(frames, bytes.Clone(segment))
		startIndex = endIndex
		b = b[startIndex:]
		endIndex = bytes.Index(b[syncLen:], syncStart) // as asserted by above, this is 0
	}
	// append the last frame e (idx algorithm couldnt find end boundary
	segment := b
	frames = append(frames, bytes.Clone(segment))
	return frames
}
