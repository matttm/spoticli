package controllers_test

import (
	"testing"

	mock_services "github.com/matttm/spoticli/spoticli-backend/internal/services/mocks"
	"go.uber.org/mock/gomock"
)

func TestAudioController_GetPresignedUrl(t *testing.T) {
	ctrl := gomock.NewController(t)

	svc := mock_services.NewMockApiAudioService(ctrl)
}
