# VX Editor

A modern vi text editor written in Go with syntax highlighting, search, undo/redo, and markdown preview.

## Features

- **Modal Editing** - Classic vi-style normal, insert, and command modes
- **Syntax Highlighting** - Support for 200+ languages via Chroma
- **Split Panes** - Side-by-side panes for editing multiple files
- **Mouse Selection** - Click and drag to select text, copy with `c`, cut with `x`
- **Real-time Search** - Incremental search with live highlighting as you type
- **Find & Replace** - Interactive replace with y/n confirmation for each match
- **Undo/Redo** - Full undo history with `u` and `r`
- **Clipboard Support** - Copy with `c`, paste with `p`
- **Markdown Preview** - Full-screen preview for `.md` files
- **Fast Startup** - Native Go binary, instant launch

## Installation

### Quick Install (Recommended - MacOS/Linux)

```bash
curl -sSL https://raw.githubusercontent.com/Adelodunpeter25/vx/main/install.sh | bash
```
### Manual Installation

Download the appropriate binary from [releases](https://github.com/Adelodunpeter25/vx/releases):

- **Linux (x64)**: `vx-linux-amd64`
- **Linux (ARM64)**: `vx-linux-arm64`
- **macOS (Intel)**: `vx-darwin-amd64`
- **macOS (Apple Silicon)**: `vx-darwin-arm64`
- **Windows**: `vx-windows-amd64.exe`

Make it executable and move to PATH:

```bash
chmod +x vx-*
sudo mv vx-* /usr/local/bin/vx
```
### Build from Source

```bash
git clone https://github.com/Adelodunpeter25/vx.git
cd vx
go build -o vx cmd/vx/*.go
```
## Usage

```bash
vx [filename]         # Open file for editing
vx --help             # Show help
vx --version          # Show version
```
## Keybindings

### Normal Mode
- `h/j/k/l` - Move cursor left/down/up/right
- `w/b` - Move forward/backward by word
- `i` - Enter insert mode
- `:` - Enter command mode
- `/` or `Ctrl+F` - Search (real-time incremental search)
- `Shift+H` - Find and replace
- `n/N` - Next/previous search result
- `c` - Copy current line (or selected text if selection active)
- `x` - Cut selected text (or delete character if no selection)
- `p` - Paste (or toggle preview for .md files)
- `dd` - Delete current line
- `u` - Undo
- `r` - Redo
- `gg` - Jump to start of file
- `G` - Jump to end of file
- `Ctrl+S` - Save file
- `Ctrl+N` - Next pane
- `Ctrl+P` - Previous pane
- `Esc` - Clear selection
- `q` - Quit
- `Ctrl+C` - Force quit

### Mouse Selection
- **Click and drag** - Select text (auto-scrolls at edges)
- `c` - Copy selected text to clipboard
- `x` - Cut selected text (copy and delete)
- `Esc` or any movement key - Clear selection

### Search Mode
- Type to search - Results highlight in real-time as you type
- `Enter` - Exit search mode (keep highlights)
- `Esc` - Cancel search

### Replace Mode
- `Shift+H` - Start find and replace
- Type search term, press `Enter`
- Type replacement term, press `Enter`
- For each match:
  - `y` - Replace this match
  - `n` - Skip this match
  - `q` - Quit replace mode

### Command Mode
- `:w` - Save file
- `:w filename` - Save as filename
- `:q` - Quit
- `:q!` - Force quit without saving
- `:wq` - Save and quit
- `:e filename` - Edit new file (replace current pane)
- `:b filename` - Open file in new pane
- `:db` - Close current pane (prompts to save if modified)

### Markdown Preview
- `p` - Toggle preview (in .md files(normal mode))
- `j/k` or arrows - Scroll preview

## Philosophy

VX is "vi, but modern" - keeping the classic vi modal editing experience while adding modern conveniences like syntax highlighting and better UX. It's not trying to be Vim or Neovim, just a fast, simple text editor that respectsyour muscle memory..
