package srmfile

import (
	"os"
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
	// if err := os.Remove(srmRmFileName); err != nil {
	// 	return err
	// }

	// Remove original file of .srm folder
	if err := os.Remove(srmRmFileName); err != nil {
		return err
	}

	return nil
}
