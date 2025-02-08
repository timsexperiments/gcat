// pkg/gcat/gcat.go
package gcat

import (
	"path/filepath"
	"strings"
	// go-git for git operations
	// Git objects
	// in-memory storage
)

// defaultLanguageMap maps file extensions to their corresponding languages or file types.
var defaultLanguageMap = map[string]string{
	// Programming Languages
	".go":     "Go",
	".py":     "Python",
	".js":     "JavaScript",
	".jsx":    "JavaScript (React)",
	".ts":     "TypeScript",
	".tsx":    "TypeScript (React)",
	".java":   "Java",
	".c":      "C",
	".cpp":    "C++",
	".cs":     "C#",
	".rb":     "Ruby",
	".rs":     "Rust",
	".zig":    "Zig",
	".gleam":  "Gleam",
	".ex":     "Elixir",
	".exs":    "Elixir",
	".erl":    "Erlang",
	".hrl":    "Erlang",
	".v":      "V",
	".mojo":   "Mojo",
	".lua":    "Lua",
	".r":      "R",
	".jl":     "Julia",
	".pl":     "Perl",
	".pm":     "Perl",
	".php":    "PHP",
	".swift":  "Swift",
	".kt":     "Kotlin",
	".kts":    "Kotlin Script",
	".ktm":    "Kotlin",
	".dart":   "Dart",
	".scala":  "Scala",
	".groovy": "Groovy",
	".hs":     "Haskell",
	".elm":    "Elm",
	".clj":    "Clojure",
	".cljs":   "ClojureScript",
	".cljc":   "Clojure",
	".edn":    "EDN",
	".rkt":    "Racket",
	".scm":    "Scheme",
	".ss":     "Scheme",
	".lisp":   "Lisp",
	".cl":     "Common Lisp",
	".el":     "Emacs Lisp",
	".ml":     "OCaml",
	".mli":    "OCaml Interface",
	".f90":    "Fortran",
	".f95":    "Fortran",
	".nim":    "Nim",
	".j":      "J (array programming)",
	".sas":    "SAS",
	".sml":    "Standard ML",
	".fs":     "F#",
	".vb":     "Visual Basic",
	".asm":    "Assembly",
	".s":      "Assembly",
	".nasm":   "Assembly",
	".vhd":    "VHDL",
	".vhdl":   "VHDL",

	// Markup / Web
	".html":       "HTML",
	".htm":        "HTML",
	".css":        "CSS",
	".scss":       "SCSS",
	".sass":       "Sass",
	".less":       "Less",
	".xml":        "XML",
	".xsl":        "XSLT",
	".svg":        "SVG",
	".md":         "Markdown",
	".rst":        "reStructuredText",
	".adoc":       "AsciiDoc",
	".jade":       "Jade/Pug",
	".pug":        "Pug",
	".ejs":        "EJS",
	".handlebars": "Handlebars",
	".hbs":        "Handlebars",

	// Configuration & Data Formats
	".json":  "JSON",
	".yaml":  "YAML",
	".yml":   "YAML",
	".toml":  "TOML",
	".ini":   "INI",
	".conf":  "Configuration",
	".cfg":   "Configuration",
	".plst":  "Property List",
	".plist": "Property List",
	".csv":   "CSV",
	".tsv":   "TSV",
	".env":   "Environment Variables",

	// Scripting and Command Languages
	".sh":   "Shell Script",
	".bash": "Bash",
	".zsh":  "Zsh",
	".fish": "Fish",
	".ps1":  "PowerShell",
	".bat":  "Batch",
	".cmd":  "Batch",
	".sql":  "SQL",

	// Markup Languages for Documents
	".tex": "LaTeX",
	".ltx": "LaTeX",
	".bib": "BibTeX",

	// Miscellaneous / Other
	".log":        "Log File",
	".mdown":      "Markdown",
	".markdown":   "Markdown",
	".lock":       "Lock File",
	".dockerfile": "Dockerfile",
	".makefile":   "Makefile",
	".txt":        "Text",
	".text":       "Text",
	".rtf":        "Rich Text Format",
	".pdf":        "PDF",
}

// Repository defines the interface for obtaining file listings and contents.
type Repository interface {
	GetFiles() ([]string, error)
	GetFileContent(filePath string) (string, error)
	ConcatFiles(files []string) (string, error)
	GetLanguage(filePath string) string
}

type repoCommon struct {
	languages map[string]string
}

func newRepoCommon() *repoCommon {
	rc := &repoCommon{
		languages: make(map[string]string),
	}
	for ext, lang := range defaultLanguageMap {
		rc.languages[ext] = lang
	}
	return rc
}

func (rc *repoCommon) registerLanguage(ext, languageName string) {
	rc.languages[ext] = languageName
}

func (rc *repoCommon) getLanguage(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	if lang, ok := rc.languages[ext]; ok {
		return lang
	}
	return ""
}

// Option is a functional option to modify repository settings.
type Option func(rc Repository)

func WithRegisteredLanguages(langs map[string]string) Option {
	return func(r Repository) {
		var rc *repoCommon
		if lr, ok := r.(*localRepository); ok {
			rc = lr.common
		}

		if gr, ok := r.(*gitRepository); ok {
			rc = gr.common
		}

		if rc != nil {
			for ext, name := range langs {
				rc.registerLanguage(ext, name)
			}
		}
	}
}

// OpenRepository returns a Repository from a given pathOrURL.
//
// If the input starts with "http://" or "https://", it is assumed to be a Git repository URL;
// otherwise, it is assumed to be a local folder.
func OpenRepository(pathOrURL string, opts ...Option) (Repository, error) {
	if strings.HasPrefix(pathOrURL, "http://") || strings.HasPrefix(pathOrURL, "https://") {
		return CloneGitRepository(pathOrURL, opts...)
	}
	return NewLocalRepository(pathOrURL, opts...)
}
