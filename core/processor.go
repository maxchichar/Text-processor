package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text-processor/internals"
)

// TransformationOptions defines which transformations to apply
type TransformationOptions struct {
	Tokenize    bool
	Hex         bool
	Bin         bool
	Case        bool
	Quote       bool
	Punctuation bool
	Article     bool
}

// DefaultTransformationOptions returns all transformations enabled
func DefaultTransformationOptions() TransformationOptions {
	return TransformationOptions{
		Tokenize:    true,
		Hex:         true,
		Bin:         true,
		Case:        true,
		Quote:       true,
		Punctuation: true,
		Article:     true,
	}
}

// ProcessText processes a single line of text with the given transformation options
func ProcessText(line string, options TransformationOptions) string {
	texts := []string{line}

	if options.Tokenize {
		texts = internals.ApplyTokenizer(line)
	}

	if options.Hex {
		texts = internals.ApplyHex(texts)
	}

	if options.Bin {
		texts = internals.ApplyBin(texts)
	}

	if options.Case {
		texts = internals.ApplyCases(texts)
	}

	if options.Quote {
		texts = internals.ApplyQuote(texts)
	}

	if options.Punctuation {
		texts = internals.ApplyPunctuation(texts)
	}

	if options.Article {
		texts = internals.ApplyAnRule(texts)
	}

	return strings.Join(texts, " ")
}

// ProcessFile processes an entire file with the given transformation options
func ProcessFile(inputPath, outputPath string, options TransformationOptions) error {
	inF, err := os.OpenFile(inputPath, os.O_RDONLY, 0664)
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}
	defer inF.Close()

	var b strings.Builder
	scanner := bufio.NewScanner(inF)

	for scanner.Scan() {
		line := scanner.Text()
		processed := ProcessText(line, options)
		b.WriteString(processed)
		b.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	outF, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outF.Close()

	_, err = outF.WriteString(b.String())
	if err != nil {
		return fmt.Errorf("error writing output file: %w", err)
	}

	return nil
}