package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//Cu2Extension cu2
const Cu2Extension string = "cu2"

//CueExtension cue
const CueExtension string = "cue"

//BinExtension bin
const BinExtension string = "bin"

const fileNameMaxLength int = 60

//GetFilesByExtension returns a list of the existing files with the given extension in the given directory
func GetFilesByExtension(path string, extension string) ([]string, error) {
	var files []string
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if strings.ToLower(filepath.Ext(path)) == "."+extension {
				files = append(files, path)
			}

			return nil
		})

	return files, err
}

// HasFilesWithExtension returns true if the directory has files with the given extension
func HasFilesWithExtension(path string, extension string) (bool, error) {
	var files []string
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if strings.ToLower(filepath.Ext(path)) == "."+extension {
				files = append(files, path)
			}

			return nil
		})

	return len(files) > 0, err
}

// TrimNames shortens the file names to the PSIO limit of 60 chars.
func TrimNames(dir string) error {
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			ext := filepath.Ext(path)
			if (ext != "."+BinExtension && ext != "."+CueExtension && ext != "."+Cu2Extension) || len(filepath.Base(path)) <= fileNameMaxLength {
				return nil
			}

			fmt.Printf("File name for %s is longer than 60, trimming it.\n", filepath.Base(path))

			newName := filepath.Join(filepath.Dir(path), filepath.Base(path)[:fileNameMaxLength-len(ext)]+ext)

			if _, err := os.Stat(newName); err == nil {
				fmt.Printf("can't rename %s. please rename it manually", filepath.Base(path))

				return nil
			}

			os.Rename(path, newName)

			return nil
		})

	return err
}
