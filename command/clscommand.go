package command

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
	"waldirborbajr/srm/internal/app"
)

var days string

var clsUsage = `Purge files from safety storage after "18" days if no parameter it is informed.

Usage: srm cls [OPTIONS]

Options:
	-d, --day days to purge srm safer folder
`

func clsFunc(cmd *Command, args []string, app app.Srm) {
	if len(days) == 0 {
		fmt.Println("Default value.")
		fmt.Println(" >> ", app.SrmHomeDir)
		days = "18"
	} else {
		if _, err := strconv.ParseInt(days, 10, 64); err != nil {
			fmt.Fprintf(os.Stderr, "Parameter must be a number.")
			os.Exit(-1)
		}
	}

	day, _ := strconv.Atoi(days)
	srmCleanup(app, day)
}

func NewCleanupCommand(app app.Srm) *Command {
	cmd := &Command{
		flags: flag.NewFlagSet("cls", flag.ExitOnError),
		Execute: func(cmd *Command, args []string) {
			clsFunc(cmd, args, app)
		},
	}

	cmd.flags.StringVar(&days, "day", "", "")
	cmd.flags.StringVar(&days, "d", "", "")

	cmd.flags.Usage = func() {
		fmt.Fprintf(os.Stderr, clsUsage)
	}
	return cmd
}

func isOlderThan(t time.Time, day int) bool {
	return time.Now().Sub(t) > time.Duration(day)*24*time.Hour
}

func srmCleanup(app app.Srm, day int) error {
	tmpfiles, err := ioutil.ReadDir(app.SrmHomeDir)
	if err != nil {
		return err
	}

	for _, file := range tmpfiles {
		if file.Mode().IsRegular() {
			if isOlderThan(file.ModTime(), day) {
				// srmfile.SrmRemove(file.Name())
				fmt.Println("Purging... > ", file.Name())
			}
		}
	}
	return nil
}
