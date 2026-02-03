# codemd

Convert between markdown files containing code blocks and individual source files.

## Installation

```bash
go install github.com/lynxai-team/codemd@latest
```

## Usage

### Convert source files to markdown

```bash
codemd tomd -f main.go -f helper.js -o docs.md
```

Or scan a directory:

```bash
codemd tomd -d ./src -o docs.md
```
