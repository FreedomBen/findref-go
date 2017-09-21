package main

import "flag"
import "fmt"
import (
	"os"
	"path/filepath"
	"regexp"
	//"errors"
	//"net"
)

const LightRed = "\033[1;31m"
const Green = "\033[0;32m"
const Purple = "\033[0;35m"
const Restore = "\033[0m"

var matchRegex *regexp.Regexp = nil
var fileFilter *regexp.Regexp = regexp.MustCompile(".*")

func printMatches(filename string, regexp *regexp.Regexp) {
	fmt.Println("Printing matches")
}

func printUsage() {
	fmt.Println("Usage")
}

func usageAndExit() {
	printUsage()
	os.Exit(1)
}

func printMatch(path string, lineNumber string, line string) {
	fmt.Println(Purple + path + Restore + Green + ":" + lineNumber + ":" + Restore + LightRed + line + Restore)
}

func passesFileFilter(path string) bool {
	return fileFilter.MatchString(path)
}

func processFile(path string, info os.FileInfo, err error) error {
	if passesFileFilter(path) {
		// TODO: read file and match against matchRegex
		// if lineMatches {
		//     printMatch(path, lineNumber, line)
		// }
		fmt.Println(LightRed + "Processing " + Restore + Green + "file: " + Restore + Purple + path)
	} else {
		fmt.Println("Ignoring file cause it doesn't match filter: ", path)
	}
	return nil
	//return errors.New(nil)
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
