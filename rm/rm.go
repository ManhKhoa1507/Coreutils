package main

import (
	"fmt"
	"io"
	"os"

	flag "github.com/ogier/pflag"
)

const (
	helpText = `
    Usage: rm [OPTION]...
    
    remove files (delete/unlink)
        --help     display this help and exit
        --version  output version information and exit
        -f         ignore if files do not exist, never prompt
        -i         prompt before each removal
        -r, -R, --recursive
            remove directories and their contents recursively
    `
	versionText = "rm (go-coreutils) 0.1"
)

var (
	helpOption        = flag.BoolP("help", "h", false, "help option")
	versionOption     = flag.BoolP("version", "v", false, "version")
	forceEnableOption = flag.BoolP("force", "f", false, "forces")
	interactiveOption = flag.BoolP("interactive", "i", false, "interactive")
	recursiveOption   = flag.BoolP("recursive", "r", false, "remove with recursive")
	directoryOption   = flag.BoolP("directory", "d", false, "move directory")
)

func main() {
	flag.Parse()

	// switch with option
	switch {

	// Case option -h --help
	case *helpOption:
		fmt.Println(helpText)
		os.Exit(0)

	// Case option -v --version
	case *versionOption:
		fmt.Println(versionText)
		os.Exit(0)
	}

	// get files
	files := flag.Args()

	CheckArguments(files)

	// Remove file in files
	for i := 0; i < len(files); i++ {
		RemoveFile(files[i])
	}
}

// Check arguments
func CheckArguments(files []string) {
	// Get the files[] len
	lenFiles := len(files)

	const (
		operandError     = "rm: missing file operand\nTry 'rm --help' for more information"
		destinationError = "rm: missing destination file operand after '%s'\nTry 'rm --help' for more information.\n"
	)

	// Case files len == 0 : Missing file
	if lenFiles == 0 {
		fmt.Println(operandError)
		os.Exit(0)
	}
}

// Remove file
func RemoveFile(file string) {

	const (
		notFoundError     = "Not such a file or directory"
		notDirectory      = "Not a directory"
		directoryNotEmpty = "Directory not empty"
	)

	// Get the file status
	fileStatus := CheckFileExists(file)

	// If force -f option is enable
	if *forceEnableOption {
		fmt.Println("Force is enables")
		os.RemoveAll(file)
		return

	} else {
		// If -f is disable
		switch {

		// If not found file to remove
		case fileStatus == nil:
			fmt.Println(notFoundError, file)
			os.Exit(1)

		// If file is not directory
		case !fileStatus.IsDir():

			// If interactive mode is enable
			if *interactiveOption {
				RemoveWithInteractive(file)
			} else {

				// Not have interactive mode
				os.Remove(file)
			}

		// If file is directory, need to remove with option -d (empty folder) or -r
		case fileStatus.IsDir():
			if *recursiveOption {
				if *interactiveOption {

					// Get user permission to recursive
					answer := GetUserConfirmation(file)

					// If have permission
					if answer == "y" {
						RemoveWithRecursive(file)
					} else {
						// Without permission exit
						os.Exit(0)
					}
				} else {
					// Without -i option
					RemoveWithRecursive(file)
				}

			} else if *directoryOption {
				// if option is -d  (remove empty directory)
				// Check if directory is empty
				if CheckEmptyDirectory(file) {
					os.Remove(file)

				} else {
					// Not empty directory
					fmt.Println(directoryNotEmpty)
					os.Exit(1)
				}

			} else {
				// Can't remove directory without option
				fmt.Println(notDirectory)
				os.Exit(1)
			}
		}
	}
}

// Check if file exits using os.Stat(), return file info
func CheckFileExists(file string) os.FileInfo {
	// Get the file status
	fileStatus, err := os.Stat(file)

	// If file not exist and get error
	if err != nil && os.IsNotExist(err) {
		return nil
	}

	return fileStatus
}

// Get user permission to remove file
func GetUserConfirmation(file string) string {
	answer := ""

	fmt.Printf("Remove %s ? (y/n)", file)
	fmt.Scanf("%s", &answer)

	return answer
}

// Check if have interactive mode then remove file
func RemoveWithInteractive(file string) {
	// Get user permission to remove file
	answer := GetUserConfirmation(file)

	// If have permission to remove
	if answer == "y" {
		os.Remove(file)

	} else {
		// Not have permission
		return
	}
}

// Remove with recursive
func RemoveWithRecursive(filePath string) {

	// Open directory
	directory := OpenDirectory(filePath)

	// Close directory
	defer directory.Close()

	// Move directory
	defer os.Remove(filePath)

	// Open all file in directory then delete
	for {

		// Get file name
		fileNames, err := directory.Readdirnames(100)

		// Handle err when open read dir
		if err == io.EOF || len(fileNames) == 0 {
			break
		}

		// open files in directory
		for _, name := range fileNames {
			filePath := filePath + string(os.PathSeparator) + name
			RemoveFile(filePath)
		}
	}
}

func CheckEmptyDirectory(directory string) bool {
	directoryStatus := OpenDirectory(directory)
	defer directoryStatus.Close()

	_, err := directoryStatus.Readdirnames(100)

	if err == io.EOF {
		return true
	} else {
		return false
	}
}

// Open directory then get status of directory
func OpenDirectory(directory string) *os.File {
	directoryStatus, err := os.Open(directory)

	// handle error when open file
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return directoryStatus
}
