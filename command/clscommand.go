package command

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

var days = 18

var clsUsage = `Purge files from safety storage after "18" days if no parameter it is informed.

Usage: srm cls 
	     srm cls 30

Options:
`

func CleanupCommand(srmHomeDir string) *Command {
	cmd := &Command{
		flags: flag.NewFlagSet("cls", flag.ExitOnError),
		Execute: func(cmd *Command, args []string) {
			if len(args) == 0 {
				fmt.Fprintf(os.Stdout, "Executing cleanup with default value: %v days\n", days)
			} else {
				if _, err := strconv.ParseInt(args[0], 10, 64); err != nil {
					fmt.Fprintf(os.Stderr, "Parameter must be a number.")
					os.Exit(-1)
				}
				days, _ = strconv.Atoi(args[0])
				fmt.Fprintf(os.Stdout, "Executing cleanup for files deleted bigger than: %v days\n", days)
			}

			srmCleanup(srmHomeDir)
		},
	}

	cmd.flags.Usage = func() {
		fmt.Fprintln(os.Stderr, clsUsage)
	}

	return cmd
}

func isOlderThan(t time.Time) bool {
	return time.Now().Sub(t) > time.Duration(days)*24*time.Hour
}

func srmCleanup(srmHome string) error {
	tmpfiles, err := ioutil.ReadDir(srmHome)
	if err != nil {
		return err
	}

	for _, file := range tmpfiles {
		if file.Mode().IsRegular() {
			if isOlderThan(file.ModTime()) {
				// srmfile.SrmRemove(file.Name())
				fmt.Println("Purging... > ", file.Name())
			}
		}
	}
	return nil
}
