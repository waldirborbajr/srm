package command

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"waldirborbajr/srm/internal/app"
	"waldirborbajr/srm/internal/srmfile"
)

var days string

var clsUsage = `Purge files from safety storage after "18" days if no parameter it is informed.

Usage: srm cls [OPTIONS]

Options:
	-d, --day days to purge srm safer folder
`

func clsFunc(cmd *Command, args []string, app app.Srm) {
	if len(days) == 0 {
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
	return int(time.Since(t).Hours()) > 24*day
}

func srmCleanup(app app.Srm, day int) error {
	tmpfiles, err := os.ReadDir(app.SrmHomeDir)
	if err != nil {
		return err
	}

	for _, file := range tmpfiles {
		if !file.IsDir() {
			info, err := file.Info()
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: getting file.Ifo() %s", filepath.Join(app.SrmHomeDir, file.Name()))
				return err
			}

			if isOlderThan(info.ModTime(), day) {
				fmt.Println("Removing : ", filepath.Join(app.SrmHomeDir, file.Name()))
				if err := srmfile.SrmRemove(filepath.Join(app.SrmHomeDir, file.Name())); err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: removing %s", filepath.Join(app.SrmHomeDir, file.Name()))
				}
			}
		}
	}
	return nil
}
