package services

import (
	"github.com/matttm/spoticli/spoticli-backend/internal/database"
	models "github.com/matttm/spoticli/spoticli-models"
)

func GetAllFilesOfType(fileTypeCd int) []*models.FileMetaInfo {
	return database.SelectAllFileMetaInfo()
}
