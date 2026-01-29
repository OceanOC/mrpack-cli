package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
)

func (mp mrpcli) AddToLauncher(pack ModPack) {
	var home string
	switch runtime.GOOS {
	case "linux":
		homer, err := os.UserHomeDir()
		if err != nil {
			mp.ExitCLIWithError("cannot find home", err)
		}
		home = homer
	case "windows":
		homer, err := os.UserConfigDir()
		if err != nil {
			mp.ExitCLIWithError("cannot find %APPDATA%", err)
		}
		home = homer
	default:
		color.HiRed("unsupported OS, skipping adding to launcher")
		return
	}

	minecraftLauncher := filepath.Join(home, ".minecraft")

	info, err := os.Stat(minecraftLauncher)
	if err != nil {
		mp.ExitCLIWithError("cannot access .minecraft", err)
	}
	if !info.IsDir() {
		mp.ExitCLIWithError(".minecraft is not a directory", err)
	}

	var launcherJson string
	var storeLauncher bool
	info, err = os.Stat(filepath.Join(minecraftLauncher, "launcher_profiles.json"))
	if err != nil {
		storeLauncher = true
	}
	if storeLauncher {
		info, err = os.Stat(filepath.Join(minecraftLauncher, "launcher_profiles_microsoft_store.json"))
		if err != nil {
			mp.ExitCLIWithError("launcher_profiles doesnt exist", err)
		}
		launcherJson = filepath.Join(minecraftLauncher, "launcher_profiles_microsoft_store.json")
	} else {
		launcherJson = filepath.Join(minecraftLauncher, "launcher_profiles.json")
	}

	ljs, err := os.ReadFile(launcherJson)
	if err != nil {
		mp.ExitCLIWithError("cannot read launcher_profiles", err)
	}
	var launcher MineLauncher
	err = json.Unmarshal(ljs, &launcher)
	if err != nil {
		mp.ExitCLIWithError("cannot parse launcher_profiles", err)
	}

	// TODO: add the rest
}
