package command

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var rstUsage = `Restone a specific file/directory.

Usage: srm rst file.bak

Options:
`

func RestoreCommand() *Command {
	cmd := &Command{
		flags: flag.NewFlagSet("restore", flag.ExitOnError),
		Execute: func(cmd *Command, args []string) {
			if len(args) == 0 {
				errAndExit("file name or directory is required")
			}
			file_name := args[0]

			var wd string
			var full_path string
			var err error

			if wd, err = os.Getwd(); err != nil {
				errAndExit("Unable to determinate content fullpath")
			}

			if full_path, err = filepath.Abs(filepath.Join(wd)); err != nil {
				errAndExit("Unable to determinate content fullpath")
			}

			source := "{" + strings.Replace(full_path, "/", "-", -1) + "}"

			homeDir, err := os.UserHomeDir()
			if err != nil {
				errAndExit("Failed to read home directory")
			}

			filePath := filepath.Join(homeDir, ".srm", source) + file_name
			fmt.Println(filePath)

			// 1st Copy file to safety folder
			src, err := os.Open(file_name)
			if err != nil {
				fmt.Fprint(os.Stderr, "srm: unable to save file. [src]")
				os.Exit(-1)
			}
			dst, err := os.Create(filePath)
			if err != nil {
				fmt.Fprint(os.Stderr, "srm: unable to save file. [dst]")
				os.Exit(-1)
			}

			defer dst.Close()
			_, err = io.Copy(dst, src)
			src.Close()
			if err != nil {
				fmt.Fprint(os.Stderr, "srm: unable to save file. [cpy]")
				os.Exit(-1)
			}

			// 2nd Remove file
			if err := os.Remove(file_name); err != nil {
				fmt.Fprint(os.Stderr, "srm: unable to save file. [rmv]")
				fmt.Println(err.Error())
				os.Exit(-1)
			}
			fmt.Printf("srm: '%s' was safety deleted\n", file_name)
		},
	}

	cmd.flags.Usage = func() {
		fmt.Fprintln(os.Stderr, rstUsage)
	}

	return cmd
}
