package main

import (
	"flag"
	"fmt"
	"os"
	"waldirborbajr/srm/command"
	"waldirborbajr/srm/internal/app"
	"waldirborbajr/srm/internal/srmfile"
)

var usage = `Usage: srm command [options]

srm - Safe ReMove it is a simple tool to remove file/directory safety.

Option:

Commands:
	srm - Remove a file/diretory using safe mode thats preserve file that is possible to restore
	rst - Restore a file/diretory that was deleted with safe option
	cls - Cleanup removed files after 18 days if not informed another day as parameter
	hlp - Display this help information
	ver - Prints version info to console
`

func init() {
}

func main() {
	srmHomeDir, err := srmfile.SrmHome()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: unable to get home directory: ")
		os.Exit(-1)
	}
	app := &app.Srm{SrmHomeDir: srmHomeDir}

	// Verify if .srm safety store folder exists
	srmFolderExists(*app)

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprint(usage))
	}

	flag.Parse()

	if flag.NArg() < 1 {
		usageAndExit("")
	}

	var cmd *command.Command

	switch os.Args[1] {
	case "srm":
		cmd = command.NewSafeRemoveCommand(*app)
	case "rst":
		cmd = command.NewRestoreCommand(*app)
	case "cls":
		cmd = command.NewCleanupCommand(*app)
	case "hlp":
		flag.Usage()
		os.Exit(0)
	case "ver":
		cmd = command.VersionCommand()
	default:
		fmt.Fprintf(os.Stderr, "srm: '%s' is not a srm valid command.\n execute srm hlp for help\n", os.Args[1])
		os.Exit(-1)
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

func srmFolderExists(app app.Srm) {
	_, err := os.Stat(app.SrmHomeDir)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(app.SrmHomeDir, 0700); err != nil {
				fmt.Fprint(os.Stderr, "Error creating srm restore point folder.")
				os.Exit(-1)
			}
		}
	}
}
