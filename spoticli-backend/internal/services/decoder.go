package services

import (
	"bytes"
	"fmt"

	"github.com/coder/flog"

	"github.com/matttm/spoticli/spoticli-backend/internal/constants"
)

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
	flog.Infof("Major Version %d Revision %d", major, revision)
	flags := b[5]
	flog.Infof("flags byte %08b", flags)
	// nhe following size bytes are sync safe so 7-bits
	s1 := int(b[6])
	s2 := int(b[7])
	s3 := int(b[8])
	s4 := int(b[9])

	size := (s1 << (7 * 3)) + (s2 << (7 * 2)) + (s3 << 7) + s4 + 10
	flog.Infof("ID3v2 tag size is %d bytes", size)
	return b[size:]
}

// getCurrentFrameLength
//
//	in:
//	  b []byte - the entire frame with header and ID3v2 stripped off
//	out:
//	  x int - the frames length in bytes, or -1 iff frame is invalid
func getCurrentFrameLength(b []byte) int {
	frameHeader := b[:4]
	flog.Infof("%x ", frameHeader)
	// check if this an ID3v1 tag
	if string(frameHeader[:3]) == "TAG" {
		flog.Successf("Identified ID3v1 tag")
		title := b[3:33]
		artist := b[33:63]
		album := b[63:93]
		flog.Successf("Song title: %s", string(title))
		flog.Successf("Artist: %s", string(artist))
		flog.Successf("Album: %s", string(album))
		return 0
	}
	// check for sync frame (eleven sequential ones)
	if !(frameHeader[0] == 0xFF && frameHeader[1]>>5 == 0b111) {
		flog.Errorf("Loaded frame has improper sync header")
		return -1
	}
	// first 11 bits are sync word, so skip them
	mpegVersion := ((frameHeader[1] >> 4) & 0b11)
	flog.Infof("MPEG Version %02b", mpegVersion)
	layerDesc := (frameHeader[1] >> 1) & 0b11 // getting bits 5 and 6 as xx
	flog.Infof("MPEG Layer %02b", layerDesc)

	versionStr := constants.VersionMap[int(mpegVersion)]
	layerStr := constants.LayerMap[int(layerDesc)]
	versionLayerStr := fmt.Sprintf("%s,%s", versionStr, layerStr)
	flog.Infof("versionLayerString %s", versionLayerStr)

	bitRateIndex := (frameHeader[2] >> 4) & 0b1111
	bitRate := constants.BitrateMap[bitRateIndex][versionLayerStr]

	flog.Infof("BitRateIndex %04b bitRate %d ", bitRateIndex, bitRate)
	samplingRateIndex := (frameHeader[2] >> 2) & 0b11
	samplingRate := constants.SamplingRateMap[samplingRateIndex][versionStr]

	flog.Infof("Sampling Rate Index %02b sampling rate %d ", samplingRateIndex, samplingRate)
	padding := (frameHeader[2] >> 1) & 0b1
	flog.Infof("padding bit %01b ", padding)
	// For Layer I files us this formula:
	//
	//	FrameLengthInBytes = (12 * BitRate / SampleRate + Padding) * 4
	//
	// For Layer II & III files use this formula:
	//
	//	FrameLengthInBytes = 144 * BitRate / SampleRate + Padding
	if layerDesc == 0b11 { // it is L1
		return (12*bitRate*1000/samplingRate + int(padding)) * 4
	} else {
		return 144*bitRate*1000/samplingRate + int(padding)
	}
}
