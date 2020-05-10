package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	cuetocu2 "github.com/ncirocco/cue-to-cu2"
	binmerge "github.com/ncirocco/psx-bin-merge"
	psxserialnumber "github.com/ncirocco/psx-serial-number"
)

const cu2Extension string = "cu2"
const cueExtension string = "cue"
const binExtension string = "bin"
const imagesEndpoint string = "https://ncirocco.github.io/PSIO-Library/images/covers_by_id/%s.bmp"
const fileNameMaxLength int = 60

func main() {
	defaultOriginDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	defaultDestinationDir := filepath.Join(defaultOriginDir, "output")

	autoCmd := flag.NewFlagSet("auto", flag.ExitOnError)
	autoDirPtr := autoCmd.String("dir", defaultOriginDir, "Directory containing the files to be processed")
	autoDestinationDirPtr := autoCmd.String("destinationDir", defaultDestinationDir, "Directory to store the processed files")

	mergeCmd := flag.NewFlagSet("merge", flag.ExitOnError)
	mergeDirPtr := mergeCmd.String("dir", defaultOriginDir, "Directory containing the files to be processed")
	mergeDestinationDirPtr := mergeCmd.String("destinationDir", defaultDestinationDir, "Directory to store the processed files")

	cu2Cmd := flag.NewFlagSet("cu2", flag.ExitOnError)
	cu2DirPtr := cu2Cmd.String("dir", defaultOriginDir, "Directory containing the cue files to be converted")
	cu2RemovePtr := cu2Cmd.Bool("removeCue", false, "If set to true, the original cue file will be removed")

	imagesCmd := flag.NewFlagSet("images", flag.ExitOnError)
	imagesDirPtr := imagesCmd.String("dir", defaultOriginDir, "Directory containing the bin files to get the images")

	flag.Usage = usage

	flag.Parse()

	if len(os.Args) < 2 {
		err = all(defaultOriginDir, defaultDestinationDir)
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}

	switch os.Args[1] {
	case "auto":
		autoCmd.Parse(os.Args[2:])
		err = all(*autoDirPtr, *autoDestinationDirPtr)
		if err != nil {
			log.Fatal(err)
		}
	case "merge":
		mergeCmd.Parse(os.Args[2:])
		err = mergeBins(*mergeDirPtr, *mergeDestinationDirPtr)
		if err != nil {
			log.Fatal(err)
		}
	case "cu2":
		cu2Cmd.Parse(os.Args[2:])
		err = generateCu2(*cu2DirPtr, *cu2RemovePtr)
		if err != nil {
			log.Fatal(err)
		}
	case "images":
		imagesCmd.Parse(os.Args[2:])
		err = getImages(*imagesDirPtr)
		if err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println(fmt.Sprintf("%s is not a valid command", os.Args[1]))
		os.Exit(1)
	}
}

func all(originDir string, destinationDir string) error {
	if _, err := os.Stat(originDir); os.IsNotExist(err) {
		return err
	}

	err := mergeBins(originDir, destinationDir)
	if err != nil {
		return err
	}

	err = generateCu2(destinationDir, true)
	if err != nil {
		return err
	}

	err = getImages(destinationDir)
	if err != nil {
		return err
	}

	err = trimFileNames(destinationDir)
	if err != nil {
		return err
	}

	return nil
}

func mergeBins(originDir string, destinationDir string) error {
	cueFiles, err := getFilesByExtension(originDir, cueExtension)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Attempting to merge bins for %d found CUE sheets\n\n", len(cueFiles))

	for _, cue := range cueFiles {
		fmt.Printf("Merging bin for %s\n", filepath.Base(cue))
		err := binmerge.Merge(cue, destinationDir)
		if err != nil {
			fmt.Printf("Error processing %s. Message: %s\n", cue, err)
		}
	}

	fmt.Print("\n\n")

	return nil
}

func generateCu2(dir string, removeCue bool) error {
	cueFiles, err := getFilesByExtension(dir, cueExtension)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Attempting to generate CU2 sheets for %d found CUE sheets\n\n", len(cueFiles))

	for _, cue := range cueFiles {
		fmt.Printf("Generating cu2 for %s\n", filepath.Base(cue))
		err := cuetocu2.Generate(cue, filepath.Dir(cue))
		if err != nil {
			fmt.Printf("Error processing %s. Message: %s\n", cue, err)
		}

		if removeCue {
			os.Remove(cue)
		}
	}

	fmt.Print("\n\n")

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

	return files, err
}

func usage() {
	fmt.Printf("usage: %s <command> [<args>]\n", os.Args[0])
	fmt.Println("\nCommands:")
	fmt.Println("  auto - Merges bin files and generates cu2 sheets for all the files in the given directory")
	fmt.Println("  merge - Merges all the bin files in a given directory")
	fmt.Println("  cu2 - Generates the cu2 sheet for each cue sheet in the given directory")
	fmt.Println("  image - Downloads covers for the bin files in the given directory")
	fmt.Print("\n")
	fmt.Printf("To see more details for a particular command run %s <command> -h\n", os.Args[0])
	fmt.Print("\n")
}

func getImages(dir string) error {
	bins, err := getFilesByExtension(dir, binExtension)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Attempting to download covers for %d found discs\n\n", len(bins))

	for _, bin := range bins {
		fmt.Printf("Downloading image for %s\n", filepath.Base(bin))
		serial, err := psxserialnumber.GetSerial(bin)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = downloadFile(serial, filepath.Dir(bin))
		if err != nil {
			fmt.Printf("%s for %s - serial %s\n", err, bin, serial)
			continue
		}
	}

	fmt.Print("\n\n")

	return nil
}

func downloadFile(code string, dir string) error {
	resp, err := http.Get(fmt.Sprintf(imagesEndpoint, strings.ToLower(code)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Image not found")
	}

	out, err := os.Create(filepath.Join(dir, path.Base(resp.Request.URL.String())))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// PSIO does not support files with names longer than 60 chars.
func trimFileNames(dir string) error {
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			ext := filepath.Ext(path)
			if (ext != "."+binExtension && ext != "."+cueExtension && ext != "."+cu2Extension) || len(filepath.Base(path)) <= fileNameMaxLength {
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
