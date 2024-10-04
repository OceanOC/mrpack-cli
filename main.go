package main

import (
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "mrpack-cli",
		Usage: "make and extract .mrpacks",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "add-entry",
				Aliases: []string{"entry"},
				Value:   false,
				Usage:   "If true add an entry of the modpack to the Minecraft launcher",
			},
		},
		Action: func(cCtx *cli.Context) error {
			tempd, err := createTempFolder()
			if err != nil {
				log.Fatal(err)
			}
			mpack, err := extract(tempd, cCtx.Args().Get(0))
			if err != nil {
				log.Fatal(err)
			}

			download(strings.ToLower(strings.ReplaceAll(mpack.Name, " ", "-"))+"/", mpack)

			err = copyOverrides(tempd, strings.ToLower(strings.ReplaceAll(mpack.Name, " ", "-"))+"/")
			if err != nil {
				log.Fatal(err)
			}

			err = os.RemoveAll(tempd)
			if err != nil {
				log.Fatal(err)
			}

			loader := -1
			if mpack.Dependencies.Fabric != "" {
				loader = 1
			} else if mpack.Dependencies.Neoforge != "" {
				loader = 0
			} else if mpack.Dependencies.Forge != "" {
				loader = 2
			} else if mpack.Dependencies.Quilt != "" {
				loader = 3
			}

			if cCtx.Bool("add-entry") {
				err = addEntry(strings.ToLower(strings.ReplaceAll(mpack.Name, " ", "-")), loader, mpack)
				if err != nil {
					log.Fatal(err)
				}
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
