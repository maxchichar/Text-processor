package api

import (
	"text-processor/core"
)

// Processor wraps the text processing functionality
type Processor struct{}

// NewProcessor creates a new Processor instance
func NewProcessor() *Processor {
	return &Processor{}
}

// ProcessOptions defines which transformations to apply
type ProcessOptions struct {
	Tokenize    bool `json:"tokenize"`
	Hex         bool `json:"hex"`
	Bin         bool `json:"bin"`
	Case        bool `json:"case"`
	Quote       bool `json:"quote"`
	Punctuation bool `json:"punctuation"`
	Article     bool `json:"article"`
}

// DefaultProcessOptions returns all transformations enabled
func DefaultProcessOptions() ProcessOptions {
	return ProcessOptions{
		Tokenize:    true,
		Hex:         true,
		Bin:         true,
		Case:        true,
		Quote:       true,
		Punctuation: true,
		Article:     true,
	}
}

// ToCoreOptions converts API options to core package options
func (p *ProcessOptions) ToCoreOptions() core.TransformationOptions {
	return core.TransformationOptions{
		Tokenize:    p.Tokenize,
		Hex:         p.Hex,
		Bin:         p.Bin,
		Case:        p.Case,
		Quote:       p.Quote,
		Punctuation: p.Punctuation,
		Article:     p.Article,
	}
}

// ProcessFile processes a file with the given options
func (p *Processor) ProcessFile(inputPath, outputPath string, options ProcessOptions) error {
	coreOptions := options.ToCoreOptions()
	return core.ProcessFile(inputPath, outputPath, coreOptions)
}

// ProcessText processes a single line of text with the given options
func (p *Processor) ProcessText(line string, options ProcessOptions) string {
	coreOptions := options.ToCoreOptions()
	return core.ProcessText(line, coreOptions)
}