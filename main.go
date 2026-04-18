package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"text-processor/internals"
)

func main()  {
	if len(os.Args) != 3{
		fmt.Println("✗ Usage: go run . sample.txt result.txt")
		os.Exit(1)
	}

	inFile := os.Args[1]
	outFile := os.Args[2]

	if (!(strings.HasSuffix(inFile, ".txt"))) || (!(strings.HasSuffix(outFile, ".txt"))) {
		fmt.Println("✗ Error: File must end with .txt")
		os.Exit(1)
	}

	if inFile == outFile {
		fmt.Println("✗ Error: Both file name cannot be the same")
		os.Exit(1)
	}

	inF, err := os.OpenFile(inFile, os.O_RDONLY, 0664)
	if err != nil {
		log.Fatal("✗ Error: File not opening", err)
	}
	defer inF.Close()
	
	var b strings.Builder

	scanner := bufio.NewScanner(inF)

	for scanner.Scan(){
		line := scanner.Text()
		
		texts := internals.ApplyTokenizer(line)
		// fmt.Println("Tokenizer: ", texts)
		texts = internals.ApplyHex(texts)
		// fmt.Println("Hex: ", texts)
		texts = internals.ApplyBin(texts)
		// fmt.Println("Bin: ", texts)
		texts = internals.ApplyCases(texts)
		// fmt.Println("Cases: ", texts)
		texts = internals.ApplyQuote(texts)
		// fmt.Println("Quote: ", texts)
		texts = internals.ApplyPunctuation(texts)
		// fmt.Println("Punctuation: ", texts)
		texts = internals.ApplyAnRule(texts)
		// fmt.Println("An Rule: ", texts)

		processed := strings.Join(texts, " ")
		b.WriteString(processed)
		b.WriteString("\n")
	}

	ProcessedText := b.String()

	if err := scanner.Err(); err != nil {
		log.Fatal("✗ Error: Unable to read file", err)
	}
	outF, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		log.Fatal("✗ Error: File not opening", err)
	}
	defer outF.Close()

	_, err = outF.WriteString(ProcessedText)
	if err != nil {
		log.Fatal("✗ Error: Unable to write file")
	}

	fmt.Println("TEXT PROCESSED")
}