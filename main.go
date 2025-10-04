package main

import (
	"context"
	"log"
	"os"

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
			mpack.outputDir = wd

			mpack.OpenMRPacks()
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
