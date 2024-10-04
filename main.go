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
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
