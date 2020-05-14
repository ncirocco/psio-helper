package cu2

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	cuetocu2 "github.com/ncirocco/cue-to-cu2"
	"github.com/ncirocco/psio-helper/files"
)

//Generate creates the cu2 files for all the found cue in the given directory
func Generate(dir string, removeCue bool) error {
	cueFiles, err := files.GetFilesByExtension(dir, files.CueExtension)
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
