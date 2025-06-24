package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func main() {
	if len(os.Args) == 1 {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("ERROR: Not enough arguments.")
		color.Set(color.FgCyan, color.Bold)
		fmt.Println("Hint: Try dragging and dropping the .mrpack file into the executable, press ENTER to exit.")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		color.Unset()
		os.Exit(1)
	}

	mrpack := os.Args[1]

	downPtr := flag.Bool("download", true, "Set to false to skip downloads")
	skipoptPtr := flag.Bool("skipopt", false, "Set to true to skip optional downloads")
	entryPtr := flag.Bool("entry", true, "Set to false to skip making entry in the Minecraft launcher")
	outPtr := flag.String("output", "default", "Set where the modpack will be extracted")

	flag.Parse()

	fs := flag.NewFlagSet("flags", flag.ExitOnError)

	fs.StringVar(outPtr, "output", "default", "Set where the modpack will be extracted")
	fs.BoolVar(downPtr, "download", true, "Set to false to skip downloads")
	fs.BoolVar(skipoptPtr, "skipopt", true, "Set to true to skip optional downloads")
	fs.BoolVar(entryPtr, "entry", true, "Set to false to skip making entry in the Minecraft launcher")

	fs.Parse(os.Args[2:])

	fmt.Println("mrpack-cli 1.2.0")

	var tempfolder = "mrpack-cli-" + strconv.FormatInt(rand.Int64N(99999), 10) + "/"

	switch runtime.GOOS {
	case "windows":
		tempfolder = strings.ReplaceAll(tempfolder, "/", "\\")
		tempfolder = os.Getenv("APPDATA") + "\\" + tempfolder
	case "linux":
		tempfolder = "/tmp/" + tempfolder
	}

	os.MkdirAll(tempfolder, 0700)
	unzip(mrpack, tempfolder)
	var jsonf = openjson(tempfolder + "modrinth.index.json")

	if jsonf.Game != "minecraft" {
		color.Set(color.FgRed)
		fmt.Println("ERROR: Game not supported.")
		color.Unset()
		os.Exit(1)
	}
	if jsonf.FormatVersion != 1 {
		color.Set(color.FgYellow)
		fmt.Println("WARNING: formatVersion not '1'.")
		color.Unset()
	}

	exePath, err := os.Executable()
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println("ERROR: getting executable path: ", err)
		color.Unset()
		os.Exit(1)
	}

	exePath, err = filepath.Abs(exePath)
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println("ERROR: getting absolute path: ", err)
		color.Unset()
		os.Exit(1)
	}

	var packFolder = ""

	if *outPtr == "default" {
		packFolder = filepath.Dir(exePath) + "/" + strings.ToLower(strings.ReplaceAll(jsonf.Name, " ", "-")+"/")
	} else {
		if !strings.HasSuffix(*outPtr, "\\") || !strings.HasSuffix(*outPtr, "/") {
			if runtime.GOOS == "windows" {
				packFolder = *outPtr + "\\"
			} else {
				packFolder = *outPtr + "/"
			}
		} else {
			packFolder = *outPtr
		}

		_, err = os.Stat(packFolder)
		if os.IsNotExist(err) {
			os.MkdirAll(packFolder, 0755)
		}
	}

	os.MkdirAll(packFolder+"mods/", os.ModePerm)
	os.MkdirAll(packFolder+"resourcepacks/", os.ModePerm)

	color.Set(color.FgGreen)
	fmt.Println("The modpack will be downloaded to: '" + packFolder + "'")
	color.Unset()

	if *skipoptPtr {
		jsonf.DownloadOptional = true
	}

	if *downPtr {
		downloadMods(packFolder, jsonf)
		downloadResourcePacks(packFolder, jsonf)
		downloadShaderPacks(packFolder, jsonf)
	}

	addOverrides(packFolder, tempfolder)

	if *entryPtr {
		if jsonf.Dependencies.Fabric != "" {
			installfabric(tempfolder, jsonf.Dependencies.Minecraft, jsonf.Dependencies.Fabric)
			addEntry(packFolder, strings.ToLower(strings.ReplaceAll(jsonf.Name, " ", "-")), jsonf.Name+" "+jsonf.VersionID, jsonf.Dependencies.Minecraft, jsonf.Dependencies.Fabric, "fabric-loader")
		}
		if jsonf.Dependencies.NeoForge != "" {
			installNeoforge(tempfolder, jsonf.Dependencies.NeoForge)
			addEntry(packFolder, strings.ToLower(strings.ReplaceAll(jsonf.Name, " ", "-")), jsonf.Name+" "+jsonf.VersionID, jsonf.Dependencies.Minecraft, jsonf.Dependencies.NeoForge, "neoforge")
		}
		if jsonf.Dependencies.Forge != "" {
			var launcherfolder string
			userhome, err := os.UserHomeDir()
			if err != nil {
				color.Set(color.FgRed)
				fmt.Println("ERROR: Could not get home directory")
				color.Unset()
				os.Exit(1)
			}
			switch runtime.GOOS {
			case "windows":
				launcherfolder = os.Getenv("APPDATA") + "\\.minecraft\\"
			case "linux":
				launcherfolder = userhome + "/.minecraft/"
			}

			installForge(tempfolder, launcherfolder, jsonf.Dependencies.Minecraft, jsonf.Dependencies.Forge)
			addEntry(packFolder, strings.ToLower(strings.ReplaceAll(jsonf.Name, " ", "-")), jsonf.Name+" "+jsonf.VersionID, jsonf.Dependencies.Minecraft, jsonf.Dependencies.Forge, "forge")
		}
		if jsonf.Dependencies.Quilt != "" {
			installQuilt(tempfolder, jsonf.Dependencies.Minecraft, jsonf.Dependencies.Quilt)
			addEntry(packFolder, strings.ToLower(strings.ReplaceAll(jsonf.Name, " ", "-")), jsonf.Name+" "+jsonf.VersionID, jsonf.Dependencies.Minecraft, jsonf.Dependencies.Quilt, "quilt-loader")
		}
	}

	os.RemoveAll(tempfolder)
}
