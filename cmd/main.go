package main

import (
	"flag"
	"fmt"
	"localhost/srm/command"
	"os"
)

var usage = `Usage: srm command [options]

srm - Safe Remove it is a simple tool to remove file/directory safety.

Option:

Commands:
	srm Remove a file/diretory using safe mode thats preserve file that is possible to restore
	frm Remove a file/diretory using without preserve file, this options it is unable to restore 
	ver Prints version info to console
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprint(usage))
	}

	var cmd *command.Command

	switch os.Args[1] {
	// case "srm":
	// 	cmd = command.SafeRemoveCommand()
	// case "frm":
	// 	cmd = command.ForceRemoveCommand()
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
		fmt.Fprint(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n")
	}
	flag.Usage()
	os.Exit(0)
}
