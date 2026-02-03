package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarkdownToSourceFiles(t *testing.T) {
	t.Parallel()

	doc := NewDocument()
	tempDir := t.TempDir()

	require.NoError(t, doc.ParseMarkdown("testdata/golden.md"))
	require.NoError(t, doc.ToSourceFiles(tempDir))

	filePairs := []struct {
		actual string
		wanted string
	}{
		{"main.c", "testdata/code/main.c"},
		{"helper.js", "testdata/code/helper.js"},
		{"file3.css", "testdata/code/main.css"},
	}

	for _, cf := range filePairs {
		assertEqualFile(
			t,
			filepath.Join(tempDir, cf.actual),
			cf.wanted,
		)
	}
}

func TestSourceFilesToMarkdown(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		parseFunc func(doc *Document) error
	}{
		"individual files": {
			parseFunc: func(doc *Document) error {
				return doc.ParseFiles(
					"testdata/code/helper.js",
					"testdata/code/main.c",
					"testdata/code/main.css",
				)
			},
		},
		"directory": {
			parseFunc: func(doc *Document) error {
				return doc.ParseDir("testdata/code")
			},
		},
	}

	for name, cc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			doc := NewDocument()

			require.NoError(t, cc.parseFunc(doc))

			outputPath := filepath.Join(t.TempDir(), "generated.md")

			require.NoError(t, doc.ToMarkdown(outputPath))
			assertEqualFile(t, outputPath, "testdata/golden-code.md")
		})
	}
}

func assertEqualFile(t *testing.T, actualPath, expectedPath string) {
	t.Helper()

	actualContent, err := os.ReadFile(actualPath)
	require.NoError(t, err, "read file %s", actualPath)

	expectedContent, err := os.ReadFile(expectedPath)
	require.NoError(t, err, "read file %s", expectedPath)

	assert.Equal(t, string(expectedContent), string(actualContent))
}
