package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

type MCLauncher struct {
	Profiles map[string]Profile
}

type Modrinth struct {
	Icon string `json:"icon_url"`
}

type Profile struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	Created       string `json:"created"`
	LastUsed      string `json:"lastUsed"`
	Icon          string `json:"icon"`
	LastVersionID string `json:"lastVersionID"`
	GameDir       string `json:"gameDir"`
}

func addEntry(packFolder string, loader int, mrp Mrpack) error {
	var launcherfolder string
	switch runtime.GOOS {
	case "linux":
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		launcherfolder = dirname + "/.minecraft/"
	case "windows":
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		launcherfolder = dirname + "\\.minecraft\\"
	}

	i, err := os.ReadFile(launcherfolder + "launcher_profiles.json")
	if err != nil {
		return err
	}

	lnr := MCLauncher{}

	err = json.Unmarshal(i, &lnr)
	if err != nil {
		return err
	}

	n, err := http.Get("https://api.modrinth.com/v2/project/" + strings.ToLower(strings.ReplaceAll(mrp.Name, " ", "-")))
	if err != nil {
		return err
	}

	api, err := io.ReadAll(n.Body)
	if err != nil {
		return err
	}

	var iconURI string

	result := Modrinth{}
	json.Unmarshal([]byte(api), &result)

	n, err = http.Get(result.Icon)
	if err != nil {
		return err
	}

	d, err := io.ReadAll(n.Body)
	if err != nil {
		return err
	}

	img := base64.StdEncoding.EncodeToString(d)

	iconURI = fmt.Sprintf("data:image/png;base64," + img)

	var versionId string

	switch loader {
	case 0:
		// Neoforge
		versionId = "neoforge-" + mrp.Dependencies.Quilt
	case 1:
		// Fabric
		versionId = "fabric-loader-" + mrp.Dependencies.Fabric + "-" + mrp.Dependencies.Minecraft
	case 2:
		// Forge
		versionId = mrp.Dependencies.Minecraft + "-forge-" + mrp.Dependencies.Forge
	case 3:
		// Quilt
		versionId = "quilt-loader-" + mrp.Dependencies.Quilt + "-" + mrp.Dependencies.Minecraft
	}

	lnr = MCLauncher{
		Profiles: map[string]Profile{
			strings.ToLower(strings.ReplaceAll(mrp.Name, " ", "-")): {
				Name:          mrp.Name,
				Type:          "custom",
				Created:       time.Now().Format(time.RFC3339),
				LastUsed:      time.Time{}.String(),
				Icon:          iconURI,
				GameDir:       packFolder,
				LastVersionID: versionId,
			},
		},
	}

	nlnr, err := json.MarshalIndent(lnr, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(launcherfolder+"launcher_profiles.json", nlnr, 0664)
	if err != nil {
		return err
	}

	return nil
}
