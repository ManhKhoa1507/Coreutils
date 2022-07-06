package main

import (
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/ogier/pflag"
)

const (
	helpText = `
    Usage: mv [OPTION]... [PATH]... [PATH]
       or: mv [PATH] [PATH]
       or: mv [OPTION]
    move or rename files or directories
        --help        display this help and exit
        --version     output version information and exit
        -f, --force   remove existing destination files and never prompt the user
    `

	versionText = "mv (go-coreutils) 0.1"
)

var (
	// options -f --force
	forceEnable = flag.BoolP("force", "f", false, "force to move file")
	help        = flag.BoolP("help", "h", false, "help")
	version     = flag.BoolP("version", "v", false, "mv version")
)

func main() {
	flag.Parse()

	switch {
	case *forceEnable:
		*forceEnable = true
	case *help:
		fmt.Println(helpText)
	case *version:
		fmt.Println(versionText)
	}
	if *forceEnable {
		*forceEnable = true
	}

	file := flag.Args()

	CheckArguments(file)
}

func CheckArguments(files []string) {
	// Get the files[] len
	lenFiles := len(files)

	const (
		operandError     = "mv: missing file operand\nTry 'mv -help' for more information"
		destinationError = "mv: missing destination file operand after '%s'\nTry 'mv --help' for more information.\n"
	)

	// Case files len
	switch lenFiles {

	// Missing file
	case 0:
		fmt.Println(operandError)
		os.Exit(0)

	// Missing newLocation
	case 1:

		fmt.Printf(destinationError, files[0])
		os.Exit(0)

	// If more than 2 arguments
	default:
		// Get the new location (last args)
		newLocation := files[lenFiles-1]

		// Re assign files[], len of files[]
		files = files[:lenFiles-1]
		lenFiles = len(files)

		// Move files to newLocation
		for i := 0; i < lenFiles; i++ {
			MoveFile(files[i], newLocation)
		}

		os.Exit(0)

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

// Move origin file to new location
func MoveFile(originFile string, newLocation string) {

	const (
		notFoundError = "Not such a file or directory"
	)

	newLocationStatus := CheckFileExists(newLocation)

	switch {

	// No origin file to move
	case CheckFileExists(originFile) == nil:
		fmt.Println(notFoundError, originFile)
		os.Exit(1)

	// If destination file exists
	case newLocationStatus != nil:

		// Move file to directory
		if newLocationStatus.IsDir() {

			// Get the base of origin file
			base := filepath.Base(originFile)
			fileBase := newLocation + "/" + base
			baseStatus := CheckFileExists(fileBase)

			fmt.Println(fileBase)

			// If file is exist, and option -f is disable
			if baseStatus != nil && !*forceEnable {
				MoveWithConfirmation(originFile, fileBase)

			} else if baseStatus != nil && *forceEnable {
				// Force is enable
				TryMove(originFile, fileBase)

			} else if baseStatus == nil {
				// If base file not exist (move don't need permission)
				TryMove(originFile, fileBase)
			}

		} else {
			// If file is not directory
			MoveWithConfirmation(originFile, newLocation)
		}

	// If force is enable or newLocation is nil
	default:
		TryMove(originFile, newLocation)
	}
}

func TryMove(originFile string, newLocation string) {

	const (
		linkError    = "Link error"
		pathError    = "Path error"
		syscallError = "Syscall error"
	)

	// Rename file
	err := os.Rename(originFile, newLocation)

	// Handle error
	switch err.(type) {

	// Case link is error
	case *os.LinkError:
		fmt.Println(linkError)
		os.Exit(1)

	// Path error
	case *os.PathError:
		fmt.Println(pathError)
		os.Exit(1)

	// Syscall err
	case *os.SyscallError:
		fmt.Println(syscallError)
		os.Exit(1)
	}
}

// Move if have permission
func MoveWithConfirmation(originFile string, newLocation string) {
	// Get user permission
	answer := GetUserConfirmation(newLocation)

	if answer == "y" {
		// Have permission, try to move file
		TryMove(originFile, newLocation)

	} else {
		// Not permission
		os.Exit(0)
	}
}

// Get user permission to overwrite file
func GetUserConfirmation(file string) string {
	answer := ""

	fmt.Printf("File %s is exist. Overwrite(y/n)? ", file)
	fmt.Scanf("%s", &answer)

	return answer
}
