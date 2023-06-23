package command

import (
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"waldirborbajr/srm/internal/app"
	"waldirborbajr/srm/internal/srmfile"
)

var (
	isSafe   bool
	isForce  bool
	srmUsage = `Usage: srm srm [options....]

Usage: 
	srm srm --save file.bak
	srm srm --force file.bak
	srm srm --save file*

Options:
	-s, --safe Save removed file/directory for restore
	-f, --force Remove file/directory without restore option
`
)

func NewSafeRemoveCommand(app app.Srm) *Command {
	cmd := &Command{
		flags: flag.NewFlagSet("srm", flag.ExitOnError),
		Execute: func(cmd *Command, args []string) {
			srmFunc(cmd, args, app)
		},
	}

	cmd.flags.BoolVar(&isForce, "force", false, "force remove without save")
	cmd.flags.BoolVar(&isForce, "f", false, "remove saving for restore")
	cmd.flags.BoolVar(&isSafe, "safe", false, "remove saving for restore")
	cmd.flags.BoolVar(&isSafe, "s", false, "remove saving for restore")

	cmd.flags.Usage = func() {
		fmt.Fprintf(os.Stderr, srmUsage)
	}

	return cmd
}

func srmFunc(cmd *Command, args []string, app app.Srm) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, srmUsage)
		os.Exit(-1)
	}

	if !isForce && !isSafe {
		fmt.Fprintf(os.Stderr, srmUsage)
		os.Exit(-1)
	}

	srmContentToRemove := cmd.flags.Args()

	for _, file := range srmContentToRemove {
		srmMatches, err := filepath.Glob(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}

		for _, srmFile := range srmMatches {
			// Stat the content to identify if it is file or directory
			info, err := os.Stat(srmFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(-1)
			}
			if !info.IsDir() {
				srmRemoveFile(app.SrmHomeDir, srmFile, isSafe)
			} else {
				srmRemoveDirectory(app.SrmHomeDir, srmFile, isSafe)
			}
		}
	}

	os.Exit(0)
}

func srmRemoveFile(srmHomeDir string, srmParamFileName string, hasSafe bool) {
	srmAbsPath, _ := filepath.Abs(srmParamFileName)
	idx := strings.Index(srmAbsPath, srmParamFileName)
	srmSourcePath := srmAbsPath[:idx]
	source := "{" + strings.Replace(srmSourcePath, "/", "-", -1) + "}"
	srmDestinationPath := srmHomeDir + "/" + source + srmParamFileName

	if hasSafe {
		// 1st Copy file to safety folder
		srmDoCopyFile(srmParamFileName, srmDestinationPath)

		// 2nd Compress file on target
		srmCompress(srmDestinationPath)
	}

	// 3rd Remove source file
	if err := srmfile.SrmRemove(srmParamFileName); err != nil {
		fmt.Fprint(os.Stderr, "srm: unable to save file. [rmv]")
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	if hasSafe {
		// 4th Remove target uncompressed file
		if err := srmfile.SrmRemove(srmDestinationPath); err != nil {
			fmt.Fprint(os.Stderr, "srm: unable to save file. [rmv]")
			fmt.Println(err.Error())
			os.Exit(-1)
		}
	}

	if hasSafe {
		fmt.Fprint(os.Stdout, "srm: ", srmParamFileName, " was safety deleted.\n")
	} else {
		fmt.Fprint(os.Stdout, "srm: ", srmParamFileName, " was deleted.\n")
	}
}

func srmRemoveDirectory(srmHomeDir string, srmParamFileName string, hasSafe bool) {
	srmSourcePath, _ := filepath.Abs(srmParamFileName)
	source := "{" + strings.Replace(srmSourcePath, "/", "-", -1) + "}"
	srmDestinationPath := srmHomeDir + "/" + source

	if hasSafe {
		// 1st Copy file to safety folder
		if err := srmDoCopyDirectory(srmSourcePath, srmDestinationPath); err != nil {
			fmt.Fprint(os.Stderr, "Error: saving folder")
			os.Exit(-1)
		}

		// 2nd Compress file on target
		srmCompress(srmDestinationPath)
	}

	// 3rd Remove source file
	if err := srmfile.SrmRemoveDirectory(srmSourcePath); err != nil {
		fmt.Fprint(os.Stderr, "srm: unable to save file. [rmv]")
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	if hasSafe {
		// 4th Remove target uncompressed file
		if err := srmfile.SrmRemoveDirectory(srmSourcePath); err != nil {
			fmt.Fprint(os.Stderr, "srm: unable to save file. [rmv]")
			fmt.Println(err.Error())
			os.Exit(-1)
		}
	}

	if hasSafe {
		fmt.Fprint(os.Stdout, "srm: ", srmParamFileName, " was safety deleted.\n")
	} else {
		fmt.Fprint(os.Stdout, "srm: ", srmParamFileName, " was deleted.\n")
	}
}

func srmDoCopyFile(srcFileName string, tgtPath string) {
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

func srmDoCopyDirectory(srmSourcePath string, srmDestinationPath string) error {
	return filepath.Walk(srmSourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}

		newPath := strings.Replace(path, srmSourcePath, srmDestinationPath, 1)
		if info.IsDir() {
			return os.MkdirAll(newPath, info.Mode())
		} else {
			return os.Rename(path, newPath)
		}
	})
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
