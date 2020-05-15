# PSIO Helper

PSIOHelper is a cross platform standalone program written in GoLang to prepare cue/bin files for the PSIO cardtrige. It runs in Windows, Linux and MacOS.

## Features
Merge multibin files into single bin

Convert CUE sheets to CU2 sheets

Download game covers automatically

Shorten file names longer than 60 characters 

## Download
[Windows](https://github.com/ncirocco/psio-helper/releases/download/v0.3.0/psioHelper.exe)

[Linux](https://github.com/ncirocco/psio-helper/releases/download/v0.3.0/psioHelperLinux)

[Mac](https://github.com/ncirocco/psio-helper/releases/download/v0.3.0/psioHelperMac)

## Usage

### Windows
For the most basic use, place the downloaded file in the same folder that your cue/bin files are and execute the program. This will create a folder called `output` and place each game in its own folder, merge its bin files, create a CU2 sheet and download the cover image.

For an advance use open the power shell in the directory where the PSIOHelper binary is placed and execute the program passing any of the available [commands](https://github.com/ncirocco/psio-helper/blob/master/README.md#commands).

### Linux/MacOS
Open a terminal, navigate to the location where the PSIOHelper binary is placed and execute it passing any of the available [commands](https://github.com/ncirocco/psio-helper/blob/master/README.md#commands).

## Commands
**Windows note:** In windows you will have to write `.\psioHelper.exe` instad of `./psioHelper` for this commands to run in shell.

**Linux/MacOS note:** If you downloaded the Linux or Mac version of the binary, rename it from `psioHelperLinux`/`psioHelperMac` to `psioHelper` for the following commands to work, or call them using `./psioHelperLinux`/`./psioHelperMac` as base.

`./psioHelper -h` displays all the available commands

`./psioHelper <command> -h` displays help for the given command

### Auto
`./psioHelper auto` - Merges bin files and generates cu2 sheets for all the files in the given directory. It's the default behavior if no command and arguments are passed (same that just double click the executable).

#### Arguments
`dir`: Directory containing the files to be processed. By default uses the current directory.

`destinationDir`: Directory to store the processed files, if it doesn't exists it gets created. If this argument is not passed it creates an `output` folder in the current directory.

#### Example
`./psioHelper auto -dir="MyISOsFolder" -destinationDir="newEmptyFolder"`

### Merge
`./psioHelper merge` - Merges all the bin files in a given directory. 

#### Arguments
`dir`: Directory containing the files to be processed. By default uses the current directory.

`destinationDir`: Directory to store the processed files. By default it creates an `output` folder in the current directory.

#### Example
`./psioHelper merge -dir="MyISOsFolder" -destinationDir="newEmptyFolder"`

### Cu2
`./psioHelper cu2` - Generates the cu2 sheet for each cue sheet in the given directory.

#### Arguments
`dir`: Directory containing the files to be processed. By default uses the current directory.

`removeCue`: If passed the original cue files will be removed

#### Example
`./psioHelper cu2 -dir="MyISOsFolder" -removeCue`

### Images
`./psioHelper images` - Downloads covers for the bin files in the given directory.

#### Arguments
`dir`: "Directory containing the bin files to get the images. By default uses the current directory.

#### Example
`./psioHelper` images -dir="MyISOsFolder"

### Multidisc
`./psioHelper multidisc` - Groups the discs that belong to the same game and generates the necessary multidisc.lst files

#### Arguments
`dir`: "Directory containing the bin files. By default uses the current directory.

**IMPORANT:** This command will modify your files without creating a backup.

#### Example
`./psioHelper` multidisc -dir="MyISOsFolder"


## Development
### Pre requirements
Have Go installed and configured in your local environment. More information can be found here https://golang.org/doc/install

### Download and usage
Run the following commands

`go get github.com/ncirocco/psio-helper`

`cd $GOPATH/src/github.com/ncirocco/psio-helper/`

To run the application use

`go run main.go [command] [arguments]`

To build the aplication use

`Make build`

The Windows, Linux and MacOS executables will be placed in the `build` folder

## Disclaimer
This software comes with no warranties of any kind whatsoever. Use it at your own risk.
