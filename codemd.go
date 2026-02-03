package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Document represents a collection of code blocks that can be converted
// between markdown and source files.
type Document struct {
	Blocks []CodeBlock
}

// CodeBlock represents a single code block with its metadata.
type CodeBlock struct {
	Filename string
	Language string
	Content  string
}

// NewDocument returns an empty Document.
func NewDocument() *Document {
	return &Document{
		Blocks: []CodeBlock{},
	}
}

// Reset clears all blocks from the document.
func (d *Document) Reset() {
	d.Blocks = []CodeBlock{}
}

// ParseMarkdown reads a markdown file and adds extracted code blocks to the document.
func (d *Document) ParseMarkdown(markdownPath string) error {
	data, err := os.ReadFile(markdownPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	lines := strings.Split(string(data), "\n")

	var inCodeBlock bool
	var language string
	var content strings.Builder
	var lastHeaderFilename string
	fileCounter := 0

	for _, line := range lines {
		if inCodeBlock {
			if strings.HasPrefix(strings.TrimSpace(line), "```") {
				fileCounter++

				d.Blocks = append(d.Blocks, CodeBlock{
					Filename: filename(lastHeaderFilename, language, fileCounter),
					Language: language,
					Content:  content.String(),
				})

				inCodeBlock = false
				language = ""
				content.Reset()
				lastHeaderFilename = ""
				continue
			}
			content.WriteString(line + "\n")
			continue
		}

		switch line := strings.TrimSpace(line); {
		case strings.HasPrefix(line, "```"):
			inCodeBlock = true
			language = strings.TrimSpace(line[3:])
		case strings.HasPrefix(line, "## "):
			headerText := strings.TrimSpace(line[3:])
			if headerText != "" {
				lastHeaderFilename = headerText
			}
		case line == "":
			continue
		default:
			lastHeaderFilename = ""
		}
	}

	return nil
}

// filename determines the appropriate filename for a code block.
// It prefers the header filename if its extension matches the language,
// otherwise generates a filename based on the language and counter.
func filename(headerFilename, language string, counter int) string {
	defaultName := fmt.Sprintf("file%d", counter)

	ext := determineFileExtension(language)
	if ext == "" {
		return defaultName
	}

	if headerFilename != "" && filepath.Ext(headerFilename) == ext {
		return headerFilename
	}

	return fmt.Sprintf("file%d%s", counter, ext)
}

// ParseFiles reads source files and adds them as blocks to the document.
func (d *Document) ParseFiles(filePaths ...string) error {
	for _, filePath := range filePaths {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", filePath)
		}

		filename := filepath.Base(filePath)
		language := detectSourceLanguage(filename)
		if language == "" {
			continue
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		block := CodeBlock{
			Filename: filePath,
			Language: language,
			Content:  string(content),
		}

		d.Blocks = append(d.Blocks, block)
	}

	return nil
}

// ParseDir recursively scans a directory for source files and adds them to the document.
func (d *Document) ParseDir(dir string) error {
	files, err := getFilesFromDirectory(dir)
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	sort.Strings(files)

	return d.ParseFiles(files...)
}

// ToMarkdown writes the Document as a markdown file.
func (d *Document) ToMarkdown(outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() { _ = file.Close() }()

	if _, err := file.WriteString("# Code Files\n\n"); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for i, block := range d.Blocks {
		if i > 0 {
			if _, err := file.WriteString("\n"); err != nil {
				return fmt.Errorf("failed to write spacing: %w", err)
			}
		}

		if _, err := fmt.Fprintf(file, "## %s\n\n", block.Filename); err != nil {
			return fmt.Errorf("failed to write filename header: %w", err)
		}

		if _, err := fmt.Fprintf(file, "```%s\n", block.Language); err != nil {
			return fmt.Errorf("failed to write code block start: %w", err)
		}

		if _, err := fmt.Fprintf(file, "%s```\n", block.Content); err != nil {
			return fmt.Errorf("failed to write code block: %w", err)
		}
	}

	return nil
}

// ToSourceFiles writes the Document as individual source files.
func (d *Document) ToSourceFiles(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for _, block := range d.Blocks {
		fullPath := filepath.Join(outputDir, block.Filename)

		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		if err := os.WriteFile(fullPath, []byte(block.Content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", fullPath, err)
		}
	}

	return nil
}

// defaultSkipPatterns contains patterns for paths to skip during directory scanning
var defaultSkipPatterns = []string{
	".git/",
	".gitignore",
	".gitattributes",
	".svn/",
	".hg/",

	"node_modules/",
	"vendor/",
	".venv/",
	"venv/",
	"env/",
	"__pycache__/",

	"dist/",
	"build/",
	"target/",
	"bin/",
	"obj/",
	".next/",
	"out/",
	"coverage/",

	".vscode/",
	".idea/",
	"*.swp",
	"*.swo",
	".DS_Store",
	"Thumbs.db",

	".pytest_cache/",
	".tox/",
	"*.egg-info/",
	".mypy_cache/",
	".coverage",

	".env",
	".env.*",
	".npmrc",

	"package-lock.json",
	"yarn.lock",
	"Cargo.lock",
	"pnpm-lock.yaml",
	"poetry.lock",
	"go.sum",

	"*.exe",
	"*.dll",
	"*.so",
	"*.dylib",
	"*.o",
	"*.obj",
	"*.a",
	"*.lib",

	"*.log",
	"*.tmp",
	"*.temp",

	".*",
}

// getFilesFromDirectory returns all source files in a directory recursively,
// skipping common non-source directories and files.
func getFilesFromDirectory(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if shouldSkip(path, info.IsDir(), defaultSkipPatterns) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// shouldSkip checks if a path matches any skip pattern
func shouldSkip(path string, isDir bool, patterns []string) bool {
	name := filepath.Base(path)

	for _, pattern := range patterns {
		if strings.HasSuffix(pattern, "/") {
			if !isDir {
				continue
			}
			dirPattern := strings.TrimSuffix(pattern, "/")
			if matched, _ := filepath.Match(dirPattern, name); matched {
				return true
			}
		} else {
			if matched, _ := filepath.Match(pattern, name); matched {
				if pattern == ".*" && name == "." {
					continue
				}
				return true
			}
		}
	}
	return false
}
