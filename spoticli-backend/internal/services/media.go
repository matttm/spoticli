package services

import (
	"bytes"
	"fmt"

	"github.com/matttm/spoticli/spoticli-backend/internal/constants"
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
	// nhe following size bytes are sync safe so 7-bits
	s1 := int(b[6])
	s2 := int(b[7])
	s3 := int(b[8])
	s4 := int(b[9])

	size := (s1 << (7 * 3)) + (s2 << (7 * 2)) + (s3 << 7) + s4 + 10
	fmt.Printf("ID3v2 tag size is %d bytes\n", size)
	return b[size:]
}

func getNextFrameHeaderIndex(b []byte) int {
	frameHeader := b[:32]
	// first 11 bits are sync word, so skip them
	mpegVersion := ((frameHeader[2] & 1) << 1) | (frameHeader[3] >> 7)
	fmt.Printf("MPEG Version %02b\n", mpegVersion)
	layerDesc := (frameHeader[5] >> 5) & 3 // getting bits 5 and 6 as xx
	fmt.Printf("MPEG Layer %02b\n", layerDesc)

	versionStr := constants.VersionMap[int(mpegVersion)]
	layerStr := constants.LayerMap[int(layerDesc)]
	versionLayerStr := fmt.Sprintf("%s,%s", versionStr, layerStr)
	bitRateIndex := frameHeader[5] & 15
	bitRate := constants.BitrateMap[bitRateIndex][versionLayerStr]
	samplingRateIndex := (frameHeader[6] >> 4) & 15
	samplingRate := constants.SamplingRateMap[samplingRateIndex][versionStr]
	fmt.Printf("BitRate %d \n Sampling rate %d \n", bitRate, samplingRate)
	// For Layer I files us this formula:
	//
	//	FrameLengthInBytes = (12 * BitRate / SampleRate + Padding) * 4
	//
	// For Layer II & III files use this formula:
	//
	//	FrameLengthInBytes = 144 * BitRate / SampleRate + Padding
	if layerDesc == 0b11 { // it is L1
		return (12*bitRate/samplingRate + 0) * 4
	} else {
		return 144*bitRate/samplingRate + 0
	}
}

// PartitionMp3Frames takes an entire
// mp3 file andreturns a slice of frames
func PartitionMp3Frames(b []byte) [][]byte {
	if len(b) == 0 {
		return [][]byte{}
	}
	var frames [][]byte
	for {
		nextStartIndex := getNextFrameHeaderIndex(b)
		break
		if nextStartIndex > len(b) {
			break
		}
		clip := b[:nextStartIndex]
		frames = append(frames, clip)
		b = b[nextStartIndex:]
		if len(b) <= 0 {
			break
		}
	}
	return frames
}
