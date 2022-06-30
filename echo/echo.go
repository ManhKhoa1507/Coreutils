package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	help_text string = `
    Usage: echo [OPTION]... [STRING]...
       or: echo [OPTION]
    
    display a line of text
        
        -n     do not output the trailing newline
        -e     enable interpretation of backslash escapes
        -E     disable interpretation of backslash escapes (default)
        
        --help        display this help and exit
        --version     output version information and exit
    `
	version_text = "echo (go-coreutils) 0.1"
)

var (
	enableEscapeChars  = flag.Bool("e", false, "enable interpretation of blackslas escapes")
	omitNewLine        = flag.Bool("n", false, "do not output the trailling newline")
	disableEscapeChars = flag.Bool("E", true, "disable interpretation of blackslas escapes")
	help               = flag.Bool("h", false, help_text)
	version            = flag.Bool("v", false, version_text)
)

func main() {
	flag.Parse()

	// Switch case the options
	switch {

	// Help option -h
	case *help:
		fmt.Println(help_text)

	// Print version -v
	case *version:
		fmt.Println(version_text)
	}

	EchoString()

	if !*omitNewLine {
		fmt.Println("\n")
	}
}

func EchoString() {
	content := strings.Join(flag.Args(), " ")

	// Convert the concatenated to rune type
	runeContent := []rune(content)
	fmt.Println(runeContent)

	// calc the runeContent's len
	lenRune := len(runeContent)

	special_index := 0

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
				if '9' >= character && character >= '0' && i < lenRune {
					hex := (character - '0')
					character = runeContent[i]
					i++
					if '9' >= character && character >= '0' && i <= lenRune {
						character = 16*(character-'0') + hex
					}
				}
			}
		}
		runeContent[special_index] = character
		special_index++
	}

	fmt.Println(string(content[:special_index]))

}
