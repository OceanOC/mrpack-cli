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

	"github.com/fatih/color"
)

var downloaded = 1

type ModPack struct {
	Name          string       `json:"name"`
	Game          string       `json:"game"`
	VersionID     string       `json:"versionId"`
	FormatVersion int          `json:"formatVersion"`
	Dependencies  Dependencies `json:"dependencies"`
	Files         []Files      `json:"files"`

	// Custom stuff
	DownloadOptional bool
}

type Dependencies struct {
	Fabric    string `json:"fabric-loader"`
	Quilt     string `json:"quilt-loader"`
	Minecraft string `json:"minecraft"`
	NeoForge  string `json:"neoforge"`
	Forge     string `json:"forge"`
}

type Files struct {
	Path      string      `json:"path"`
	Downloads []string    `json:"downloads"`
	Hashes    []string    `json:"hashes"`
	Env       Environment `json:"env"`
}

type Environment struct {
	Client string `json:"client"`
	Server string `json:"server"`
}

func downloadMods(packFolder string, mrpac ModPack) {
	for i := range mrpac.Files {
		modMap := mrpac.Files[i]

		path := modMap.Path

		if strings.Contains(path, "mods") {
			if !mrpac.DownloadOptional {
				if modMap.Env.Client == "optional" && modMap.Env.Server != "required" {
					color.Set(color.FgYellow)
					fmt.Println("Skipping optional client mod: ", strings.Split(path, "/")[1])
					color.Unset()
					downloaded++
					continue
				} else if modMap.Env.Server == "optional" && modMap.Env.Client != "required" {
					color.Set(color.FgYellow)
					fmt.Println("Skipping optional server mod: ", strings.Split(path, "/")[1])
					color.Unset()
					downloaded++
					continue
				} else if modMap.Env.Server == "optional" && modMap.Env.Client == "optional" {
					color.Set(color.FgYellow)
					fmt.Println("Skipping optional mod: ", strings.Split(path, "/")[1])
					color.Unset()
					downloaded++
					continue
				}
			}
		}

		color.Set(color.FgGreen)
		fmt.Print("Downloading mod: ")
		color.Set(color.Bold)
		fmt.Print(strings.Split(path, "/")[1])
		color.Set(color.ResetBold)
		fmt.Println(" (" + strconv.FormatInt(int64(downloaded), 10) + "/" + strconv.FormatInt(int64(len(mrpac.Files)), 10) + ")")
		color.Unset()
		out, err := os.Create(packFolder + "mods/" + strings.Split(path, "/")[1])
		if err != nil {
			color.Set(color.FgRed)
			fmt.Println("Skipping...")
			fmt.Println("ERROR: Could not make mod:", err)
			color.Unset()
			downloaded++
			continue
		}
		defer out.Close()
		for i := range modMap.Downloads {
			dwn := modMap.Downloads[i]

			resp, err := http.Get(dwn)
			if err != nil {
				color.Set(color.FgRed)
				fmt.Println("ERROR: Could not download mod:", err)
				fmt.Println("Skipping...")
				color.Unset()
				downloaded++
				continue
			}
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				color.Set(color.FgRed)
				fmt.Println("ERROR: Could not copy mod data:", err)
				fmt.Println("Skipping...")
				color.Unset()
				downloaded++
				continue
			}
			defer resp.Body.Close()

			readSHA256(packFolder, path, modMap, "mods/")
		}
		downloaded++
	}
}

func downloadShaderPacks(packFolder string, mrpac ModPack) {
	for i := range mrpac.Files {
		modMap := mrpac.Files[i]

		path := modMap.Path

		if strings.Contains(path, "shaderpacks") {
			color.Set(color.FgYellow)
			fmt.Print("Downloading shaderpack: ")
			color.Set(color.Bold)
			fmt.Print(strings.Split(path, "/")[1])
			color.Set(color.ResetBold)
			fmt.Println("(" + strconv.FormatInt(int64(downloaded), 10) + "/" + strconv.FormatInt(int64(len(mrpac.Files)), 10) + ")")
			color.Unset()
			out, err := os.Create(packFolder + "shaderpacks/" + strings.Split(path, "/")[1])
			if err != nil {
				color.Set(color.FgRed, color.Bold)
				fmt.Println("Could not make shaderpacks folder:", err)
				fmt.Println("Skipping...")
				downloaded++
				continue
			}
			defer out.Close()
			for i := range modMap.Downloads {

				resp, err := http.Get(modMap.Downloads[i])
				if err != nil {
					color.Set(color.FgRed)
					fmt.Println("ERROR: Could not download shaderpack:", err)
					fmt.Println("Skipping...")
					color.Unset()
					downloaded++
					continue
				}

				_, err = io.Copy(out, resp.Body)
				if err != nil {
					color.Set(color.FgRed)
					fmt.Println("ERROR: Could not copy shaderpack data:", err)
					fmt.Println("Skipping...")
					color.Unset()
					downloaded++
					continue
				}
				defer resp.Body.Close()

				readSHA256(packFolder, path, modMap, "shaderpacks/")
			}
		}
	}
}

