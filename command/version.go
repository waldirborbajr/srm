package command

import (
	"flag"
	"fmt"
	"os"
)

var versionUsage = `Print the app version and build info for the current context.

Usage: srm version [options]

Options:
  --short  If true, print just the version number. Default false.
`

var (
	build   = "???"
	version = "???"
	short   = false
)

var versionFunc = func(cmd *Command, args []string) {
	if short {
		fmt.Printf("srm via  v%s/2023", version)
	} else {
		fmt.Printf("srm via  v%s/2023, build: %s", version, build)
	}
	os.Exit(0)
}

func VersionCommand() *Command {
	cmd := &Command{
		flags:   flag.NewFlagSet("ver", flag.ExitOnError),
		Execute: versionFunc,
	}

	cmd.flags.BoolVar(&short, "short", false, "")
	cmd.flags.BoolVar(&short, "s", false, "")

	cmd.flags.Usage = func() {
		fmt.Fprintln(os.Stderr, versionUsage)
	}

	return cmd
}
