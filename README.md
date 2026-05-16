# VX Editor

A modern vi text editor written in Go with syntax highlighting, search, undo/redo, and markdown preview.

## Features

- **Modal Editing** - Classic vi-style normal, insert, and command modes
- **Syntax Highlighting** - Support for 200+ languages via Chroma
- **Split Panes** - Side-by-side panes for editing multiple files
- **File Browser** - Toggleable left sidebar for navigating folders/files
- **File Palette** - `Ctrl+P` fuzzy picker for quickly opening files
- **Mouse Selection** - Click and drag to select text, copy with `c`, cut with `x`
- **Real-time Search** - Incremental search with live highlighting as you type
- **Find & Replace** - Interactive replace with y/n confirmation for each match
- **Undo/Redo** - Full undo history with `Ctrl+Z` / `Ctrl+Y`
- **Clipboard Support** - Copy with `Ctrl+C`, paste with `Ctrl+V`
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

#### Navigation
- `h/j/k/l` - Move cursor left/down/up/right
- `w/b` - Move forward/backward by word
- `gg` / `G` - Jump to start / end of file
- `Ctrl+U` / `Ctrl+D` - Page up / page down

#### Editing
- `i` - Enter insert mode
- `x` / `Ctrl+X` - Cut (selected text or character)
- `p` - Paste (or toggle preview for .md files)
- `u` / `Ctrl+Z` - Undo
- `r` / `Ctrl+Y` - Redo
- `dd` - Delete current line

#### Clipboard
- `Ctrl+C` - Copy (selected text or current line)
- `Ctrl+V` - Paste
- `c` - Copy (selected text or current line)
- `x` / `Ctrl+X` - Cut

#### Search & Replace
- `/` or `Ctrl+F` - Search (real-time incremental)
- `n` / `N` - Next / previous search result
- `Shift+H` - Find and replace

#### File & Pane Management
- `Ctrl+S` - Save file
- `Ctrl+O` - Open file
- `Ctrl+W` - Close current pane
- `Ctrl+N` - Next pane
- `Ctrl+P` - Open file palette
- `Ctrl+B` - Toggle file browser sidebar
- `Esc` - Clear selection
- `:` - Enter command mode

#### Selection
- `Ctrl+A` - Select all

#### Misc
- `q` / `Ctrl+Q` - Quit
- `Ctrl+C` - Copy

### Insert Mode
- `Esc` - Return to normal mode
- `Ctrl+S` - Save
- `Ctrl+Z` - Undo
- `Ctrl+Y` - Redo
- Arrow keys - Navigate

### Mouse Selection
- **Click and drag** - Select text (auto-scrolls at edges)
- `Ctrl+C` - Copy selected text to clipboard
- `Ctrl+X` - Cut selected text
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
- `:f` - Toggle file browser sidebar
- `:cd [path]` - Change the current directory
- `:set:show-hidden` - Show hidden files in file browser
- `:set:hide-hidden` - Hide hidden files in file browser

### Markdown Preview
- `p` - Toggle preview (in .md files, normal mode)
- `j/k` or arrows - Scroll
- `gg` / `G` - Start / end
- `Ctrl+U` / `Ctrl+D` - Page up / page down

## Philosophy

VX is "vi, but modern" - keeping the classic vi modal editing experience while adding modern conveniences like syntax highlighting and better UX. It's not trying to be Vim or Neovim, just a fast, simple text editor that respects your muscle memory.