func downloadResourcePacks(packFolder string, mrpac ModPack) {
	for i := range mrpac.Files {
		modMap := mrpac.Files[i]

		path := modMap.Path

		if strings.Contains(path, "resourcepack") {
			color.Set(color.FgBlue)
			fmt.Print("Downloading resourcepack:")
			color.Set(color.Bold)
			fmt.Print(strings.Split(path, "/")[1])
			color.Set(color.ResetBold)
			fmt.Println("(" + strconv.FormatInt(int64(downloaded), 10) + "/" + strconv.FormatInt(int64(len(mrpac.Files)), 10) + ")")
			color.Unset()
			out, err := os.Create(packFolder + "resourcepacks/" + strings.Split(path, "/")[1])
			if err != nil {
				color.Set(color.FgRed)
				fmt.Println("ERROR: Could not download resourcepack:", err)
				fmt.Println("Skipping...")
				color.Unset()
				downloaded++
				continue
			}
			defer out.Close()
			for i := range modMap.Downloads {
				//fmt.Println(dwn.(string))
				resp, err := http.Get(modMap.Downloads[i])
				if err != nil {
					color.Set(color.FgRed)
					fmt.Println("ERROR: Could not download resourcepack:", err)
					fmt.Println("Skipping...")
					color.Unset()
					downloaded++
					continue
				}
				_, err = io.Copy(out, resp.Body)
				if err != nil {
					color.Set(color.FgRed)
					fmt.Println("ERROR: Could not copy resourcepack data:", err)
					fmt.Println("Skipping...")
					color.Unset()
					downloaded++
					continue
				}
				defer resp.Body.Close()

				readSHA256(packFolder, path, modMap, "resourcepacks/")
			}
			downloaded++
		}
	}
}

func addOverrides(packFolder string, tempFolder string) {
	fmt.Println("Copy: "+tempFolder+"overrides", "->", packFolder)
	switch runtime.GOOS {
	case "linux":
		cmd, err := exec.Command("/bin/sh", "-c", "cp -r "+tempFolder+"overrides/* "+packFolder).Output()
		if err != nil {
			color.Set(color.FgRed)
			fmt.Println("ERROR: Could not copy overrides:", err)
			color.Unset()
		}
		_ = cmd
	case "windows":
		cmd, err := exec.Command("robocopy", tempFolder+"overrides", packFolder, "/s").Output()
		if err != nil && err.Error() != "exit status 3" {
			color.Set(color.FgRed)
			fmt.Println("ERROR: Could not copy overrides:", err)
			color.Unset()
		}
		_ = cmd
	}
}

func readSHA256(packFolder, path string, modMap Files, ft string) error {
	filetype := ft
	f, err := os.Open(packFolder + filetype + strings.Split(path, "/")[1])
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println("ERROR: Could not open file:", err)
		color.Unset()
		return err
	}

	has := sha512.New()
	if _, err := io.Copy(has, f); err != nil {
		color.Set(color.FgRed)
		fmt.Println("ERROR: Could not copy file:", err)
		color.Unset()
		return err
	}
	f.Close()

	for i := range modMap.Hashes {
		fhas := modMap.Hashes[i]
		if fhas != hex.EncodeToString(has.Sum(nil)) {
			color.Set(color.FgRed, color.Bold)
			fmt.Println("Warning: Potentially Modified File")
			fmt.Println("The file hash doesn't match what is recorded in the .mrpack, which may indicate a fake or modified version. Please verify the file’s source and ensure it’s from a trusted provider. (e.g., Modrinth)")
			color.Unset()
		}
	}
	return nil
}
