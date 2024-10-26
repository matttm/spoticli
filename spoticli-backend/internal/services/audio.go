// Package services
package services

import (
	"errors"
	"fmt"
	"io"
)

// GetPresignedUrl gets a presigned url
// for downloading an object from s3
func GetPresignedUrl(id int) (string, error) {
	t, _ := GetTrack(id)
	key := t.Title
	svc := GetStorageService()
	req, err := svc.GetPresignedUrl(key)
	if err != nil {
		panic(err)
	}
	return req.URL, nil
}

// GetAudio gets a io.Reader containing an entire
// audio object on success and a *int referring to
// content's size
func GetAudio(id int) ([]byte, *int64, error) {
	t, _ := GetTrack(id)
	key := t.Title
	svc := GetStorageService()
	// TODO: rewrite and use getaudiopart
	res, err := svc.DownloadFile(key, nil)
	if err != nil {
		return nil, nil, err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}
	return body, res.ContentLength, nil
}

// getAudioPart is a helper function in services which
// invokes the downloading and turning the reader into bytexs
//
// _range input, must be in the form "bytes=<start>-<end>"
func getAudioPart(id int, _range string) ([]byte, *int64, error) {
	t, _ := GetTrack(id)
	key := t.Title
	svc := GetStorageService()
	// TODO: put file in redis
	res, err := svc.DownloadFile(key, &_range)
	if err != nil {
		return nil, nil, err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}
	fmt.Printf("Content-Range from AWS %s\n", *res.ContentRange)
	return body, res.ContentLength, nil
}

// StreamAudioSegment
func StreamAudioSegment(id int, start, end *int) ([]byte, *int, *int, error) {
	t, _ := GetTrack(id)
	filesize := t.FileSize
	// key := t.Title
	if *start == 0 {
		*end = GetConfigValue[int]("STREAM_SEGMENT_SIZE")
	}
	// TODO: PUT START/END VALIDATION IN A VALIDATOR
	fmt.Printf("start %d, end %d, filesize %d\n", *start, *end, filesize)
	if *start >= *end || *end > filesize+1 {
		return nil, nil, nil, errors.New("Invalid range header")
	}
	key := t.Title
	svc := GetStorageService()
	// TODO: rewrite and use getaudiopart
	res, err := svc.DownloadFile(key, nil)
	if err != nil {
		return nil, nil, nil, err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}
	// the following  blobk is in testing TODO:
	body = ReadID3v2Header(body)
	frames := PartitionMp3Frames(body)
	fmt.Printf("Frame count: %d\n", len(frames))
	// end test NOTE:
	fr := flatten(frames[30:70])
	length := len(fr)
	return fr, &length, &filesize, nil
}

func flatten(arr [][]byte) []byte {
	result := []byte{}
	for _, subarr := range arr {
		result = append(result, subarr...)
	}
	return result
}
