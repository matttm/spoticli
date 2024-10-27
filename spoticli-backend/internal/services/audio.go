// Package services
package services

// GetPresignedUrl gets a presigned url
// for downloading an object from s3
func GetPresignedUrl(id int) (string, error) {
	t, _ := GetTrack(id)
	key := t.Title
	svc := GetStorageService()
	url, err := svc.GetPresignedUrl(key)
	if err != nil {
		panic(err)
	}
	return url, nil
}

// GetAudio gets a io.Reader containing an entire
// audio object on success and a *int referring to
// content's size
func GetAudio(id int) ([]byte, *int64, error) {
	t, _ := GetTrack(id)
	key := t.Title
	svc := GetStorageService()
	// TODO: rewrite and use getaudiopart
	body, err := svc.DownloadFile(key, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	length := int64(len(body))
	return body, &length, nil
}

// getAudioPart is a helper function in services which
// invokes the downloading and turning the reader into bytexs
//
// _range input, must be in the form "bytes=<start>-<end>"
func getAudioPart(id int) ([]byte, *int64, error) {
	t, _ := GetTrack(id)
	key := t.Title
	svc := GetStorageService()
	// TODO: put file in redis
	body, err := svc.DownloadFile(key, nil, nil)
	if err != nil {
		panic(err)
	}
	length := int64(len(body))
	return body, &length, nil
}

// StreamAudioSegment
func StreamAudioSegment(id int, start, end *int64) ([]byte, *int, *int, error) {
	t, _ := GetTrack(id)
	filesize := t.FileSize
	// key := t.Title
	if *start == 0 {
		*end = GetConfigValue[int64]("STREAM_SEGMENT_SIZE")
	}
	key := t.Title
	svc := GetStorageService()
	// TODO: rewrite and use getaudiopart
	segment, err := svc.DownloadFile(key, start, end)
	if err != nil {
		return nil, nil, nil, err
	}
	length := len(segment)
	return segment, &length, &filesize, nil
}
