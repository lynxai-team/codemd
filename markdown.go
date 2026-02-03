package main

import (
	"path/filepath"
	"strings"
)

var extToLanguage = map[string]string{
	// Programming languages
	".go":    "go",
	".js":    "javascript",
	".ts":    "typescript",
	".py":    "python",
	".java":  "java",
	".c":     "c",
	".cpp":   "cpp",
	".cxx":   "cpp",
	".cs":    "csharp",
	".php":   "php",
	".rb":    "ruby",
	".rs":    "rust",
	".swift": "swift",
	".kt":    "kotlin",
	".scala": "scala",
	".pl":    "perl",
	".lua":   "lua",
	".r":     "r",

	// Web technologies
	".html": "html",
	".htm":  "html",
	".css":  "css",
	".xml":  "xml",
	".json": "json",
	".yml":  "yaml",
	".yaml": "yaml",
	".toml": "toml",

	// Shell/scripting
	".sh":   "bash",
	".ps1":  "powershell",
	".bat":  "batch",
	".cmd":  "batch",
	".fish": "fish",
	".zsh":  "zsh",

	// Database
	".sql": "sql",

	// Markup/config
	".md":         "markdown",
	".tex":        "latex",
	".ini":        "ini",
	".properties": "properties",
	".dockerfile": "dockerfile",
	".makefile":   "makefile",
	".gitignore":  "gitignore",

	// Data formats
	".csv":   "csv",
	".jsonl": "jsonl",
	".tsv":   "tsv",

	// Other formats
	".txt":  "text",
	".diff": "diff",
	".log":  "log",
	".conf": "conf",
}

var languageToExt = map[string]string{
	// Programming languages
	"go":         ".go",
	"javascript": ".js",
	"js":         ".js",
	"typescript": ".ts",
	"ts":         ".ts",
	"python":     ".py",
	"py":         ".py",
	"java":       ".java",
	"c":          ".c",
	"cpp":        ".cpp",
	"c++":        ".cpp",
	"cxx":        ".cpp",
	"csharp":     ".cs",
	"c#":         ".cs",
	"cs":         ".cs",
	"php":        ".php",
	"ruby":       ".rb",
	"rb":         ".rb",
	"rust":       ".rs",
	"rs":         ".rs",
	"swift":      ".swift",
	"kotlin":     ".kt",
	"kt":         ".kt",
	"scala":      ".scala",
	"perl":       ".pl",
	"lua":        ".lua",
	"r":          ".r",

	// Web technologies
	"html": ".html",
	"css":  ".css",
	"xml":  ".xml",
	"json": ".json",
	"yaml": ".yml",
	"yml":  ".yml",
	"toml": ".toml",

	// Shell/scripting
	"bash":       ".sh",
	"shell":      ".sh",
	"sh":         ".sh",
	"powershell": ".ps1",
	"ps1":        ".ps1",
	"batch":      ".bat",
	"cmd":        ".bat",
	"bat":        ".bat",
	"fish":       ".fish",
	"zsh":        ".zsh",

	// Database
	"sql":        ".sql",
	"mysql":      ".sql",
	"postgresql": ".sql",
	"postgres":   ".sql",
	"sqlite":     ".sql",

	// Markup/config
	"markdown":   ".md",
	"md":         ".md",
	"latex":      ".tex",
	"tex":        ".tex",
	"ini":        ".ini",
	"properties": ".properties",
	"dockerfile": ".dockerfile",
	"docker":     ".dockerfile",
	"makefile":   ".makefile",
	"make":       ".makefile",
	"gitignore":  ".gitignore",

	// Data formats
	"csv":   ".csv",
	"jsonl": ".jsonl",
	"tsv":   ".tsv",

	// Other formats
	"text":   ".txt",
	"txt":    ".txt",
	"plain":  ".txt",
	"diff":   ".diff",
	"log":    ".log",
	"conf":   ".conf",
	"config": ".conf",
}

// detectSourceLanguage returns the programming language based on the filename extension.
func detectSourceLanguage(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if language, ok := extToLanguage[ext]; ok {
		return language
	}

	return ""
}

// determineFileExtension returns the appropriate file extension for a given programming language.
func determineFileExtension(language string) string {
	key := strings.ToLower(language)
	if ext, ok := languageToExt[key]; ok {
		return ext
	}

	return ""
}
