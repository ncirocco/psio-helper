package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ncirocco/psio-helper/bin"
	"github.com/ncirocco/psio-helper/cu2"
	"github.com/ncirocco/psio-helper/files"
	"github.com/ncirocco/psio-helper/images"
	"github.com/ncirocco/psio-helper/multidisc"
)

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

	multidiscCmd := flag.NewFlagSet("multidisc", flag.ExitOnError)
	multidiscDirPtr := multidiscCmd.String("dir", defaultOriginDir, "Directory containing the multidisc bin files")

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
		err = bin.Merge(*mergeDirPtr, *mergeDestinationDirPtr)
		if err != nil {
			log.Fatal(err)
		}
	case "cu2":
		cu2Cmd.Parse(os.Args[2:])
		err = cu2.Generate(*cu2DirPtr, *cu2RemovePtr)
		if err != nil {
			log.Fatal(err)
		}
	case "images":
		imagesCmd.Parse(os.Args[2:])
		err = images.GetImages(*imagesDirPtr)
		if err != nil {
			log.Fatal(err)
		}
	case "multidisc":
		multidiscCmd.Parse(os.Args[2:])
		err = multidisc.Multidisc(*multidiscDirPtr)
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

	err := bin.Merge(originDir, destinationDir)
	if err != nil {
		return err
	}

	err = cu2.Generate(destinationDir, true)
	if err != nil {
		return err
	}

	err = multidisc.Multidisc(destinationDir)
	if err != nil {
		return err
	}

	err = images.GetImages(destinationDir)
	if err != nil {
		return err
	}

	err = files.TrimNames(destinationDir)
	if err != nil {
		return err
	}

	return nil
}

func usage() {
	fmt.Printf("usage: %s <command> [<args>]\n", os.Args[0])
	fmt.Println("\nCommands:")
	fmt.Println("  auto - Merges bin files and generates cu2 sheets for all the files in the given directory")
	fmt.Println("  merge - Merges all the bin files in a given directory")
	fmt.Println("  cu2 - Generates the cu2 sheet for each cue sheet in the given directory")
	fmt.Println("  image - Downloads covers for the bin files in the given directory")
	fmt.Println("  multidisc - Groups the discs that belong to the same game and generates the necessary multidisc.lst files")
	fmt.Print("\n")
	fmt.Printf("To see more details for a particular command run %s <command> -h\n", os.Args[0])
	fmt.Print("\n")
}
