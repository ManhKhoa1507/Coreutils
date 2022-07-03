package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	helpText = `
    Usage: echo [OPTION]... [STRING]...
       or: echo [OPTION]
    
    display a line of text
        
        -n     do not output the trailing newline
        -e     enable interpretation of backslash escapes
        -E     disable interpretation of backslash escapes (default)
        
        --help        display this help and exit
        --version     output version information and exit
    `
	versionText = "echo (go-coreutils) 0.1"
)

var (
	enableEscapeChars  = flag.Bool("e", false, "enable interpretation of blackslas escapes")
	omitNewLine        = flag.Bool("n", false, "do not output the trailling newline")
	disableEscapeChars = flag.Bool("E", true, "disable interpretation of blackslas escapes")
	help               = flag.Bool("h", false, helpText)
	version            = flag.Bool("v", false, versionText)
)

func main() {
	flag.Parse()

	// Switch case the options
	switch {

	// Help option -h
	case *help:
		fmt.Println(helpText)

	// Print version -v
	case *version:
		fmt.Println(versionText)
	}

	EchoString()

	// Option -n no new line
	if !*omitNewLine {
		fmt.Println("\n")
	}
}

func EchoString() {
	content := strings.Join(flag.Args(), " ")

	// Convert the concatenated to rune type
	runeContent := []rune(content)

	specialContent := BackSlasEscape(runeContent)
	fmt.Println(specialContent)
}

// Handle the blackslas
func BackSlasEscape(runeContent []rune) string {
	specialIndex := 0
	specialContent := ""
	lenRune := len(runeContent)

	for i := 0; i < lenRune; {

		// get the character
		character := runeContent[i]
		i++

		// check option if enable/disable escapeChars and special character
		if (*enableEscapeChars || !(*disableEscapeChars)) && character == '\\' {

			// get the character
			character = runeContent[i]

			switch character {

			// If speacial is '\a' -> '\\'
			case 'a':
				character = '\a'

			case 'b':
				character = '\b'

			case 'c':
				os.Exit(0)

			case 'e':
				character = '\x1B'

			case 'f':
				character = '\f'

			case 'n':
				character = '\n'

			case 'r':
				character = '\r'

			case 't':
				character = '\t'

			case 'v':
				character = '\v'

			case '\\':
				character = '\\'

			// Convert to ascii format
			case 'x':
				i++
				character = ConvertToASCII(runeContent, character, i, lenRune)
			}
		}
		runeContent[specialIndex] = character
		specialIndex++
	}

	specialContent = string(runeContent[:specialIndex])
	return specialContent
}

// Convert to ASCII from Hex
func ConvertToASCII(runeContent []rune, character rune, index int, lenRune int) rune {

	if '9' >= character && character >= '0' && index < lenRune {
		hex := (character - '0')
		character = runeContent[index]
		index++

		// Convert back to ASCII
		if '9' >= character && character >= '0' && index <= lenRune {
			character = 16*(character-'0') + hex
		}
	}
	return character
}
