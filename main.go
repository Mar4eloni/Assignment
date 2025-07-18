package main

import (
	"assignment/parser"
	"bufio"
	"fmt"
	"log"
	"os"
)

// usage prints the correct command usage to stderr
func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <filename>\n", os.Args[0])
	fmt.Println(os.Stderr, "<Filename> Path to file containing email addresses to parse")
}

// main is the entry point that handles:
// - Command-line argument validation
// - File operations
// - Email parsing
// - Output formatting
func main() {

	// Validate exactly one argument is provided (program name + filename)
	if len(os.Args) != 2 {
		usage()
		os.Exit(1) // Exit status 1 for argument errors
	}

	filename := os.Args[1]

	// for testing os.Open("testdata/emails.txt")
	file, err := os.Open(filename)

	if err != nil {
		log.Fatalf("Error openingfile: %v", err)
	}
	defer file.Close() // must ensure cosure otherwise problems when main exits :D

	scanner := bufio.NewScanner(file)
	var emails []parser.EmailParts

	for scanner.Scan() {
		line := scanner.Text()
		emails = append(emails, parser.ParseLine(line))

	}

	// Generate and output JSON
	jsonOutput, err := parser.FormatResultsToJson(emails)
	if err != nil {
		log.Fatalf("error formating JSON: %v", err)
	}

	fmt.Println(string(jsonOutput))

}
