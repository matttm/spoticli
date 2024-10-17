package services

import (
	"io"
)

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
func GetAudioPart(id int, _range string) ([]byte, *int64, error) {
	t, _ := GetTrack(id)
	key := t.Title
	svc := GetStorageService()
	res, err := svc.DownloadFile(key, &_range)
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
