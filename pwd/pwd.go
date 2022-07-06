package main

import (
	"fmt"
	"io/fs"
	"os"
	"syscall"

	flag "github.com/ogier/pflag"
)

const (
	helpText = `
	NAME
		   pwd - print name of current/working directory
	SYNOPSIS
		   pwd [OPTION]...
	DESCRIPTION
		   Print the full filename of the current working directory.
		   -L, --logical
				  use PWD from environment, even if it contains symlinks
		   -P, --physical
				  avoid all symlinks
		   --help display this help and exit
		   --version
				  output version information and exit
		   If no option is specified, -P is assumed.
	`
	versionText = "pwd (Go coreutils) 0.1"
)

// Get option from terminal using flag
var (
	// Option -L for logical path
	logical = flag.BoolP("logical", "L", false, "")

	// Option -P for physical path
	physical = flag.BoolP("physical", "P", false, "")

	// Option -P for version
	version = flag.BoolP("version", "V", false, "")
)

func main() {
	flag.Usage = func() {
		// Print help option
		fmt.Println(helpText)
		os.Exit(1)
	}

	flag.Parse()

	if len(flag.Args()) != 0 {
		fmt.Println("No arguments")
		os.Exit(1)
	}

	switch {

	// Case logical (-L) option
	case *logical:
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error when using option -L")
			os.Exit(1)
		}

		fmt.Println(dir)

	// Case physical (-P) option
	case *physical:
		dir, err := GetwdWithoutSymLinks()

		// Handle error when use option -P
		if err != nil {
			fmt.Println("Error when using option -P")
			os.Exit(1)
		}

		fmt.Println(dir)

	// Case version (-V) option
	case *version:
		fmt.Println(versionText)

	// Default is -L
	default:
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error when using option -L")
			os.Exit(1)
		}

		fmt.Println(dir)
	}
}

func GetwdWithoutSymLinks() (dir string, err error) {
	// Use runtime.GOOS to detect the OS system

	// if $PWD is set and matches with "./" (current path)
	// Lstat returns a FileInfo describing the named file.
	// If the file is a symbolic link, the returned FileInfo describes the symbolic link. Lstat makes no attempt to follow the link.
	// If there is an error, it will be of type *PathError.
	directory := ""

	currentDir, _ := GetFileLstat(".")

	pwd := os.Getenv("PWD")
	isPWD := CheckPWD(pwd, currentDir)

	// Check if pwd != "" and pwd is not set to another value
	if isPWD {
		directory, err = GetPWD(pwd, currentDir)
	} else {
		directory, err = GetDir(currentDir)
	}

	// Root case with ends is / and no parents

	return directory, err
}

// Get the $PWD if set
func GetPWD(pwdDir string, currentDir fs.FileInfo) (dir string, err error) {

	// Check $PWD is set and first character is /
	if len(pwdDir) > 0 && pwdDir[0] == '/' {

		// Check the Lstat of $PWD
		d, _ := GetFileLstat(pwdDir)

		// If $PWD is same with dir (./)
		if os.SameFile(currentDir, d) {
			return pwdDir, nil
		}
	}
	return "", nil
}

func CheckPWD(pwd string, currentDir fs.FileInfo) bool {

	if len(pwd) > 0 && pwd[0] == '/' {

		pwdDir, _ := GetFileLstat(pwd)

		if len(pwd) > 0 && os.SameFile(pwdDir, currentDir) {
			return true
		}
	}
	return false
}

func GetDir(currentDir fs.FileInfo) (dir string, err error) {

	// Declare directory path, root path
	directory := ""
	root, _ := GetFileLstat("/")

	// Declare and add parent
	for parent := ".."; ; parent = "../" + parent {

		// Check parent if path for too long
		if len(parent) >= 1024 {
			return "", syscall.ENAMETOOLONG
		}

		// Open Parent folder
		parentDirectory, _ := OpenFile(parent)
		defer parentDirectory.Close()

		// Get parent directory stat, files contains in directory
		parentStat, _ := GetFileStat(&parentDirectory)
		files, _ := GetDirName(&parentDirectory)

		GetParentDirectory(files, parent, directory, currentDir, parentStat)

		// Check if parent directory = root, then break the loop
		if os.SameFile(parentStat, root) {
			break
		}
	}

	return directory, err
}

// Get the fileStatus using os.Lstat
func GetFileLstat(filePath string) (fs.FileInfo, error) {
	file, err := os.Lstat(filePath)

	// Handle error when Lstat the file
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return file, err
}

// Open file
func OpenFile(filePath string) (os.File, error) {
	file, err := os.Open(filePath)

	// Handle error when open file
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return *file, err
}

// Get the dirname
func GetDirName(filePath *os.File) ([]string, error) {
	files, err := filePath.Readdirnames(100)

	// Handle error when read dir name
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return files, err
}

// Get the fileStatus using os.Lstat
func GetFileStat(filePath *os.File) (fs.FileInfo, error) {
	file, err := filePath.Stat()

	// Handle error when Lstat the file
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return file, err
}

// Check file in parent directory if similar to the files before
func GetParentDirectory(files []string, parent string, directory string, currentDir fs.FileInfo, parentStat fs.FileInfo) {
	for _, file := range files {

		filePath, _ := GetFileLstat(parent + "/" + file)

		if os.SameFile(filePath, currentDir) {
			// If same file with currentDirectory add name and directory path
			directory = "/" + file + directory
			currentDir = parentStat
		}
	}
}
