package srmfile

import (
	"os"
)

func SrmHome() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homeDir + "/.srm", nil
}
