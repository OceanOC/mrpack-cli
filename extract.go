package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

func extract(tempd, packFolder string) (Mrpack, error) {
	mp := Mrpack{}

	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("tar", "-xf", packFolder, "-C", tempd)
		cmd.Run()
	case "linux":
		if _, err := os.Stat("/bin/unzip"); err == nil {
			_, err := exec.Command("unzip", packFolder, "-d", tempd).Output()
			if err != nil {
				return mp, err
			}
		} else {
			fmt.Println("unzip not found.")
			os.Exit(2)
		}
	}

	f, err := os.Open(tempd + "modrinth.index.json")
	if err != nil {
		return mp, err
	}

	i, err := io.ReadAll(f)
	if err != nil {
		return mp, err
	}

	err = json.Unmarshal(i, &mp)
	if err != nil {
		return mp, err
	}

	return mp, nil
}
