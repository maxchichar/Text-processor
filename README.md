# Text-processor

A robust, command-line text transformation tool written in **Go** that cleans and normalizes messy text with smart rules inspired by classic text processing challenges.

It intelligently handles punctuation spacing, quote normalization, number base conversions, case transformations, and English article rules, turning chaotic input into clean, readable output.

## Features

- **Tokenizer** — Properly separates attached flags (e.g. `he(up)`, `1E(hex)`)
- **Case Transformations** — Supports `(up)`, `(low)`, `(cap)`, and numbered variants like `(cap, 4)` or `(up, 2)`
- **Number Conversion** — `(hex)` and `(bin)` to decimal
- **Quote Cleaning** — Handles mixed single/double quotes with extra spaces (e.g. `'   " chibueze '   "` → `'"chibueze"'`)
- **Punctuation Normalization** — Fixes spacing, ellipses, multiple marks, and commas
- **Smart "a/an" Rule** — Corrects articles before vowels and common exceptions
- Clean pipeline architecture with proper transformation order
- MIT Licensed

## Installation & Usage

### Prerequisites
- Go 1.21+

### Clone & Run

git clone https://github.com/maxchichar/Text-processor.git
cd Text-processor

# Build
go build -o text-processor

# Run
./text-processor sample.txt output.txt

Or run directly:bash

go run main.go sample.txt output.txt

ExampleInput (sample.txt):text

A elephant walked by. he(up) is a great programmer.
1E(hex) in decimal is 30. Test: '   " chibueze '   "

Output:text

An elephant walked by. HE is a great programmer.
30 in decimal is 30. Test: '"chibueze"'

Project Structurebash
```
Text-processor/
├── main.go                 # Entry point & processing pipeline
├── go.mod
├── internals/              # Core transformation logic
│   ├── tokenizer.go
│   ├── cases.go
│   ├── hex.go
│   ├── bin.go
│   ├── quote.go
│   ├── punctuation.go
│   └── anrule.go
├── sample.txt              # Sample input
├── expected_result.txt     # Expected clean output
├── .gitignore
└── LICENSE
````
Transformation OrderThe processor follows this strict order for accurate results:Tokenization (separates flags)
Hex & Bin conversion
Case transformations (up/low/cap)
Quote cleaning
Punctuation normalization
"a/an" rule correction

RoadmapSupport for stdin/stdout
More advanced flags ((low, N), (cap, N))
Configuration file support
Comprehensive unit tests
CLI flags for selective transformations

LicenseThis project is licensed under the MIT License — see the LICENSE file for details.
**Built with  in Go**


Made for learning and mastering clean text processing. Star the repo if you find it useful!
