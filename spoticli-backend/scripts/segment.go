package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/matttm/spoticli/spoticli-backend/internal/services"
)

const DIRECTORY_NAME = "segments"
const SEGMENT_LIMIT = 20

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	filename := flag.String("file", "", "file to segment")

	flag.Parse()
	*filename = filepath.Join(homeDir, *filename)
	if *filename == "" {
	}
	fmt.Printf("opening dir %s...\n", *filename)
	file, err := os.Open(*filename)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	b = services.ReadID3v2Header(b)
	frames := services.PartitionMp3Frames(b)
	if len(frames) <= 0 {
		panic("No frames")
	}
	fmt.Printf("Frame count %d", len(frames))

	// make dir if itdowsnt exist
	_, err = os.Stat(DIRECTORY_NAME)
	if err != nil {
		os.Mkdir(DIRECTORY_NAME, 777)
	}
	for i := range min(len(frames), SEGMENT_LIMIT) {
		*filename = fmt.Sprintf("%s/%d.frame", DIRECTORY_NAME, i)
		file, err = os.Create(*filename)
		defer file.Close()
		if err != nil {
			panic(err)
		}
		file.Write(frames[i])
	}

}
