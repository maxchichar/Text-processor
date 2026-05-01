// By MAXCHICHAR on April 18th 2026

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text-processor/core"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("✗ Usage: go run . sample.txt result.txt")
		os.Exit(1)
	}

	inFile := os.Args[1]
	outFile := os.Args[2]

	if (!strings.HasSuffix(inFile, ".txt")) || (!strings.HasSuffix(outFile, ".txt")) {
		fmt.Println("✗ Error: File must end with .txt")
		os.Exit(1)
	}

	if inFile == outFile {
		fmt.Println("✗ Error: Both file name cannot be the same")
		os.Exit(1)
	}

	options := core.DefaultTransformationOptions()
	err := core.ProcessFile(inFile, outFile, options)
	if err != nil {
		log.Fatal("✗ Error:", err)
	}

	fmt.Println("TEXT PROCESSED")
}