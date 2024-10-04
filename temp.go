package main

import (
	"math/rand/v2"
	"os"
	"strconv"
)

func createTempFolder() (string, error) {
	tempfolder := "mrpack-cli-" + strconv.FormatInt(rand.Int64N(99999), 10)

	tempdir, err := os.MkdirTemp("", tempfolder)
	if err != nil {
		return "", err
	}
	return tempdir + "/", nil
}
