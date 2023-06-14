package main

import (
	"flag"
	"fmt"
	"os"
	"waldirborbajr/srm/command"
	"waldirborbajr/srm/internal/srmfile"
)

var srmHomeDir string

var usage = `Usage: srm command [options]

srm - Safe Remove it is a simple tool to remove file/directory safety.

Option:

Commands:
	srm - Remove a file/diretory using safe mode thats preserve file that is possible to restore
	rst - Restore a file/diretory that was deleted with safe option
	cls - Cleanup removed files after 18 days if not informed another day as parameter
	ver - Prints version info to console
`

func init() {
	var err error

	srmHomeDir, err = srmfile.SrmHome()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: unable to get home directory: ")
		os.Exit(-1)
	}
}

func main() {
	// Verify if .srm exists
	srmFolderExists()

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprint(usage))
	}

	var cmd *command.Command

	switch os.Args[1] {
	case "srm":
		cmd = command.SafeRemoveCommand(srmHomeDir)
	case "rst":
		cmd = command.RestoreCommand()
	case "cls":
		cmd = command.CleanupCommand(srmHomeDir)
	case "ver":
		cmd = command.VersionCommand()
	default:
		usageAndExit(fmt.Sprintf("srm: '%s' is not a srm valid command.\n", os.Args[1]))
	}

	cmd.Init(os.Args[2:])
	cmd.Run()
}

func usageAndExit(msg string) {
	if msg != "" {
		if err := fmt.Errorf("%s", msg); err != nil {
			os.Exit(-1)
		}
		// fmt.Fprint(os.Stderr, msg)
		// fmt.Fprintf(os.Stderr, "\n")
	}
	flag.Usage()
	os.Exit(0)
}

func srmFolderExists() {
	_, err := os.Stat(srmHomeDir)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(srmHomeDir, 0700); err != nil {
				fmt.Fprint(os.Stderr, "Error creating srm restore point folder.")
				os.Exit(-1)
			}
		}
	}
}
