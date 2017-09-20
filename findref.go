package main

import "flag"
import "fmt"
import (
    "regexp"
    "os"
)

func listAllFiles() {
    fmt.Println("Listing all files")
}

func printMatches(filename string, regexp *regexp.Regexp) {
    fmt.Println("Printing matches")
}

func regex_start_filenames(args ...string) string {
    return "world"
}

func printUsage() {
    fmt.Println("Usage")
}

func usageAndExit() {
    printUsage()
    os.Exit(1)
}

func main() {
    matchCasePtr := flag.Bool("match-case", false, "Match regex case (if unset smart-case is used)")

    flag.Parse()

    fmt.Println("matchCase: ", *matchCasePtr)
    fmt.Println("tail: ", flag.Args())

    if len(flag.Args()) < 1 {
        fmt.Errorf("Must specify regex to match against files")
        usageAndExit()
    } else if len(flag.Args()) > 3 {
        fmt.Errorf("Too many args")
        usageAndExit()
    }

    //regex := flag.Args()
    listAllFiles();

    fmt.Println("args: ", os.Args)
    //printMatches();
}
