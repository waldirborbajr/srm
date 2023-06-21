package command

import (
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"waldirborbajr/srm/internal/app"
	"waldirborbajr/srm/internal/srmfile"
)

var srmUsage = `Removes a specific file/directory.

Usage: srm srm file.bak

Options:
`

func NewSafeRemoveCommand(app app.Srm) *Command {
	cmd := &Command{
		flags: flag.NewFlagSet("srm", flag.ExitOnError),
		Execute: func(cmd *Command, args []string) {
			if len(args) == 0 {
				errAndExit("file name or directory is required")
			}

			file_name := flag.Args()
			// file_name := args[0]

			for _, file := range file_name {
				srmMatches, err := filepath.Glob(file)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				}

				for _, match := range srmMatches {
					srmPathAbs, _ := filepath.Abs(match)
					idx := strings.Index(srmPathAbs, match)
					srmSourcePath := srmPathAbs[:idx]
					source := "{" + strings.Replace(srmSourcePath, "/", "-", -1) + "}"
					filePath := app.SrmHomeDir + "/" + source + match

					// 1st Copy file to safety folder
					srmDoCopy(match, filePath)

					// 2nd Compress file on target
					srmCompress(filePath)

					// 3rd Remove source file
					if err := srmfile.SrmRemove(match); err != nil {
						fmt.Fprint(os.Stderr, "srm: unable to save file. [rmv]")
						fmt.Println(err.Error())
						os.Exit(-1)
					}

					// 4th Remove target uncompressed file
					if err := srmfile.SrmRemove(filePath); err != nil {
						fmt.Fprint(os.Stderr, "srm: unable to save file. [rmv]")
						fmt.Println(err.Error())
						os.Exit(-1)
					}

					fmt.Fprint(os.Stdout, "srm: ", match, " was safety deleted.\n")
				}
			}

			os.Exit(0)
			// srmPathAbs, _ := filepath.Abs(file_name)
			// idx := strings.Index(srmPathAbs, file_name)
			// srmSourcePath := srmPathAbs[:idx]

			// source := "{" + strings.Replace(srmSourcePath, "/", "-", -1) + "}"

			// filePath := app.SrmHomeDir + "/" + source + file_name

			// // 1st Copy file to safety folder
			// srmDoCopy(file_name, filePath)

			// // 2nd Compress file on target
			// srmCompress(filePath)

			// // 3rd Remove source file
			// if err := srmfile.SrmRemove(file_name); err != nil {
			// 	fmt.Fprint(os.Stderr, "srm: unable to save file. [rmv]")
			// 	fmt.Println(err.Error())
			// 	os.Exit(-1)
			// }

			// // 4th Remove target uncompressed file
			// if err := srmfile.SrmRemove(filePath); err != nil {
			// 	fmt.Fprint(os.Stderr, "srm: unable to save file. [rmv]")
			// 	fmt.Println(err.Error())
			// 	os.Exit(-1)
			// }

			// fmt.Fprint(os.Stdout, "srm: ", file_name, " was safety deleted.")
			// os.Exit(0)
		},
	}

	cmd.flags.Usage = func() {
		fmt.Fprintln(os.Stderr, srmUsage)
	}

	return cmd
}

func srmDoCopy(srcFileName string, tgtPath string) {
	src, err := os.Open(srcFileName)
	if err != nil {
		fmt.Fprint(os.Stderr, "srm: unable to save file. [src]")
		os.Exit(-1)
	}

	dst, err := os.Create(tgtPath)
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
}

func srmCompress(srmFileToCompress string) {
	// Create a new file "example.zlib"
	file, err := os.Create(srmFileToCompress + ".zlib")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a new zlib writer with the best compression level
	writer, err := zlib.NewWriterLevel(file, zlib.BestCompression)
	if err != nil {
		panic(err)
	}
	defer writer.Close()

	// Open the input file "example.txt"
	inputFile, err := os.Open(srmFileToCompress)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	// Copy the contents of the input file to the writer
	io.Copy(writer, inputFile)
}
