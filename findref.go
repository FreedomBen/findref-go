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

const Debug = false

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
		fmt.Println(Green + message + Restore)
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
		debug(LightRed + "Processing " + Restore + Green + "file: " + Restore + Purple + path)

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
					debug("Not processing binary file: " + path)
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
		debug("Ignoring file cause it doesn't match filter: " + path)
	}
	// TODO:  Return error if relevant
	return nil
}

func getMatchRegex(matchCase bool, usersRegex string) *regexp.Regexp {
	// If -match-case is not set, figure out smartcase
	if !matchCase && !regexp.MustCompile("[A-Z]").MatchString(usersRegex) {
		return regexp.MustCompile("(?i)" + usersRegex) // make regex case insensitive
	} else {
		return regexp.MustCompile(usersRegex)
	}
}

func main() {
	matchCasePtr := flag.Bool("match-case", false, "Match regex case (if unset smart-case is used)")

	flag.Parse()

	// add --ignore-case
	// add mc for alias match-case
	// add ic for alias ignore-case

	fmt.Println("matchCase: ", *matchCasePtr)
	fmt.Println("tail: ", flag.Args())

	rootDir := "."

	if len(flag.Args()) < 1 {
		fmt.Println("Must specify regex to match against files")
		usageAndExit()
	} else if len(flag.Args()) > 3 {
		fmt.Println("Too many args")
		usageAndExit()
	} else {
		matchRegex = getMatchRegex(*matchCasePtr, flag.Args()[0])

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
