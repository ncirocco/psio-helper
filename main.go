package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	binmerge "github.com/ncirocco/psx-bin-merge"
	psxserialnumber "github.com/ncirocco/psx-serial-number"
)

func main() {
	cueFiles, err := getCueFiles()
	if err != nil {
		log.Fatal(err)
	}

	for _, cue := range cueFiles {
		err = binmerge.Merge(cue, "carpeta")
		if err != nil {
			log.Fatal(err)
		}
	}

	if true {
		return
	}

	bins, err := checkDirectory()
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

func getCueFiles() ([]string, error) {
	var bins []string
	err := filepath.Walk("./ISOs",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filepath.Ext(path) == ".cue" {
				bins = append(bins, path)
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	return bins, err
}

func checkDirectory() ([]string, error) {
	var bins []string
	err := filepath.Walk("./ISOs",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filepath.Ext(path) == ".bin" {
				bins = append(bins, path)
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	return bins, err
}
