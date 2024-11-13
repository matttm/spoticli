// Package services
package services

import "github.com/matttm/spoticli/spoticli-backend/internal/database"

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
	body, err := svc.DownloadFile(key, nil)
	if err != nil {
		return nil, nil, err
	}
	length := int64(len(body))
	return body, &length, nil
}

// StreamAudioSegment
func StreamAudioSegment(id int, start, end *int64) ([]byte, *int, *int64, error) {
	t, _ := GetTrack(id)
	var filesize int64
	// key := t.Title
	if *start == 0 {
		*end = GetConfigService().GetConfigValueInt64("STREAM_SEGMENT_SIZE")
	}
	if *start >= int64(t.FileSize) {
		panic("Invalid start pos")
		var b []byte
		return b, nil, nil, nil
	}
	key := t.Title
	svc := GetStorageService()
	// TODO: rewrite and use getaudiopart
	segment, filesize, err := svc.StreamFile(key, start, end)
	if err != nil {
		return nil, nil, nil, err
	}
	length := len(segment)
	return segment, &length, &filesize, nil
}
func UploadMusicThroughPresigned(track_name string) string {
	db := database.GetDatabase()
	svc := GetStorageService()
	tx, _ := db.Begin()
	database.InsertFileMetaInfo(tx, track_name, *TRACKS_BUCKET_NAME, 1)
	url, err := svc.PostPresignedUrl(track_name)
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	// TODO: delegate tx finalizatipn to bg task to check for upload
	_ = tx.Commit()
	return *url
}
