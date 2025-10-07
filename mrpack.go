package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
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
				mpack.modpackDir = filepath.Join(mpack.outputDir, strings.ReplaceAll(strings.ToLower(mp.Name), " ", "-"))

				if !mpack.nodownload {
					mp.DownloadMods(mpack.modpackDir)
				}

				manifestFound = true
			}
		}

		if !manifestFound {
			mpack.ExitCLI("modrinth.index.json not found, not a mrpack", Error)
		}

		err = os.MkdirAll(mpack.modpackDir, 0755)
		if err != nil {
			mpack.ExitCLIWithError("cannot create directory", err)
		}

		// Looping through the zip file twice is not very efficient but i cannot get extraction folder before the manifest
		for _, rf := range r.File {
			if strings.HasPrefix(rf.Name, "overrides") && !rf.FileInfo().IsDir() {
				c := color.New(color.FgHiCyan).Add(color.Bold)
				c.Printf("Extracting override: %s\n", rf.Name)

				fo, err := rf.Open()
				if err != nil {
					mpack.ExitCLIWithError("cannot open override", err)
				}
				err = os.MkdirAll(filepath.Dir(filepath.Join(mpack.modpackDir, rf.Name)), os.ModePerm)
				if err != nil {
					mpack.ExitCLIWithError("cannot create directory", err)
				}
				file, err := os.Create(filepath.Join(mpack.modpackDir, rf.Name))
				if err != nil {
					mpack.ExitCLIWithError("cannot create file", err)
				}
				defer file.Close()

				_, err = io.Copy(file, fo)
				if err != nil {
					mpack.ExitCLIWithError("cannot write to file", err)
				}
			}
		}

		if !mpack.nolauncher {
			mpack.AddToLauncher(mp)
		}

		fmt.Println("File completed" + " '" + fp + "'" + " (" + strconv.FormatInt(int64(i+1), 10) + "/" + strconv.FormatInt(int64(len(mpack.files)), 10) + ")")
	}
}
