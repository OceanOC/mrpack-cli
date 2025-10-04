package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
)

type PrintLevel int

const (
	Success PrintLevel = iota
	Warning
	Error
)

func (mp mrpcli) ExitCLIWithError(msg string, err error) {
	c := color.New(color.FgHiRed, color.Bold)
	c.Println("ERROR: " + msg)
	fmt.Printf("%s\n", err)
	if !mp.automated {
		fmt.Println("Press ENTER to exit")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
	os.Exit(2)
}

func (mp mrpcli) ExitCLI(msg string, lvl PrintLevel) {
	switch lvl {
	case Success:
		c := color.New(color.FgHiGreen, color.Bold)
		c.Println(msg)
		if !mp.automated {
			fmt.Println("Press ENTER to exit")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
		os.Exit(0)
	case Warning:
		c := color.New(color.FgHiYellow, color.Bold)
		c.Println("WARNING: " + msg)
		if !mp.automated {
			fmt.Println("Press ENTER to exit")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
		os.Exit(1)
	case Error:
		c := color.New(color.FgHiRed, color.Bold)
		c.Println("ERROR: " + msg)
		if !mp.automated {
			fmt.Println("Press ENTER to exit")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
		os.Exit(2)
	default:
		panic(fmt.Sprintf("unexpected main.ExitLevel: %#v", lvl))
	}

}
