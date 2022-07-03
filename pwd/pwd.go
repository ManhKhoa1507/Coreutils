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
	dot, err := os.Lstat(".")

	// Handle the error when Lstat
	if err != nil {
		return "", err
	}

	pwd := os.Getenv("PWD")
	isPWD := CheckPWD(pwd, dot)

	// Check if pwd != "" and pwd is not set to another value
	if isPWD {
		directory, err = GetPWD(pwd, dot)
	} else {
		directory, err = GetDir(dot)
	}

	// Root case with ends is / and no parents

	return directory, err
}

// Get the $PWD if set
func GetPWD(pwdDir string, dot fs.FileInfo) (dir string, err error) {

	// Check $PWD is set and first character is /
	if len(dir) > 0 && dir[0] == '/' {

		// Check the Lstat of $PWD
		d, err := os.Lstat(dir)

		// If $PWD is same with dir (./)
		if err == nil && os.SameFile(dot, d) {
			return dir, nil
		}
	}
	return "", nil
}

func CheckPWD(pwd string, dot fs.FileInfo) bool {

	if len(pwd) > 0 && pwd[0] == '/' {
		pwdDir, err := os.Lstat(pwd)

		// Handle error when Lstat pwd
		if err != nil {
			fmt.Println("Error when using Lstat")
		}

		if len(pwd) > 0 && os.SameFile(pwdDir, dot) && err == nil {
			return true
		}
	}
	return false
}

func GetDir(dot fs.FileInfo) (dir string, err error) {

	// Declare directory path, root path
	directory := ""
	root, err := os.Lstat("/")

	// Handle error when using Lstat
	if err != nil {
		return "", err
	}

	// Declare and add parent
	for parent := ".."; ; parent = "../" + parent {

		// Check parent if path for too long
		if len(parent) >= 1024 {
			return "", syscall.ENAMETOOLONG
		}

		// open Parent folder
		fatherDir, err := os.Open(parent)

		// Handle the error
		if err != nil {
			return "", err
		}

		for {
			names, err := fatherDir.Readdirnames(100)

			// Handle error when read dir name
			if err != nil {
				return "", err
			}

			for _, name := range names {

				d, _ := os.Lstat(parent + "/" + name)

				if os.SameFile(d, dot) {
					// If same file with dot add name and directory path
					directory = "/" + name + directory
					goto Found

					// fmt.Println(directory)
				}
			}
		}

	Found:
		parentDir, err := fatherDir.Stat()

		// Handle error when using Stat()
		if err != nil {
			return "", err
		}

		// Close father dir
		fatherDir.Close()

		// Check if parent directory = root, then break the loop
		if os.SameFile(parentDir, root) {
			break
		}

		dot = parentDir
	}
	return directory, err
}
