package multidisc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ncirocco/psio-helper/files"
	psxserialnumber "github.com/ncirocco/psx-serial-number"
)

type multiDisc struct {
	name string
	disc string
}

var multidiscMap = map[string]multiDisc{
	"SLUS_012.01": {"Alone in the Dark - The New Nightmare", "Disc 1"},
	"SLUS_013.77": {"Alone in the Dark - The New Nightmare", "Disc 2"},
	"SLUS_012.53": {"Arc the Lad Collection - Arc the Lad III", "Disc 1"},
	"SLUS_012.54": {"Arc the Lad Collection - Arc the Lad III", "Disc 2"},
	"SLUS_010.30": {"Armored Core - Master of Arena", "Disc 1"},
	"SLUS_010.81": {"Armored Core - Master of Arena", "Disc 2"},
	"SLUS_000.83": {"BrainDead 13", "Disc 1"},
	"SLUS_001.71": {"BrainDead 13", "Disc 2"},
	"SCUS_947.00": {"Chronicles of the Sword", "Disc 1"},
	"SCUS_947.01": {"Chronicles of the Sword", "Disc 2"},
	"SLUS_010.41": {"Chrono Cross", "Disc 1"},
	"SLUS_010.80": {"Chrono Cross", "Disc 2"},
	"SLUS_005.43": {"Colony Wars", "Disc 1"},
	"SLUS_005.54": {"Colony Wars", "Disc 2"},
	"SLUS_003.79": {"Command & Conquer", "GDI"},
	"SLUS_004.10": {"Command & Conquer", "NOD"},
	"SLUS_004.31": {"Command & Conquer - Red Alert", "Allies"},
	"SLUS_004.85": {"Command & Conquer - Red Alert", "Soviet"},
	"SLUS_006.65": {"Command & Conquer - Red Alert - Retaliation", "Allies"},
	"SLUS_006.67": {"Command & Conquer - Red Alert - Retaliation", "Soviet"},
	"SLUS_008.98": {"Countdown Vampires", "Disc 1"},
	"SLUS_011.99": {"Countdown Vampires", "Disc 2"},
	"SLUS_011.51": {"Covert Ops - Nuclear Dawn", "Disc 1"},
	"SLUS_011.57": {"Covert Ops - Nuclear Dawn", "Disc 2"},
	"SLUS_001.28": {"D", "Disc 1"},
	"SLUS_001.73": {"D", "Disc 2"},
	"SLUS_001.74": {"D", "Disc 3"},
	"SLUS_014.40": {"Dracula - The Last Sanctuary", "Disc 1"},
	"SLUS_014.43": {"Dracula - The Last Sanctuary", "Disc 2"},
	"SLUS_012.84": {"Dracula - The Resurrection", "Disc 1"},
	"SLUS_013.16": {"Dracula - The Resurrection", "Disc 2"},
	"SLUS_010.92": {"Dragon Valor", "Disc 1"},
	"SLUS_011.64": {"Dragon Valor", "Disc 2"},
	"SLUS_012.06": {"Dragon Warrior VII", "Disc 1"},
	"SLUS_013.46": {"Dragon Warrior VII", "Disc 2"},
	"SLUS_011.61": {"Driver 2", "Disc 1"},
	"SLUS_013.18": {"Driver 2", "Disc 2"},
	"SLUS_010.72": {"Evil Dead - Hail to the King", "Disc 1"},
	"SLUS_013.26": {"Evil Dead - Hail to the King", "Disc 2"},
	"SLUS_009.20": {"Fear Effect", "Disc 1"},
	"SLUS_010.56": {"Fear Effect", "Disc 2"},
	"SLUS_010.57": {"Fear Effect", "Disc 3"},
	"SLUS_010.58": {"Fear Effect", "Disc 4"},
	"SLUS_012.66": {"Fear Effect 2 - Retro Helix", "Disc 1"},
	"SLUS_012.75": {"Fear Effect 2 - Retro Helix", "Disc 2"},
	"SLUS_012.76": {"Fear Effect 2 - Retro Helix", "Disc 3"},
	"SLUS_012.77": {"Fear Effect 2 - Retro Helix", "Disc 4"},
	"SLUS_012.51": {"Final Fantasy IX", "Disc 1"},
	"SLUS_012.95": {"Final Fantasy IX", "Disc 2"},
	"SLUS_012.96": {"Final Fantasy IX", "Disc 3"},
	"SLUS_012.97": {"Final Fantasy IX", "Disc 4"},
	"SCUS_941.63": {"Final Fantasy VII", "Disc 1"},
	"SCUS_941.64": {"Final Fantasy VII", "Disc 2"},
	"SCUS_941.65": {"Final Fantasy VII", "Disc 3"},
	"SLUS_008.92": {"Final Fantasy VIII", "Disc 1"},
	"SLUS_009.08": {"Final Fantasy VIII", "Disc 2"},
	"SLUS_009.09": {"Final Fantasy VIII", "Disc 3"},
	"SLUS_009.10": {"Final Fantasy VIII", "Disc 4"},
	"SLUS_001.01": {"Fox Hunt", "Disc 1"},
	"SLUS_001.75": {"Fox Hunt", "Disc 2"},
	"SLUS_000.17": {"Fox Hunt", "Disc 3"},
	"SLUS_005.44": {"G-Police", "Disc 1"},
	"SLUS_005.56": {"G-Police", "Disc 2"},
	"SLUS_009.86": {"Galerians", "Disc 1"},
	"SLUS_010.98": {"Galerians", "Disc 2"},
	"SLUS_010.99": {"Galerians", "Disc 3"},
	"SLUS_003.19": {"Golden Nugget", "Disc 1"},
	"SLUS_005.55": {"Golden Nugget", "Disc 2"},
	"SCUS_944.57": {"Grandia", "Disc 1"},
	"SCUS_944.65": {"Grandia", "Disc 2"},
	"SLUS_006.96": {"Heart of Darkness", "Disc 1"},
	"SLUS_007.41": {"Heart of Darkness", "Disc 2"},
	"SLUS_001.20": {"The Hive", "Disc 1"},
	"SLUS_001.82": {"The Hive", "Disc 2"},
	"SLUS_012.94": {"In Cold Blood", "Disc 1"},
	"SLUS_013.14": {"In Cold Blood", "Disc 2"},
	"SLUS_008.94": {"Juggernaut", "Disc 1"},
	"SLUS_009.88": {"Juggernaut", "Disc 2"},
	"SLUS_009.89": {"Juggernaut", "Disc 3"},
	"SLUS_010.51": {"Koudelka", "Disc 1"},
	"SLUS_011.00": {"Koudelka", "Disc 2"},
	"SLUS_011.01": {"Koudelka", "Disc 3"},
	"SLUS_011.02": {"Koudelka", "Disc 4"},
	"SCUS_944.91": {"The Legend of Dragoon", "Disc 1"},
	"SCUS_945.84": {"The Legend of Dragoon", "Disc 2"},
	"SCUS_945.85": {"The Legend of Dragoon", "Disc 3"},
	"SCUS_945.86": {"The Legend of Dragoon", "Disc 4"},
	"SLUS_006.28": {"Lunar - Silver Star Story Complete", "Disc 1"},
	"SLUS_008.99": {"Lunar - Silver Star Story Complete", "Disc 2"},
	"SLUS_010.71": {"Lunar 2 - Eternal Blue Complete", "Disc 1"},
	"SLUS_012.39": {"Lunar 2 - Eternal Blue Complete", "Disc 2"},
	"SLUS_012.40": {"Lunar 2 - Eternal Blue Complete", "Disc 3"},
	"SLUS_005.94": {"Metal Gear Solid", "Disc 1"},
	"SLUS_007.76": {"Metal Gear Solid", "Disc 2"},
	"SCUS_944.04": {"Novastorm", "Disc 1"},
	"SCUS_944.07": {"Novastorm", "Disc 2"},
	"SLUS_007.10": {"Oddworld - Abe's Exoddus", "Disc 1"},
	"SLUS_007.31": {"Oddworld - Abe's Exoddus", "Disc 2"},
	"SLUS_006.62": {"Parasite Eve", "Disc 1"},
	"SLUS_006.68": {"Parasite Eve", "Disc 2"},
	"SLUS_010.42": {"Parasite Eve II", "Disc 1"},
	"SLUS_010.55": {"Parasite Eve II", "Disc 2"},
	"SLUS_001.65": {"Psychic Detective", "Disc 1"},
	"SLUS_001.66": {"Psychic Detective", "Disc 2"},
	"SLUS_001.67": {"Psychic Detective", "Disc 3"},
	"SLUS_004.21": {"Resident Evil 2", "Leon"},
	"SLUS_005.92": {"Resident Evil 2", "Claire"},
	"SLUS_007.48": {"Resident Evil 2 - Dual Shock", "Leon"},
	"SLUS_007.56": {"Resident Evil 2 - Dual Shock", "Claire"},
	"SLUS_006.81": {"Rival Schools - United by Fate", "Arcade"},
	"SLUS_007.71": {"Rival Schools - United by Fate", "Evolution"},
	"SLUS_005.35": {"Riven - The Sequel to Myst", "Disc 1"},
	"SLUS_005.63": {"Riven - The Sequel to Myst", "Disc 2"},
	"SLUS_005.64": {"Riven - The Sequel to Myst", "Disc 3"},
	"SLUS_005.65": {"Riven - The Sequel to Myst", "Disc 4"},
	"SLUS_005.80": {"Riven - The Sequel to Myst", "Disc 5"},
	"SLUS_004.68": {"Shadow Madness", "Disc 1"},
	"SLUS_007.18": {"Shadow Madness", "Disc 2"},
	"SLUS_000.28": {"Shockwave Assault", "Shockwave - Invasion Earth"},
	"SLUS_001.37": {"Shockwave Assault", "Shockwave - Operation Jumpgate"},
	"SCUS_944.21": {"Star Ocean - The Second Story", "Disc 1"},
	"SCUS_944.22": {"Star Ocean - The Second Story", "Disc 2"},
	"SLUS_003.81": {"Star Wars - Rebel Assault II - The Hidden Empire", "Disc 1"},
	"SLUS_003.86": {"Star Wars - Rebel Assault II - The Hidden Empire", "Disc 2"},
	"SLUS_004.23": {"Street Fighter Collection", "Disc 1"},
	"SLUS_005.84": {"Street Fighter Collection", "Disc 2"},
	"SCUS_944.51": {"Syphon Filter 2", "Disc 1"},
	"SCUS_944.92": {"Syphon Filter 2", "Disc 2"},
	"SLUS_013.55": {"Tales of Destiny II", "Disc 1"},
	"SLUS_013.67": {"Tales of Destiny II", "Disc 2"},
	"SLUS_013.68": {"Tales of Destiny II", "Disc 3"},
	"SLUS_008.45": {"Thousand Arms", "Disc 1"},
	"SLUS_008.58": {"Thousand Arms", "Disc 2"},
	"SLUS_011.56": {"Valkyrie Profile", "Disc 1"},
	"SLUS_011.79": {"Valkyrie Profile", "Disc 2"},
	"SCUS_944.84": {"Wild Arms 2", "Disc 1"},
	"SCUS_944.98": {"Wild Arms 2", "Disc 2"},
	"SLUS_000.19": {"Wing Commander III - Heart of the Tiger", "Disc 1"},
	"SLUS_001.34": {"Wing Commander III - Heart of the Tiger", "Disc 2"},
	"SLUS_001.35": {"Wing Commander III - Heart of the Tiger", "Disc 3"},
	"SLUS_001.36": {"Wing Commander III - Heart of the Tiger", "Disc 4"},
	"SLUS_002.70": {"Wing Commander IV - The Price of Freedom", "Disc 1"},
	"SLUS_002.71": {"Wing Commander IV - The Price of Freedom", "Disc 2"},
	"SLUS_002.72": {"Wing Commander IV - The Price of Freedom", "Disc 3"},
	"SLUS_002.73": {"Wing Commander IV - The Price of Freedom", "Disc 4"},
	"SLUS_009.15": {"The X-Files", "Disc 1"},
	"SLUS_009.49": {"The X-Files", "Disc 2"},
	"SLUS_009.50": {"The X-Files", "Disc 3"},
	"SLUS_009.51": {"The X-Files", "Disc 4"},
	"SLUS_006.64": {"Xenogears", "Disc 1"},
	"SLUS_006.69": {"Xenogears", "Disc 2"},
	"SLUS_007.16": {"You Don't Know Jack", "Disc 1"},
	"SLUS_007.62": {"You Don't Know Jack", "Disc 2"},
}

