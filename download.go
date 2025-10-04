package main

import (
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func (mp ModPack) DownloadMods() {
	c := color.New(color.FgGreen)
	for i, mod := range mp.Files {
		if !strings.Contains(mod.Path, "mods") {
			continue
		}

		c.Println(" (" + strconv.FormatInt(int64(i), 10) + "/" + strconv.FormatInt(int64(len(mp.Files)), 10) + ")")

		// TODO: download mods
	}
}
