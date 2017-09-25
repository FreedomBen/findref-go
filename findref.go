package main

import "flag"
import "fmt"
import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

const LightRed = "\033[1;31m"
const Green = "\033[0;32m"
const Purple = "\033[0;35m"
const Restore = "\033[0m"

var Debug bool = false
var IncludeHidden bool = false

var hiddenFileRegex *regexp.Regexp = regexp.MustCompile(`(^|\/)\.`)
var matchRegex *regexp.Regexp = nil
var fileFilter *regexp.Regexp = regexp.MustCompile(".*")

func printUsage() {
	fmt.Println("Usage")
}

func usageAndExit() {
	printUsage()
	os.Exit(1)
}

func debug(message string) {
	if Debug {
		fmt.Println(message)
	}
}

func printMatch(path string, lineNumber int, line []byte, match []int) {
	// TODO:  Use match for coloring between match[0] and match[1]
	fmt.Println(Purple + path + Restore + Green + ":" + strconv.Itoa(lineNumber) + ":" + Restore + string(line[:match[0]]) + LightRed + string(line[match[0]:match[1]]) + Restore + string(line[match[1]:]))
}

func passesFileFilter(path string) bool {
	return fileFilter.MatchString(path)
}

func processFile(path string, info os.FileInfo, err error) error {
	if passesFileFilter(path) {
		debug(Green + "Processing file: " + Restore + path)

		// Ignore hidden files unless the IncludeHidden flag is set
		if !IncludeHidden && hiddenFileRegex.MatchString(path) {
		    debug(Green + "Hidden file '" + Restore + path + Green + "' not processed")
		    return nil
		}

		file, err := os.Open(path)
		// TODO: Handle err
		if err != nil {
			//log.Fatal(err)
			//debug("Error 1: ", err)
			return nil
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var lineNumber int = 0
		for scanner.Scan() {
			line := scanner.Bytes()
			for _, el := range line {
				if el == 0 {
					// This is a binary file.  Skip it!
					debug(Green + "Not processing binary file: " + Restore + path)
					return nil
				}
			}
			lineNumber++
			matchIndex := matchRegex.FindIndex(line)
			if matchIndex != nil {
				// we have a match! loc == nil means no match so just ignore that case
				printMatch(path, lineNumber, line, matchIndex)
			}
		}

		// TODO: Handle err
		if err := scanner.Err(); err != nil {
			//debug("Error 2: ", err)
			return nil
		}
	} else {
		debug(Green + "Ignoring file cause it doesn't match filter: " + Restore + path)
	}
	// TODO:  Return error if relevant
	return nil
}

func getMatchRegex(ignoreCase bool, matchCase bool, usersRegex string) *regexp.Regexp {
	// If ignore case is set, ignore the case of the regex.
	// if match-case is not set, use smartcase which means if it's all lower case be case-insensitive,
	// but if there's capitals then be case-sensitive
	if ignoreCase || (!matchCase && !regexp.MustCompile("[A-Z]").MatchString(usersRegex)) {
		return regexp.MustCompile("(?i)" + usersRegex) // make regex case insensitive
	} else {
		return regexp.MustCompile(usersRegex)
	}
}

func determineMatchCase(matchCasePtr *bool, mcPtr *bool) {
	if *mcPtr {
		*matchCasePtr = true
	}
}

func determineIgnoreCase(ignoreCasePtr *bool, icPtr *bool) {
	if *icPtr {
		*ignoreCasePtr = true
	}
}

func main() {
	matchCasePtr := flag.Bool("match-case", false, "Match regex case (if unset smart-case is used)")
	mcPtr := flag.Bool("mc", false, "Alias for --match-case")
	ignoreCasePtr := flag.Bool("ignore-case", false, "Ignore case in regex (overrides smart-case)")
	icPtr := flag.Bool("ic", false, "Alias for --ignore-case")
	debugPtr := flag.Bool("debug", false, "Enable debug mode")
	hiddenPtr := flag.Bool("hidden", false, "Include hidden files and files in hidden directories")

	flag.Parse()

	// add --ignore-case
	// add mc for alias match-case
	// add ic for alias ignore-case
	// add --debug for turning on debug mode
	// add -h|--hidden for searching hidden directories

	determineMatchCase(matchCasePtr, mcPtr)
	determineIgnoreCase(ignoreCasePtr, icPtr)

	IncludeHidden = *hiddenPtr
	Debug = *debugPtr

	fmt.Println("matchCase: ", *matchCasePtr)
	fmt.Println("ignoreCase: ", *ignoreCasePtr)
	fmt.Println("hidden: ", IncludeHidden)
	fmt.Println("debug: ", Debug)
	fmt.Println("tail: ", flag.Args())

	rootDir := "."

	if len(flag.Args()) < 1 {
		fmt.Println("Must specify regex to match against files")
		usageAndExit()
	} else if len(flag.Args()) > 3 {
		fmt.Println("Too many args")
		usageAndExit()
	} else {
		matchRegex = getMatchRegex(*ignoreCasePtr, *matchCasePtr, flag.Args()[0])

		if len(flag.Args()) >= 2 {
			rootDir = flag.Args()[1]
		}
		if len(flag.Args()) == 3 {
			fileFilter = regexp.MustCompile(flag.Args()[2])
		}
	}

	// TODO: Switch to powerwalk for performance:  https://github.com/stretchr/powerwalk
	filepath.Walk(rootDir, processFile)
}
