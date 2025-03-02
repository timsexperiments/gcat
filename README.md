# gcat

**gcat** is a command‑line tool written in Go for concatenating files from either a Git repository (remote or local) or a local folder. It provides an interactive file selection UI, supports file–type detection based on file extensions, allows ignoring hidden or unwanted files (using a `.gitignore` or user-defined ignore patterns), and can optionally copy the concatenated output to your clipboard.

This tool is perfect for creating a single, LLM–friendly string containing multiple files’ contents (with each file’s path and detected language as a header) for further analysis or feeding to other tools.

---

## Features

- **Unified Source Handling:**  
  Works with Git repositories and local directories. Use a URL (starting with `http://` or `https://`) for Git or pass a local folder path.

- **Interactive File Selection:**  
  Uses [Survey](https://github.com/AlecAivazis/survey/v2) for an interactive multi-select prompt. The prompt shows a sorted list of files (limited to 10 visible options) and allows you to toggle your selection with the spacebar.

- **File Concatenation:**  
  Concatenates selected files into a single output string. Each file is preceded by its file path and a naive language detection header based on its extension.

- **Clipboard Support:**  
  Optionally copy the result directly to your system clipboard using [golang.design/x/clipboard](https://pkg.go.dev/golang.design/x/clipboard).

- **Ignore Hidden/Unwanted Files:**  
  The local repository implementation filters out hidden files (names starting with a dot) and files/directories defined in default ignore patterns (e.g. `.git`).

---

## Installation

### Option 1: Using Go Install (recommended)

```bash
go install github.com/timsexperiments/gcat/cmd/gcat@latest
```

### Option 2: Pre-built binaries

Download the pre-built binary for your platform from the [releases page](https://github.com/timsexperiments/gcat/releases).

### Option 3: Build from source

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/timsexperiments/gcat.git
   cd gcat
   ```

2. **Install Dependencies:**

   Ensure you have Go installed (version 1.23.5 or later). Then run:

   ```bash
   go build -o gcat ./cmd/gcat
   ```

   This command compiles the CLI tool and produces an executable named gcat.

## Usage

The tool accepts the source as a required argument (either a Git repository URL or a local folder path). Use the optional --copy (or -c) flag to copy the concatenated result to the clipboard instead of printing it to the terminal.

### Examples

- **Concatenate files from a remote Git repository:**

  ```bash
  ./gcat https://github.com/googleapis/api-linter
  ```

- **Concatenate files from a local directory:**

  ```bash
  ./gcat /path/to/local/folder
  ```

- **Copy the concatenated output to the clipboard (do not print to console):**

  ```bash
  ./gcat --copy https://github.com/googleapis/api-linter
  ```

## How It Works

1. **Source Detection:**

   The source argument is checked—if it starts with http:// or https:// it is treated as a Git repository; otherwise, it is assumed to be a local folder.

2. **Repository Handling:**

   - **Git Repositories:**

     The tool clones the repository shallowly (in-memory) using go-git.

   - **Local Repositories:**

     It performs a file-walk starting from the given folder, ignoring hidden files (names starting with a dot) and files/directories defined by default ignore patterns.

3. **File Selection:**

   An interactive selection UI (implemented using Survey’s MultiSelect prompt) displays the list of files in alphabetical order (limited to 10 visible options) and allows you to toggle your file selections with the spacebar. Confirm your selection with Enter.

4. **Concatenation:**

   Selected files are read and concatenated into a single string. Each file's section includes its file path, a naive language detection header (based on extension), followed by the file contents.

5. **Output:**

   - When the `--copy` flag is not set, the concatenated result is printed to the console.

   - When --copy is specified, the result is copied to the clipboard (using [golang.design/x/clipboard](https://pkg.go.dev/golang.design/x/clipboard)) instead of printing.

## Project Structure

```
gcat/
├── cmd/
│   └── gcat/
│       └── main.go
├── internal/
│   ├── cli/
│   │   └── selector.go
│   └── clipboard/
│       └── clipboard.go
└── pkg/
    └── gcat/
        ├── gcat.go
        ├── local.go
        └── remote.go
```

## License

Distributed under the MIT License. See [LICENSE](./LICENSE.md) for more information.
