package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/fatih/color"
)

func addEntry(packfolder, packname, fancypackname, gamever, fabricver, loader string) error {
	var launcherfolder string

	userhome, err := os.UserHomeDir()
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println("ERROR: Could not get home directory")
		color.Unset()
		return err
	}
	switch runtime.GOOS {
	case "windows":
		launcherfolder = os.Getenv("APPDATA") + "\\.minecraft\\"
	case "linux":
		launcherfolder = userhome + "/.minecraft/"
	}

	color.Set(color.FgGreen)
	fmt.Println("Getting icon from Modrinth...")
	color.Unset()

	n, err := http.Get("https://api.modrinth.com/v2/project/" + packname)
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println("ERROR: Could not communicate with Modrinth API:", err)
		color.Unset()
		return err
	}

	api, err := io.ReadAll(n.Body)
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println("ERROR: Could not read Modrinth API response body:", err)
		color.Unset()
		return err
	}

	var iconURI string

	var js map[string]interface{} = openjsonfromstring(string(api))
	if uri, ok := js["icon_url"].(string); ok {
		n, err := http.Get(uri)
		if err != nil {
			color.Set(color.FgRed)
			fmt.Println("ERROR: Could not get icon:", err)
			color.Unset()
		}

		d, err := io.ReadAll(n.Body)
		if err != nil {
			color.Set(color.FgRed)
			fmt.Println("ERROR: Could not read icon:", err)
			color.Unset()
		}

		img := base64.StdEncoding.EncodeToString(d)

		iconURI = fmt.Sprintf("data:image/png;base64," + img)
	}

	color.Set(color.FgGreen)
	fmt.Println("Adding entry to Minecraft launcher")
	color.Unset()

	var ljs map[string]interface{} = openjson(launcherfolder + "launcher_profiles.json")
	if pf, ok := ljs["profiles"].(map[string]interface{}); ok {
		pf[packname] = map[string]interface{}{
			"name":          fancypackname,
			"type":          "custom",
			"created":       time.Now().Format(time.RFC3339),
			"lastUsed":      time.Time{},
			"icon":          iconURI,
			"gameDir":       packfolder,
			"lastVersionId": loader + "-" + fabricver + "-" + gamever,
		}
	}

	ujs, err := json.MarshalIndent(ljs, "", "  ")
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println("ERROR: marshalling JSON:", err)
		color.Unset()
		return err
	}
	_ = ujs

	err = os.WriteFile(launcherfolder+"launcher_profiles.json", ujs, 0664)
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println("ERROR: writing file:", err)
		color.Unset()
		return err
	}

	return nil
}
