package controllers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	mock_services "github.com/matttm/spoticli/spoticli-backend/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const presignedUrl = "aws.s3/resource"

func TestAudioController_GetPresignedUrl_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mock_services.NewMockApiAudioService(ctrl)
	audioService = svc

	// Create a new request recorder
	w := httptest.NewRecorder()
	// Create a mock request
	req := httptest.NewRequest(http.MethodGet, "/audio/1", nil)
	svc.EXPECT().GetPresignedUrl(1).Return(presignedUrl, nil).Times(1)

	// Create a new test server
	handler := mux.NewRouter()
	handler.Path("/audio/{id:[0-9]+}").Methods(http.MethodGet).HandlerFunc(GetPresignedUrl)
	server := httptest.NewServer(handler)
	defer server.Close()

	handler.ServeHTTP(w, req) // Assert on the status code
	assert.Equalf(t, w.Code, http.StatusOK, "Expected status code 200, got %d", w.Code)

	// Parse the response body (if needed)
	var got = w.Body.String()
	assert.Equalf(t, presignedUrl, got, "Expected presigned url %s, got %s", presignedUrl, got)
}
func TestAudioController_GetAudio_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mock_services.NewMockApiAudioService(ctrl)
	audioService = svc

	// Create a new request recorder
	w := httptest.NewRecorder()
	// Create a mock request
	audio := "Some fake audio data"
	b := []byte(audio)
	req := httptest.NewRequest(http.MethodGet, "/audio/proxy/1", nil)
	length := int64(len(b))
	svc.EXPECT().GetAudio(1).Return(b, &length, nil).Times(1)

	// Create a new test server
	handler := mux.NewRouter()
	handler.Path("/audio/proxy/{id:[0-9]+}").Methods(http.MethodGet).HandlerFunc(GetAudio)
	server := httptest.NewServer(handler)
	defer server.Close()

	handler.ServeHTTP(w, req) // Assert on the status code
	assert.Equalf(t, w.Code, http.StatusOK, "Expected status code 200, got %d", w.Code)

	// Parse the response body (if needed)
	data, _ := io.ReadAll(w.Body)
	body := string(data)
	h := w.Header()
	assert.Equal(t, "audio/mp3", h.Get("Content-Type"), "Expected Content-Type of 'audio/mp3'")
	assert.Equal(t, fmt.Sprint(length), h.Get("Content-Length"), "Expected Content-Length of bleh")
	assert.Equal(t, audio, body, "Expected")
}
func TestAudioController_GetAudioSegment_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mock_services.NewMockApiAudioService(ctrl)
	audioService = svc

	// Create a new request recorder
	w := httptest.NewRecorder()
	// Create a mock request
	audio := "Some fake audio data"
	b := []byte(audio)
	var start int64 = 0
	var end int64 = 0
	_range := fmt.Sprintf("bytes=%d-%d", start, end)
	part := b[:1]
	length := len(part)
	req := httptest.NewRequest(http.MethodGet, "/audio/proxy/stream/1", nil)
	filesize := int64(len(b))
	svc.EXPECT().StreamAudioSegment(1, &start, &end).Return(part, &length, &filesize, nil).Times(1)

	// Create a new test server
	handler := mux.NewRouter()
	handler.Path("/audio/proxy/stream/{id:[0-9]+}").Methods(http.MethodGet).HandlerFunc(GetAudioPart)
	server := httptest.NewServer(handler)
	defer server.Close()

	// setting range
	req.Header.Add("Range", _range)
	handler.ServeHTTP(w, req) // Assert on the status code
	assert.Equalf(t, w.Code, http.StatusPartialContent, "Expected status code 200, got %d", w.Code)

	// Parse the response body (if needed)
	data, _ := io.ReadAll(w.Body)
	h := w.Header()
	assert.Equal(t, "audio/mp3", h.Get("Content-Type"), "Expected Content-Type of 'audio/mp3'")
	assert.Equal(t, fmt.Sprint(length), h.Get("Content-Length"), "Expected Content-Length of bleh")
	assert.Equal(t, string(data), "S", "Expected []byte with content x: got y")
}
