package services

import (
	models "github.com/matttm/spoticli/spoticli-models"
)

type FileMetaInfoServiceWrap struct{}

func (s *FileMetaInfoServiceWrap) GetAllFilesOfType(fileTypeCd int) []*models.FileMetaInfo {
	return GetAllFilesOfType(fileTypeCd)
}
func (s *FileMetaInfoServiceWrap) GetFileById(id int) (*models.FileMetaInfo, error) {
	return GetFileById(id)
}
