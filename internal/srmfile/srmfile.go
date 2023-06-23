package srmfile

import (
	"os"
	"strings"
)

func SrmRemove(srmRmFileName string) error {
	if err := os.Remove(srmRmFileName); err != nil {
		return err
	}
	return nil
}

func SrmRemoveDirectory(srmRmFileName string) error {
	if err := os.RemoveAll(srmRmFileName); err != nil {
		return err
	}
	return nil
}

func SrmCleanup(srmRmFileName string) error {
	// Remove .zlib file
	if err := os.Remove(srmRmFileName); err != nil {
		return err
	}

	// Remove original file fom .srm folder
	if err := os.Remove(strings.TrimSuffix(srmRmFileName, ".zlib")); err != nil {
		return err
	}

	return nil
}
