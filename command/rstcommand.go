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
	"waldirborbajr/srm/internal/srmfile"
)

var rstUsage = `Restore a specific file/directory.

Usage: srm rst [file_name]

Options:
`

func RestoreCommand() *Command {
	cmd := &Command{
		flags: flag.NewFlagSet("rst", flag.ExitOnError),
		Execute: func(cmd *Command, args []string) {
			if len(args) == 0 {
				errAndExit("file name or directory is required")
			}
			file_name := args[0]

			homeDir, err := os.UserHomeDir()
			if err != nil {
				errAndExit("Failed to read home directory")
			}

			srmPathHome := filepath.Join(homeDir, ".srm")

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
							srmUncompress(filepath.Join(srmPathHome, rawContent))
							srmRestore(filepath.Join(srmPathHome, rawContent), filepath.Join(srmPath, srmFileName))
							srmfile.SrmCleanup(filepath.Join(srmPathHome, rawContent))
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
		srmFileName := srmFile[posStart+1 : posEnd]
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

func isExistsPath(targetPath string) bool {
	if _, err := os.Stat(targetPath); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	return false
}

func srmRestore(srmSrc string, srmTgt string) {
	// 1st uncompress file to retore

	// 1rawContentst Copy file to safety folder
	src, err := os.Open(srmSrc)
	if err != nil {
		fmt.Fprint(os.Stderr, "srm: unable to save file. [src]")
		os.Exit(-1)
	}
	dst, err := os.Create(srmTgt)
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
	// if err := os.Remove(file_name); err != nil {
	// 	fmt.Fprint(os.Stderr, "srm: unable to save file. [rmv]")
	// 	fmt.Println(err.Error())
	// 	os.Exit(-1)
	// }
	// fmt.Printf("srm: '%s' was safety deleted\n", file_name)
}

func srmUncompress(srmCompressedFile string) {
	filename := srmCompressedFile

	zlibfile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer zlibfile.Close()

	newfilename := strings.TrimSuffix(filename, ".zlib")
	decompressedFile, err := os.Create(newfilename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer decompressedFile.Close()

	reader, err := zlib.NewReader(zlibfile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer reader.Close()

	if _, err := io.Copy(decompressedFile, reader); err != nil {
		fmt.Println(err)
		return
	}
}
