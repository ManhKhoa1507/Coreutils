package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	flag "github.com/ogier/pflag"
)

const (
	helpText = `
    Usage: head [OPTION]... [FILE]...
       
    Print the first 10 lines of each FILE to standard output. With more than one FILE, precede each with a header giving the file name. With no FILE, or when FILE is -, read standard input.
    
    
    Mandatory arguments to long options are mandatory for short options too.
       -help        display this help and exit
       -version     output version information and exit
       -c, --bytes=K
              output the first K bytes; or use -n +K to output starting with the Kth byte.
       -n, -lines=K
              output the first K lines; or use -n +K to output starting with 
the Kth
       -q, -quiet, -silent
              never output headers giving file names
`
	versionText = "head (go-coreutils) 0.1"
)

var (
	helpOption    = flag.BoolP("help", "h", false, "help option")
	versionOption = flag.BoolP("version", "v", false, "version")
	bytesOption   = flag.IntP("bytes", "c", 0, "bytes")
	linesOption   = flag.IntP("lines", "n", 10, "number of lines")
)

func main() {
	flag.Parse()

	switch {

	// case --help -h
	case *helpOption:
		fmt.Print(helpText)
		os.Exit(0)

	// case -v --version
	case *versionOption:
		fmt.Println(versionText)
		os.Exit(0)
	}

	// Get the files
	files := flag.Args()

	// Check for valid options
	// CheckOption()

	// Read heading of all files
	CheckFile(files)
}

// Check one or many file
func CheckFile(files []string) {
	lenFiles := len(files)
	switch {

	case lenFiles == 1:
		PrintOption(files[0])

	case lenFiles > 1:
		for _, file := range files {
			PrintFileName(file)
			PrintOption(file)
		}
	}

}

// Handle option -c or -n
func PrintOption(filePath string) {

	const (
		combineError = "head: can't combine line and byte counts"
		illegalError = "illegal arguments"
	)

	file := OpenFile(filePath)
	defer file.Close()

	switch {

	// Not valid lines or bytes
	case *bytesOption < 0 || *linesOption < 0:
		fmt.Println(illegalError)
		os.Exit(1)

	// Can't set 2 options -n -c at the same time
	case *bytesOption > 0 && *linesOption != 10:
		fmt.Println(combineError)
		os.Exit(1)

	// Option -c print bytes
	case *bytesOption > 0:
		PrintHeadingBytes(file)

	// Option -n print lines
	case *linesOption >= 0:
		PrintHeadingLine(file)

	// Default: Print first 10 lines
	default:
		PrintHeadingLine(file)
	}
}

// Print file name
func PrintFileName(file string) {
	fmt.Printf("============>%s<===============\n", file)
}

// Print number heading content of files
func PrintHeadingLine(file os.File) {

	// Get the fileReader io.Reader
	fileReader := &file
	reader := bufio.NewScanner(fileReader)
	lineCount := 0

	for lineCount < *linesOption && reader.Scan() {
		fmt.Println(string(reader.Bytes()))
		lineCount++
	}
}

// Print heading bytes
func PrintHeadingBytes(file os.File) {
	fileReader := &file
	reader := io.LimitReader(fileReader, int64(*bytesOption))
	io.Copy(os.Stdout, reader)
}

// Open specific file
func OpenFile(file string) os.File {
	fileTest, err := os.Open(file)

	// Handle error when open file
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return *fileTest
}
