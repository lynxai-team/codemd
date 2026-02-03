// codemd-lite is a lightweight tool to convert between markdown files
// containing code blocks and individual source code files.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	cmd := CommandMdCodeLite()

	if err := cmd.Execute(); err != nil {
		cmd.PrintErr(err)
		os.Exit(1)
	}
}

func CommandMdCodeLite() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "codemd",
		Short: "A tool to convert between markdown and source files",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	cmd.AddCommand(CommandToCode())
	cmd.AddCommand(CommandToMarkdown())

	return cmd
}

func CommandToCode() *cobra.Command {
	var mdFile string
	var outputDir string

	cmd := &cobra.Command{
		Use:     "tocode",
		Short:   "Convert markdown file to source files",
		Example: `codemd tocode -i docs.md -o src`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if mdFile == "" {
				return fmt.Errorf("markdown file is required (use --input or -i)")
			}
			if outputDir == "" {
				return fmt.Errorf("output directory is required (use --output or -o)")
			}

			doc := NewDocument()

			cmd.Printf("Parsing markdown file: %s\n", mdFile)
			if err := doc.ParseMarkdown(mdFile); err != nil {
				return fmt.Errorf("failed to parse markdown: %w", err)
			}

			cmd.Printf("Extracting %d code blocks to: %s\n", len(doc.Blocks), outputDir)
			if err := doc.ToSourceFiles(outputDir); err != nil {
				return fmt.Errorf("failed to extract files: %w", err)
			}

			cmd.Printf("Successfully extracted:")
			for _, block := range doc.Blocks {
				cmd.Printf("  - %s (%s)\n", block.Filename, block.Language)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&mdFile, "input", "i", "", "Input markdown file to extract from")
	_ = cmd.MarkFlagRequired("input")

	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for extracted files")
	_ = cmd.MarkFlagRequired("output")

	return cmd
}

func CommandToMarkdown() *cobra.Command {
	var outputFile string
	var directory string
	var files []string

	cmd := &cobra.Command{
		Use:     "tomd",
		Short:   "Convert source files to markdown",
		Example: `codemd tomd -f file1 -f file2 -o output.md`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if outputFile == "" {
				return fmt.Errorf("output file is required (use --output or -o)")
			}

			doc := NewDocument()

			if directory == "" && len(files) == 0 {
				return fmt.Errorf("must specify either --dir or --files (or both)")
			}

			if directory != "" {
				cmd.Printf("Scanning directory: %s\n", directory)
				if err := doc.ParseDir(directory); err != nil {
					return fmt.Errorf("failed to scan directory: %w", err)
				}
			}

			if len(files) > 0 {
				cmd.Printf("Parsing %d files\n", len(files))
				if err := doc.ParseFiles(files...); err != nil {
					return fmt.Errorf("failed to parse files: %w", err)
				}
			}

			cmd.Printf("Generating markdown with %d code blocks: %s\n", len(doc.Blocks), outputFile)
			if err := doc.ToMarkdown(outputFile); err != nil {
				return fmt.Errorf("failed to generate markdown: %w", err)
			}

			cmd.Printf("Successfully generated markdown with:")
			for _, block := range doc.Blocks {
				cmd.Printf("  - %s (%s)\n", block.Filename, block.Language)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output markdown file")
	_ = cmd.MarkFlagRequired("output")

	cmd.Flags().StringVarP(&directory, "dir", "d", "", "Directory to scan for source files")
	cmd.Flags().StringSliceVarP(&files, "files", "f", nil, "Files to include")

	return cmd
}
