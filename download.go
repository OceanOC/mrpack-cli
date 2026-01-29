package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func (mp ModPack) DownloadMods(outFolder string) {
	err := os.MkdirAll(filepath.Join(outFolder, "mods"), 0755)
	if err != nil {
		log.Fatalln(err)
	}

	c := color.New(color.FgHiYellow).Add(color.Bold)
	for i, mod := range mp.Files {
		if !strings.Contains(mod.Path, "mods") {
			continue
		}
		color.Green("downloading mod '" + filepath.Base(mod.Path) + "' (" + strconv.FormatInt(int64(i), 10) + "/" + strconv.FormatInt(int64(len(mp.Files)), 10) + ")")

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
				c.Println("could not download mod '" + down + "'")
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

func (mp ModPack) DownloadResourcePacks(outFolder string) {
	err := os.MkdirAll(filepath.Join(outFolder, "resourcepacks"), 0755)
	if err != nil {
		log.Fatalln(err)
	}

	c := color.New(color.FgHiYellow).Add(color.Bold)
	for i, mod := range mp.Files {
		if !strings.Contains(mod.Path, "resourcepacks") {
			continue
		}
		color.HiGreen("downloading resourcepack '" + filepath.Base(mod.Path) + "' (" + strconv.FormatInt(int64(i), 10) + "/" + strconv.FormatInt(int64(len(mp.Files)), 10) + ")")

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
				c.Println("could not download resourcepack '" + down + "'")
				continue
			}

			out, err := os.Create(filepath.Join(outFolder, mod.Path))
			if err != nil {
				c.Println("could not create file '" + filepath.Join(outFolder, mod.Path) + "', (" + err.Error() + ")")
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

func (mp ModPack) DownloadShaders(outFolder string) {
	err := os.MkdirAll(filepath.Join(outFolder, "shaderpacks"), 0755)
	if err != nil {
		log.Fatalln(err)
	}

	c := color.New(color.FgHiYellow).Add(color.Bold)
	for i, mod := range mp.Files {
		if !strings.Contains(mod.Path, "shaderpacks") {
			continue
		}
		color.HiGreen("downloading shader '" + filepath.Base(mod.Path) + "' (" + strconv.FormatInt(int64(i), 10) + "/" + strconv.FormatInt(int64(len(mp.Files)), 10) + ")")

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
				c.Println("could not download shader '" + down + "'")
				continue
			}

			out, err := os.Create(filepath.Join(outFolder, mod.Path))
			if err != nil {
				c.Println("could not create file '" + filepath.Join(outFolder, mod.Path) + "', (" + err.Error() + ")")
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
