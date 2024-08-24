package main

import (
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
		color.Unset()
		os.Exit(1)
	}

	mrpack := os.Args[1]

	downPtr := flag.Bool("download", true, "Set to false to skip downloads")
	entryPtr := flag.Bool("entry", true, "Set to false to skip making entry in the Minecraft launcher")
	outPtr := flag.String("output", "default", "Set where the modpack will be extracted")

	flag.Parse()

	fs := flag.NewFlagSet("flags", flag.ExitOnError)

	fs.StringVar(outPtr, "output", "default", "Set where the modpack will be extracted")
	fs.BoolVar(downPtr, "download", true, "Set to false to skip downloads")
	fs.BoolVar(entryPtr, "entry", true, "Set to false to skip making entry in the Minecraft launcher")

	fs.Parse(os.Args[2:])

	fmt.Println("mrpack-cli 1.0.0-alpha.2")

	var tempfolder = "mrpack-cli-" + strconv.FormatInt(rand.Int64N(99999), 10) + "/"

	if runtime.GOOS == "windows" {
		tempfolder = strings.ReplaceAll(tempfolder, "/", "\\")
		tempfolder = os.Getenv("APPDATA") + "\\" + tempfolder
	} else if runtime.GOOS == "linux" {
		tempfolder = "/tmp/" + tempfolder
	}
	os.MkdirAll(tempfolder, 0700)
	unzip(mrpack, tempfolder)
	//var jsonf
	var jsonf map[string]interface{} = openjson(tempfolder + "modrinth.index.json")

	if jsonf["game"] != "minecraft" {
		color.Set(color.FgRed)
		fmt.Println("ERROR: Game not supported.")
		color.Unset()
		os.Exit(1)
	}
	if jsonf["formatVersion"].(float64) != 1 {
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
		packFolder = filepath.Dir(exePath) + "/" + strings.ToLower(strings.ReplaceAll(jsonf["name"].(string), " ", "-")+"/")
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
			os.MkdirAll(packFolder, 755)
		}
	}

	os.MkdirAll(packFolder+"mods/", os.ModePerm)
	os.MkdirAll(packFolder+"resourcepacks/", os.ModePerm)

	color.Set(color.FgGreen)
	fmt.Println("The modpack will be downloaded to: '" + packFolder + "'")
	color.Unset()

	if *downPtr {
		downloadMods(packFolder, jsonf)
		downloadResourcePacks(packFolder, jsonf)
	}

	addOverrides(packFolder, tempfolder)

	if *entryPtr {
		if vern, ok := jsonf["dependencies"].(map[string]interface{}); ok {
			if vern["fabric-loader"] != nil {
				installfabric(tempfolder, vern["minecraft"].(string), vern["fabric-loader"].(string))
				addEntry(packFolder, strings.ToLower(strings.ReplaceAll(jsonf["name"].(string), " ", "-")), jsonf["name"].(string)+" "+jsonf["versionId"].(string), vern["minecraft"].(string), vern["fabric-loader"].(string), "fabric-loader")
			}
		}
	}

	os.RemoveAll(tempfolder)
}
