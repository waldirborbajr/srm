package command

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var rstUsage = `Restore a specific file/directory.

Usage: srm rst [file_name]

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

			homeDir, err := os.UserHomeDir()
			if err != nil {
				errAndExit("Failed to read home directory")
			}

			srmPath := filepath.Join(homeDir, ".srm")

			// List files from srm restore point
			fileToRestore, err := os.ReadDir(srmPath)
			if err != nil {
				errAndExit("Failed to read restore poit directory > " + err.Error())
			}

			for _, content := range fileToRestore {
				rawContent := content.Name()
				srmFileName := splitFileName(rawContent)
				srmPath := splitPath(rawContent)
				if srmFileName != "" && srmPath != "" {
					if file_name == srmFileName {
						fmt.Println(" restoring > " + srmFileName + " to > " + srmPath)
					}
				}
			}

			// restorePath := filepath.Join(homeDir, ".srm", source) + file_name
			// fmt.Println(restorePath)

			// 1st Copy file to safety folder
			// src, err := os.Open(file_name)
			// if err != nil {
			// 	fmt.Fprint(os.Stderr, "srm: unable to save file. [src]")
			// 	os.Exit(-1)
			// }
			// dst, err := os.Create(filePath)
			// if err != nil {
			// 	fmt.Fprint(os.Stderr, "srm: unable to save file. [dst]")
			// 	os.Exit(-1)
			// }

			// defer dst.Close()
			// _, err = io.Copy(dst, src)
			// src.Close()
			// if err != nil {
			// 	fmt.Fprint(os.Stderr, "srm: unable to save file. [cpy]")
			// 	os.Exit(-1)
			// }

			// 2nd Remove file
			// if err := os.Remove(file_name); err != nil {
			// 	fmt.Fprint(os.Stderr, "srm: unable to save file. [rmv]")
			// 	fmt.Println(err.Error())
			// 	os.Exit(-1)
			// }
			// fmt.Printf("srm: '%s' was safety deleted\n", file_name)
		},
	}

	cmd.flags.Usage = func() {
		fmt.Fprintln(os.Stderr, rstUsage)
	}

	return cmd
}

func splitFileName(rawPath string) string {
	srmFile := rawPath
	pos := strings.Index(srmFile, "}")
	if pos != -1 && pos < len(srmFile)-1 {
		srmFileName := srmFile[pos+1:]
		return srmFileName
	}
	return ""
}

func splitPath(rawPath string) string {
	pattern := regexp.MustCompile(`\{([^}]+)\}`)
	matches := pattern.FindStringSubmatch(rawPath)

	if len(matches) > 1 {
		content := strings.Replace(matches[1], "-", "/", -1)
		return content
	}
	return ""
}
