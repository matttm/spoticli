package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	mock_services "github.com/matttm/spoticli/spoticli-backend/internal/services/mocks"
	"go.uber.org/mock/gomock"
)

const presignedUrl = "aws.s3/resource"

func TestAudioController_GetPresignedUrl(t *testing.T) {
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
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	// Parse the response body (if needed)
	var responseBody []byte
	json.NewDecoder(w.Body).Decode(&responseBody)

}
