package cu2

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	cuetocu2 "github.com/ncirocco/cue-to-cu2"
	"github.com/ncirocco/psio-helper/files"
)

const concurrentOpenedFiles = 25

var sem = make(chan struct{}, concurrentOpenedFiles)

//Generate creates the cu2 files for all the found cue in the given directory
func Generate(dir string, removeCue bool) error {
	cueFiles, err := files.GetFilesByExtension(dir, files.CueExtension)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Attempting to generate CU2 sheets for %d found CUE sheets\n\n", len(cueFiles))

	var wg sync.WaitGroup
	for _, cue := range cueFiles {
		wg.Add(1)
		go generate(&wg, cue, removeCue)
	}

	wg.Wait()

	fmt.Print("\n\n")

	return nil
}

func generate(wg *sync.WaitGroup, cue string, removeCue bool) {
	sem <- struct{}{}
	defer func() { <-sem }()
	defer wg.Done()

	fmt.Printf("Generating cu2 for %s\n", filepath.Base(cue))
	err := cuetocu2.Generate(cue, filepath.Dir(cue))
	if err != nil {
		fmt.Printf("Error processing %s. Message: %s\n", cue, err)
		return
	}

	if removeCue {
		os.Remove(cue)
	}
}