// Multidisc moves all the multidisc files for a game into the same directory
func Multidisc(dir string) error {
	hasCue, err := files.HasFilesWithExtension(dir, files.CueExtension)
	if err != nil {
		return err
	}

	if hasCue {
		return errors.New("can't call multidisc with cue files, first merge your bin files and generate cu2 sheets")
	}

	fmt.Println("Looking for multidisc games")
	bins, err := files.GetFilesByExtension(dir, files.BinExtension)
	if err != nil {
		log.Fatal(err)
	}

	createdDirs := make(map[string][]string)
	for _, bin := range bins {
		serial, err := psxserialnumber.GetSerial(bin)
		if err != nil {
			continue
		}

		if disc, ok := multidiscMap[serial]; ok {
			destination := filepath.Join(dir, disc.name)
			err = os.MkdirAll(destination, os.ModePerm)

			binName := fmt.Sprintf("%s.bin", disc.disc)
			createdDirs[destination] = append(createdDirs[destination], binName)
			err = os.Rename(bin, filepath.Join(destination, binName))
			if err != nil {
				return err
			}
			err = os.Rename(strings.Replace(bin, ".bin", ".cu2", 1), filepath.Join(destination, fmt.Sprintf("%s.cu2", disc.disc)))
			if err != nil {
				return err
			}

			files, _ := ioutil.ReadDir(filepath.Dir(bin))
			if len(files) == 0 {
				os.Remove(filepath.Dir(bin))
			}
		}
	}

	fmt.Printf("%d multidisc games found. creating multidisc.lst.\n", len(createdDirs))

	for d, bins := range createdDirs {
		f, err := os.Create(filepath.Join(d, "multidisc.lst"))
		if err != nil {
			return err
		}
		defer f.Close()

		for _, bin := range bins {
			f.WriteString(bin + "\n")
		}
	}

	fmt.Print("\n\n")

	return nil
}
