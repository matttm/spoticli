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

func getCurrentFrameHeaderLength(b []byte) int {
	frameHeader := b[:4]
	fmt.Printf("%x \n", frameHeader)
	// first 11 bits are sync word, so skip them
	mpegVersion := ((frameHeader[1] >> 4) & 0b11)
	fmt.Printf("MPEG Version %02b\n", mpegVersion)
	layerDesc := (frameHeader[1] >> 1) & 0b11 // getting bits 5 and 6 as xx
	fmt.Printf("MPEG Layer %02b\n", layerDesc)

	versionStr := constants.VersionMap[int(mpegVersion)]
	layerStr := constants.LayerMap[int(layerDesc)]
	versionLayerStr := fmt.Sprintf("%s,%s", versionStr, layerStr)
	fmt.Printf("versionLayerString %s\n", versionLayerStr)

	bitRateIndex := (frameHeader[2] >> 4) & 0b1111
	bitRate := constants.BitrateMap[bitRateIndex][versionLayerStr]
	if _, err := fmt.Printf("BitRateIndex %04b bitRate %d \n", bitRateIndex, bitRate); err != nil {
		panic(err)
	}

	samplingRateIndex := (frameHeader[2]) & 0b1111
	samplingRate := constants.SamplingRateMap[samplingRateIndex][versionStr]
	if _, err := fmt.Printf("Sampling Rate Index %04b sampling rate %d \n", samplingRateIndex, samplingRate); err != nil {
		panic(err)
	}
	padding := frameHeader[3] >> 7
	fmt.Printf("padding bit %01b\n", padding)
	// For Layer I files us this formula:
	//
	//	FrameLengthInBytes = (12 * BitRate / SampleRate + Padding) * 4
	//
	// For Layer II & III files use this formula:
	//
	//	FrameLengthInBytes = 144 * BitRate / SampleRate + Padding
	//  err := os.Stdout.Sync()
	//  if err != nil {
	//  	panic(err)
	//  } // Force flush the output buffer
	if layerDesc == 0b11 { // it is L1
		return (12*bitRate*1000/samplingRate + int(padding)) * 4
	} else {
		return 144*bitRate*1000/samplingRate + int(padding)
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
		currentFrameLength := getCurrentFrameHeaderLength(b)
		clip := b[:currentFrameLength]
		frames = append(frames, clip)
		b = b[currentFrameLength:]
		if len(b) <= 0 {
			break
		}
		fmt.Printf("Frame count %d frame length %d \n\n", len(frames), currentFrameLength)
	}
	return frames
}
