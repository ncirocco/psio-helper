package bin

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/ncirocco/psio-helper/files"
	binmerge "github.com/ncirocco/psx-bin-merge"
)

//Merge merges the bin files in the origin directory and outputs them in the destination directory
func Merge(originDir string, destinationDir string) error {
	cueFiles, err := files.GetFilesByExtension(originDir, files.CueExtension)
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
