package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"
)

func (mp ModPack) DownloadFiles(outFolder string) {
	err := os.MkdirAll(filepath.Join(outFolder, "mods"), 0755)
	if err != nil {
		log.Fatalln(err)
	}

	c := color.New(color.FgHiYellow).Add(color.Bold)
	for i, mod := range mp.Files {
		err := os.MkdirAll(filepath.Join(outFolder, filepath.Dir(mod.Path)), 0755)
		if err != nil {
			log.Fatalln(err)
		}

		color.HiGreen("downloading file '" + filepath.Base(mod.Path) + "' [" + strconv.FormatFloat(float64(mod.FileSize)/1000000, 'f', 2, 64) + " MB] (" + strconv.FormatInt(int64(i), 10) + "/" + strconv.FormatInt(int64(len(mp.Files)), 10) + ")")

		for _, down := range mod.Downloads {
			furl, err := url.Parse(down)
			if err != nil {
				c.Println("could not parse URL '" + down + "'")
				continue
			}
			if mp.NotInWhitelist(furl.Host) {
				c.Println("download link is not from modrinth '" + down + "'")
				continue
			}

			resp, err := http.Get(down)
			if err != nil {
				c.Println("could not download file '" + down + "'")
				continue
			}

			out, err := os.Create(filepath.Join(outFolder, mod.Path))
			if err != nil {
				c.Println("could not create file '" + filepath.Join(outFolder, mod.Path) + "'")
				continue
			}
			defer out.Close()

			_, err = io.Copy(out, resp.Body)
			if err != nil {
				c.Println("could not write to file '" + filepath.Join(outFolder, mod.Path) + "'")
				continue
			}
		}
	}
}
