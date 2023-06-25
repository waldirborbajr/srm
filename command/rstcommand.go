package command

import (
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"waldirborbajr/srm/internal/app"
	"waldirborbajr/srm/internal/srmfile"
)

var rstUsage = `Restore a specific file/directory.

Usage: srm rst [file_name]

Options:
`

func NewRestoreCommand(app app.Srm) *Command {
	cmd := &Command{
		flags: flag.NewFlagSet("rst", flag.ExitOnError),
		Execute: func(cmd *Command, args []string) {
			if len(args) == 0 {
				errAndExit("file name or directory is required")
			}
			file_name := args[0]

			srmPathHome := app.SrmHomeDir

			// List files from srm restore point
			fileToRestore, err := os.ReadDir(srmPathHome)
			if err != nil {
				errAndExit("Failed to read restore point directory > " + err.Error())
			}

			for _, content := range fileToRestore {
				rawContent := content.Name()
				srmFileName := splitFileName(rawContent)
				srmPath := splitPath(rawContent)

				if srmFileName != "" && srmPath != "" {
					if file_name == srmFileName {
						if isExistsPath(srmPath) {

							// 1st uncompress file
							if err := srmUncompress(filepath.Join(srmPathHome, rawContent)); err != nil {
								fmt.Fprintf(os.Stderr, "Error: %v", err)
								os.Exit(-1)
							}

							// 2nd move file from safety folder to destination
							if err := srmRestore(filepath.Join(srmPathHome, rawContent), filepath.Join(srmPath, srmFileName)); err != nil {
								fmt.Fprintf(os.Stderr, "Error: %v", err)
								os.Exit(-1)
							}

							if err := srmfile.SrmCleanup(filepath.Join(srmPathHome, rawContent)); err != nil {
								fmt.Fprintf(os.Stderr, "Error: %v", err)
								os.Exit(-1)
							}
						}
					}
				}
			}
		},
	}

	cmd.flags.Usage = func() {
		fmt.Fprintln(os.Stderr, rstUsage)
	}

	return cmd
}

func splitFileName(rawPath string) string {
	srmFile := rawPath
	posStart := strings.Index(srmFile, "}")
	posEnd := strings.Index(srmFile, ".zlib")

	if posStart != -1 && posEnd != -1 {
		return srmFile[posStart+1 : posEnd]
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

func isExistsPath(targetPath string) bool {
	if _, err := os.Stat(targetPath); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	return false
}

func srmRestore(srmSrc string, srmTgt string) error {
	position := strings.Index(srmSrc, ".zlib")
	if position != -1 {
		srmSrc = srmSrc[:position]
	}

	if err := os.Rename(srmSrc, srmTgt); err != nil {
		return err
	}

	return nil
}

func srmUncompress(srmCompressedFile string) error {
	filename := srmCompressedFile

	zlibfile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer zlibfile.Close()

	newfilename := strings.TrimSuffix(filename, ".zlib")
	decompressedFile, err := os.Create(newfilename)
	if err != nil {
		return err
	}
	defer decompressedFile.Close()

	reader, err := zlib.NewReader(zlibfile)
	if err != nil {
		return err
	}
	defer reader.Close()

	if _, err := io.Copy(decompressedFile, reader); err != nil {
		return err
	}

	return nil
}
