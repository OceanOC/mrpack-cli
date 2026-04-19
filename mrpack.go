package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"crypto/sha512"
	"encoding/hex"
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
					mp.DownloadFiles(mpack.modpackDir)
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
				oN := strings.Replace(rf.Name, "overrides/", "", 1)

				c := color.New(color.FgHiCyan).Add(color.Bold)
				c.Printf("Extracting override: %s\n", oN)

				fo, err := rf.Open()
				if err != nil {
					mpack.ExitCLIWithError("cannot open override", err)
				}
				err = os.MkdirAll(filepath.Dir(filepath.Join(mpack.modpackDir, oN)), os.ModePerm)
				if err != nil {
					mpack.ExitCLIWithError("cannot create directory", err)
				}
				file, err := os.Create(filepath.Join(mpack.modpackDir, oN))
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

		fmt.Println("Checking Mod hashes...")

		for _, mod := range mp.Files {
			modJar := filepath.Join(mpack.modpackDir, mod.Path)

			for sha, hash := range mod.Hashes {
				if sha == "sha1" {
					// sha1 should not be used to do file checksums
					continue
				}

				file, err := os.ReadFile(modJar)
				if err != nil {
					fmt.Println("could not read file \"" + modJar + "\"")
					continue
				}

				su := sha512.Sum512(file)
				sum := hex.EncodeToString(su[:])

				if sum != hash {
					color.HiRed("WARNING: mod \"" + modJar + "\" does not match mrpack hash")
					color.HiRed("This could mean that this file has been tampered with and may possibly contain malware. Please make sure you trust the mrpack author and have reviewed the mrpack contents before continuing.")

					if mpack.automated {
						mpack.ExitCLI("Mod files have been tampered with.", Error)
					}

					fmt.Println("FHASH: " + string(sum))
					fmt.Println("MRHASH: " + hash)

					fmt.Println("Please select [1-2] and press ENTER")
					fmt.Println("  1) Abort    [DEFAULT]")
					fmt.Println("  2) Continue")

					num, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
					if err != nil {
						mpack.ExitCLIWithError("could not get option", err)
					}

					if strings.Contains(string(num), "1") && !strings.Contains(string(num), "2") {
						mpack.ExitCLI("Mod files have been tampered with.", Error)
					} else if !strings.Contains(string(num), "1") && strings.Contains(string(num), "2") {
						color.HiRed("continuing...")
						continue
					} else {
						mpack.ExitCLI("Mod files have been tampered with.", Error)
					}
				} else {
					fmt.Println("mod \"" + mod.Path + "\" verified.")
				}
			}
		}

		if !mpack.nolauncher {
			mpack.AddToLauncher(mp)
		}

		fmt.Println("File completed" + " '" + fp + "'" + " (" + strconv.FormatInt(int64(i+1), 10) + "/" + strconv.FormatInt(int64(len(mpack.files)), 10) + ")")
	}
}
