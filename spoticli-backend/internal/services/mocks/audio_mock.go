// Code generated by MockGen. DO NOT EDIT.
// Source: internal/services/audio_api.go
//
// Generated by this command:
//
//	mockgen -source=internal/services/audio_api.go -destination=internal/services/mocks/audio_mock.go
//

// Package mock_services is a generated GoMock package.
package mock_services

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockApiAudioService is a mock of ApiAudioService interface.
type MockApiAudioService struct {
	ctrl     *gomock.Controller
	recorder *MockApiAudioServiceMockRecorder
	isgomock struct{}
}

// MockApiAudioServiceMockRecorder is the mock recorder for MockApiAudioService.
type MockApiAudioServiceMockRecorder struct {
	mock *MockApiAudioService
}

// NewMockApiAudioService creates a new mock instance.
func NewMockApiAudioService(ctrl *gomock.Controller) *MockApiAudioService {
	mock := &MockApiAudioService{ctrl: ctrl}
	mock.recorder = &MockApiAudioServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApiAudioService) EXPECT() *MockApiAudioServiceMockRecorder {
	return m.recorder
}

// GetAudio mocks base method.
func (m *MockApiAudioService) GetAudio(id int) ([]byte, *int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAudio", id)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(*int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetAudio indicates an expected call of GetAudio.
func (mr *MockApiAudioServiceMockRecorder) GetAudio(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAudio", reflect.TypeOf((*MockApiAudioService)(nil).GetAudio), id)
}

// GetPresignedUrl mocks base method.
func (m *MockApiAudioService) GetPresignedUrl(id int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPresignedUrl", id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPresignedUrl indicates an expected call of GetPresignedUrl.
func (mr *MockApiAudioServiceMockRecorder) GetPresignedUrl(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPresignedUrl", reflect.TypeOf((*MockApiAudioService)(nil).GetPresignedUrl), id)
}

// StreamAudioSegment mocks base method.
func (m *MockApiAudioService) StreamAudioSegment(id int, start, end *int64) ([]byte, *int, *int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StreamAudioSegment", id, start, end)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(*int)
	ret2, _ := ret[2].(*int64)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// StreamAudioSegment indicates an expected call of StreamAudioSegment.
func (mr *MockApiAudioServiceMockRecorder) StreamAudioSegment(id, start, end any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StreamAudioSegment", reflect.TypeOf((*MockApiAudioService)(nil).StreamAudioSegment), id, start, end)
}