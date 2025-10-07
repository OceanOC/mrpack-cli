package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
)

func main() {
	var mpack mrpcli
	mpack.modPacks = make(map[string]ModPack)
	cmd := &cli.Command{
		Name:                   "mrpack-cli",
		Usage:                  "download and add Modrinth modpacks to the vanilla launcher",
		UseShortOptionHandling: true,
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name:        "mrpackFile",
				Max:         -1,
				Destination: &mpack.files,
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "automated",
				Aliases:     []string{"auto"},
				Value:       false,
				Usage:       "removes 'Press ENTER to exit', useful for automated scripts",
				Destination: &mpack.automated,
			},
			&cli.BoolFlag{
				Name:        "no-color",
				Usage:       "removes color",
				Destination: &color.NoColor,
			},
			&cli.BoolFlag{
				Name:        "nodownload",
				Aliases:     []string{"d"},
				Usage:       "dont download mods",
				Destination: &mpack.nodownload,
				Value:       false,
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Usage:       "mrpack(s) download final location",
				DefaultText: "working directory",
				Value:       "",
				Destination: &mpack.outputDir,
			},
		},
		Action: func(context.Context, *cli.Command) error {
			wd, err := os.Getwd()
			if err != nil {
				mpack.ExitCLIWithError("Couldn't get working directory", err)
			}
			if mpack.outputDir == "" {
				mpack.outputDir = wd
			} else {
				mpack.outputDir = filepath.Clean(mpack.outputDir)
			}

			mpack.OpenMRPacks()

			if len(mpack.files) == 1 {
				mpack.ExitCLI(strconv.Itoa(len(mpack.files))+" modpack extracted and downloaded", Success)
			} else {
				mpack.ExitCLI(strconv.Itoa(len(mpack.files))+" modpacks extracted and downloaded", Success)
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
