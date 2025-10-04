package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func (mpack mrpcli) OpenMRPacks() {
	for i, file := range mpack.files {
		fp := filepath.Clean(file)
		var mp ModPack

		r, err := zip.OpenReader(fp)
		if err != nil {
			mpack.ExitCLIWithError("couldn't open "+fp+"", err)
		}
		defer r.Close()

		manifestFound := false
		for _, rf := range r.File {
			if rf.FileInfo().Name() == "modrinth.index.json" {
				fo, err := rf.Open()
				if err != nil {
					mpack.ExitCLIWithError("couldn't open modrinth.index.json", err)
				}
				buf := new(bytes.Buffer)
				buf.ReadFrom(fo)

				json.Unmarshal(buf.Bytes(), &mp)

				mp.DownloadMods()

				manifestFound = true
			} else if strings.HasPrefix(rf.Name, "overrides/") {
				c := color.New(color.FgHiCyan).Add(color.Bold)
				c.Printf("Extracting override: %s\n", rf.Name)

				// TODO: extract overrides
			}
		}

		if !manifestFound {
			mpack.ExitCLI("modrinth.index.json not found, not a mrpack", Error)
		}

		fmt.Println("File completed" + " '" + fp + "'" + " (" + strconv.FormatInt(int64(i+1), 10) + "/" + strconv.FormatInt(int64(len(mpack.files)), 10) + ")")
		mpack.modPacks[fp] = mp
	}
}
