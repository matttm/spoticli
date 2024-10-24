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
						Name:  "ls",
						Usage: "ls",
						Action: func(cCtx *cli.Context) error {
							fmt.Println("Feature not implemented")
							return nil
						},
					},
					{
						Name:  "play",
						Usage: "play <song-title>",
						Action: func(cCtx *cli.Context) error {
							handler.StreamSong("1")
							return nil
						},
					},
					{
						Name:  "download",
						Usage: "download <song-title>",
						Action: func(cCtx *cli.Context) error {
							handler.DownloadSong("1")
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
