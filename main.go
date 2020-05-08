package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	cuetocu2 "github.com/ncirocco/cue-to-cu2"
	binmerge "github.com/ncirocco/psx-bin-merge"
	psxserialnumber "github.com/ncirocco/psx-serial-number"
)

func main() {
	originPath := "."
	outputPath := "output"

	if len(os.Args) > 1 {
		originPath = os.Args[1]
	}

	if len(os.Args) > 2 {
		outputPath = os.Args[2]
	}

	if _, err := os.Stat(originPath); os.IsNotExist(err) {
		log.Fatal(err)
	}

	err := mergeBins(originPath, outputPath)
	if err != nil {
		log.Fatal(err)
	}

	err = generateCu2(outputPath)
	if err != nil {
		log.Fatal(err)
	}
}

func mergeBins(originPath string, outputPath string) error {
	cueFiles, err := getFilesByExtension(originPath, "cue")
	if err != nil {
		log.Fatal(err)
	}

	for _, cue := range cueFiles {
		err := binmerge.Merge(cue, outputPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateCu2(outputPath string) error {
	mergedCueFiles, err := getFilesByExtension(outputPath, "cue")
	if err != nil {
		log.Fatal(err)
	}

	for _, cue := range mergedCueFiles {
		err := cuetocu2.Generate(cue, filepath.Dir(cue))
		os.Remove(cue)
		if err != nil {
			return err
		}
	}

	return nil
}

func getFilesByExtension(path string, extension string) ([]string, error) {
	var files []string
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filepath.Ext(path) == "."+extension {
				files = append(files, path)
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	return files, err
}

func getSerial(outputPath string) {
	bins, err := getFilesByExtension(outputPath, "bin")
	if err != nil {
		log.Fatal(err)
	}

	for _, bin := range bins {
		serial, err := psxserialnumber.GetSerial(bin)
		if err != nil {
			//log.Fatal(err)
			continue
		}

		fmt.Println(serial)
	}
}
