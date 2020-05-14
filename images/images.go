package images

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ncirocco/psio-helper/files"
	psxserialnumber "github.com/ncirocco/psx-serial-number"
)

const imagesEndpoint string = "https://ncirocco.github.io/PSIO-Library/images/covers_by_id/%s.bmp"

// GetImages downloads the images for the given directory
func GetImages(dir string) error {
	bins, err := files.GetFilesByExtension(dir, files.BinExtension)
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
