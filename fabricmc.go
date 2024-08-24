package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
)

func installfabric(tempfolder, gameversion, loaderver string) error {
	var java_home = "java"
	if runtime.GOOS == "windows" {
		java_home = os.Getenv("JAVA_HOME") + "\\bin\\java.exe"
	}

	out, err := os.Create(tempfolder + "/fabric-installer.jar")
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("ERROR: Could not create fabric installer file")
		color.Unset()
		return err
	}

	color.Set(color.FgCyan)
	fmt.Println("Downloading Fabric installer...")
	color.Unset()

	resp, err := http.Get("https://maven.fabricmc.net/net/fabricmc/fabric-installer/1.0.1/fabric-installer-1.0.1.jar")
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("ERROR: Could not download fabric installer from " + resp.Request.URL.RequestURI())
		color.Unset()
		return err
	}

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("ERROR: Could not copy fabric installer")
		color.Unset()
		return err
	}
	_ = n

	color.Set(color.FgCyan, color.Bold)
	fmt.Println("Running Fabric installer... (This might take a while)")
	color.Unset()

	fi := exec.Command(java_home, "-jar", tempfolder+"fabric-installer.jar", "client", "-mcversion", gameversion, "-loader", loaderver, "-noprofile")

	output, err := fi.CombinedOutput()
	if err != nil {
		fmt.Printf("Fabric installer failed: %v\n", err)
	}

	fmt.Printf("%s\n", output)

	return nil
}
