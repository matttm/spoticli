package services

import (
	models "github.com/matttm/spoticli/spoticli-models"
)

type FileMetaInfoServiceApi interface {
	GetAllFilesOfType(fileTypeCd int) []*models.FileMetaInfo
	GetFileById(id int) (*models.FileMetaInfo, error)
}
