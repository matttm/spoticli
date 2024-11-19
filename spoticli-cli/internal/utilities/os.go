package utilities

import (
	"fmt"
	"os"
	"strings"

	"github.com/matttm/spoticli/spoticli-cli/internal/models"
)

func CollectFiles(path, fileExt string) []*models.FileInfo {
	files := make([]*models.FileInfo, 0)
	collectFilesFromFileInfo(path, &files, fileExt)
	return files
}

func collectFilesFromFileInfo(dsc string, files *[]*models.FileInfo, fileExt string) {
	info, err := os.Stat(dsc)
	if err != nil {
		panic(err)
	}
	if info.IsDir() {
		_files, err := os.ReadDir(dsc)
		if err != nil {
			panic(err)
		}
		for _, file := range _files {
			abspath := fmt.Sprintf("%s/%s", dsc, file.Name())
			//  fmt.Printf("recursing on %s\n", abspath)
			collectFilesFromFileInfo(abspath, files, fileExt)
		}
	} else if strings.Contains(dsc, fileExt) {
		*files = append(*files, &models.FileInfo{Path: dsc, Size: info.Size()})
	}
}
