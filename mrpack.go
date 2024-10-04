package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type Mrpack struct {
	FormatVersion int          `json:"formatVersion"`
	Game          string       `json:"game"`
	VersionID     string       `json:"versionId"`
	Name          string       `json:"name"`
	Files         []File       `json:"files"`
	Dependencies  Dependencies `json:"dependencies"`
}

type File struct {
	Path      string      `json:"path"`
	Hashes    Hash        `json:"hashes"`
	Env       Environment `json:"env"`
	Downloads []string    `json:"downloads"`
}

type Environment struct {
	Client string `json:"client"`
	Server string `json:"server"`
}

type Hash struct {
	SHA512 string `json:"sha512"`
}

type Dependencies struct {
	Minecraft string `json:"minecraft"`
	Fabric    string `json:"fabric-loader"`
	Quilt     string `json:"quilt-loader"`
	Forge     string `json:"forge"`
	Neoforge  string `json:"neoforge"`
}

func download(packFolder string, mrp Mrpack) error {
	downloaded := 1
	for i := range mrp.Files {
		mod := mrp.Files[i]

		fmt.Printf("Downloading (%v/%v): ", strconv.FormatInt(int64(downloaded), 10), strconv.FormatInt(int64(len(mrp.Files)), 10))
		fmt.Println(strings.Split(mod.Path, "/")[1])

		err := os.MkdirAll(packFolder+strings.Split(mod.Path, "/")[0], 0777)
		if err != nil {
			fmt.Println("Skipping...")
			fmt.Println("ERROR: Could not make mod folder:", err)
			downloaded++
			continue
		}

		out, err := os.Create(packFolder + mod.Path)
		if err != nil {
			fmt.Println("Skipping...")
			fmt.Println("ERROR: Could not make mod:", err)
			downloaded++
			continue
		}
		for i := range mod.Downloads {
			dwn := mod.Downloads[i]

			resp, err := http.Get(dwn)
			if err != nil {
				fmt.Println("ERROR: Could not download mod:", err)
				fmt.Println("Skipping...")
				downloaded++
				continue
			}
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				fmt.Println("ERROR: Could not copy mod data:", err)
				fmt.Println("Skipping...")
				downloaded++
				continue
			}
			defer resp.Body.Close()
			defer out.Close()

			o, err := os.Open(packFolder + mod.Path)
			if err != nil {
				fmt.Println("Skipping...")
				fmt.Println("ERROR: Could not make mod:", err)
				downloaded++
				continue
			}

			has := sha512.New()
			if _, err := io.Copy(has, o); err != nil {
				fmt.Println("ERROR: Could not copy file:", err)
				return err
			}

			fhas := mod.Hashes.SHA512
			if fhas != hex.EncodeToString(has.Sum(nil)) {
				fmt.Println("Warning: Potentially Modified File")
				fmt.Println("The file hash doesn't match what is recorded in the .mrpack, which may indicate a fake or modified version. Please verify the file’s source and ensure it’s from a trusted provider. (e.g., Modrinth)")
			}
		}
		downloaded++

	}
	return nil
}

func copyOverrides(tempd, packFolder string) error {
	fmt.Println("Copy: "+tempd+"overrides", "->", packFolder)
	switch runtime.GOOS {
	case "linux":
		cmd, err := exec.Command("/bin/sh", "-c", "cp -r "+tempd+"overrides/* "+packFolder).Output()
		if err != nil {
			return err
		}
		_ = cmd
	case "windows":
		cmd, err := exec.Command("robocopy", tempd+"overrides", packFolder, "/s").Output()
		if err != nil && err.Error() != "exit status 3" {
			return err
		}
		_ = cmd
	}
	return nil
}
