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

func downloadMods(packFolder string, jsonf map[string]interface{}) {
	if mods, ok := jsonf["files"].([]interface{}); ok {
		for _, mod := range mods {
			if modMap, ok := mod.(map[string]interface{}); ok {
				path := modMap["path"].(string)

				if strings.Contains(path, "mods") {
					color.Set(color.FgGreen)
					fmt.Print("Downloading mod: ")
					color.Set(color.Bold)
					fmt.Print(strings.Split(path, "/")[1])
					color.Set(color.ResetBold)
					fmt.Println(" (" + strconv.FormatInt(int64(downloaded), 10) + "/" + strconv.FormatInt(int64(len(mods)), 10) + ")")
					color.Unset()
					out, err := os.Create(packFolder + "mods/" + strings.Split(path, "/")[1])
					if err != nil {
						panic(err)
					}
					defer out.Close()
					for _, dwn := range modMap["downloads"].([]interface{}) {
						//fmt.Println(dwn.(string))
						resp, err := http.Get(dwn.(string))
						if err != nil {
							color.Set(color.FgRed)
							fmt.Println("ERROR: Could not download mod:", err)
							color.Unset()
							break
						}
						n, err := io.Copy(out, resp.Body)
						if err != nil {
							color.Set(color.FgRed)
							fmt.Println("ERROR: Could not copy mod data:", err)
							color.Unset()
							break
						}
						defer resp.Body.Close()
						_ = n // Just to make the compiler shut up

						readSHA256(packFolder, path, modMap, true)
					}
					downloaded++
				}
			}
		}
	} else {
		fmt.Println("Expected 'files' to be a slice, but got something else.")
	}
}

func downloadResourcePacks(packFolder string, jsonf map[string]interface{}) {
	if mods, ok := jsonf["files"].([]interface{}); ok {
		for _, mod := range mods {
			if modMap, ok := mod.(map[string]interface{}); ok {
				path := modMap["path"].(string)

				if strings.Contains(path, "resourcepack") {
					color.Set(color.FgBlue)
					fmt.Print("Downloading resourcepack:")
					color.Set(color.Bold)
					fmt.Print(strings.Split(path, "/")[1])
					color.Set(color.ResetBold)
					fmt.Println("(" + strconv.FormatInt(int64(downloaded), 10) + "/" + strconv.FormatInt(int64(len(mods)), 10) + ")")
					color.Unset()
					out, err := os.Create(packFolder + "resourcepacks/" + strings.Split(path, "/")[1])
					if err != nil {
						panic(err)
					}
					defer out.Close()
					for _, dwn := range modMap["downloads"].([]interface{}) {
						//fmt.Println(dwn.(string))
						resp, err := http.Get(dwn.(string))
						if err != nil {
							color.Set(color.FgRed)
							fmt.Println("ERROR: Could not download resourcepack:", err)
							color.Unset()
							break
						}
						n, err := io.Copy(out, resp.Body)
						if err != nil {
							color.Set(color.FgRed)
							fmt.Println("ERROR: Could not copy resourcepack:", err)
							color.Unset()
							break
						}
						defer resp.Body.Close()
						_ = n // Just to make the compiler shut up

						readSHA256(packFolder, path, modMap, false)
					}
					downloaded++
				}
			}
		}
	} else {
		fmt.Println("Expected 'files' to be a slice, but got ")
		os.Exit(1)
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

func readSHA256(packFolder, path string, modMap map[string]interface{}, ft bool) error {
	var filetype string
	if ft {
		filetype = "mods/"
	} else {
		filetype = "resourcepacks/"
	}
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

	if hashes, ok := modMap["hashes"].(map[string]interface{}); ok {
		fhas := hashes["sha512"].(string)
		if fhas != hex.EncodeToString(has.Sum(nil)) {
			color.Set(color.FgRed)
			fmt.Println("Warning: Potentially Modified File")
			fmt.Println("The file hash doesn't match what is recorded in the .mrpack, which may indicate a fake or modified version. Please verify the file’s source and ensure it’s from a trusted provider. (e.g., Modrinth)")
			fmt.Println("File deleted.")
			color.Unset()
			os.Remove(packFolder + filetype + strings.Split(path, "/")[1])
		}
	}
	return nil
}
