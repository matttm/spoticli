package main

import (
	"fmt"
	"log"
	"os"

	"github.com/matttm/spoticli/spoticli-cli/internal/handler"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "auth",
				Aliases: []string{},
				Usage:   "authentication",
				Subcommands: []*cli.Command{
					{
						Name:  "login",
						Usage: "authenticate yourself",
						Action: func(cCtx *cli.Context) error {
							fmt.Println("Feature not implemented")
							return nil
						},
					},
					{
						Name:  "logout",
						Usage: "revoke self-authentication",
						Action: func(cCtx *cli.Context) error {
							fmt.Println("Feature not implemented")
							return nil
						},
					},
				},
			},
			{
				Name:    "song",
				Aliases: []string{},
				Usage:   "song <action>",
				Subcommands: []*cli.Command{
					{
						Name:        "upload",
						Usage:       "upload <path>",
						Description: "Uploads a directory of music using a presigned url",
						Action: func(cCtx *cli.Context) error {
							path := cCtx.Args().Get(0)
							handler.UploadMusic(path)
							return nil
						},
					},
					{
						Name:        "play",
						Usage:       "play",
						Description: "Prompts for a song choice and streams it",
						Action: func(cCtx *cli.Context) error {
							handler.StreamSong()
							return nil
						},
					},
					{
						Name:        "download",
						Usage:       "download>",
						Description: "Prompts for a song and downloads it",
						Action: func(cCtx *cli.Context) error {
							handler.DownloadSong()
							return nil
						},
					},
					{
						Name:        "ls",
						Usage:       "ls",
						Description: "Lists songs available to download/atream",
						Action: func(cCtx *cli.Context) error {
							cd := 1 //   cCtx.Args().Get(0)
							handler.ListFiles(cd)
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
