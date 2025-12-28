# VX Editor - Pre-Release Improvements & Polish

## Critical Issues to Fix

### Stability
- Handle very large files gracefully (current full-buffer highlighting will crash on huge files)
- Add error recovery for corrupted files
- Prevent crashes on invalid UTF-8 sequences
- Handle terminal resize during editing without losing cursor position

### Performance
- Cache syntax highlighting only when buffer is modified (currently re-highlights too often)
- Optimize rendering to only redraw changed lines
- Add lazy loading for large files (don't load entire file into memory)

### Data Safety
- Warn before overwriting existing files with :w filename
- Add backup file creation (.swp or ~ files)
- Detect if file changed on disk since opening
- Handle disk full errors gracefully when saving

## User Experience Improvements

### Better Feedback
- Show file size and line count in status bar
- Display "Saving..." message during save operations
- Show error messages for longer (currently disappear on next render)
- Add visual feedback when at start/end of file

### Missing Vi Commands
- Add 'a' (append after cursor)
- Add 'o' and 'O' (open line below/above)
- Add 'x' (delete character)
- Add 'dd' (delete line)
- Add 'u' (undo) - this is critical
- Add '0' and '$' (start/end of line)
- Add 'gg' and 'G' (start/end of file)
- Add 'w' and 'b' (word navigation)

### Command Mode Enhancements
- Add :e filename (edit new file)
- Add :set number (toggle line numbers)
- Add :set syntax on/off (toggle highlighting)
- Add command history (up/down arrows in command mode)
- Tab completion for filenames

### Visual Polish
- Add line numbers (optional, togglable)
- Show current line highlight (subtle background)
- Better color scheme (current colors may not work on all terminals)
- Add visual mode for selecting text
- Show "-- INSERT --" style mode indicator

## Code Quality

### Testing
- Add unit tests for buffer operations
- Add tests for bracket matching
- Test with various file encodings
- Test on different terminal emulators

### Documentation
- Add README with installation instructions
- Document keybindings
- Add usage examples
- Create man page

### Build & Distribution
- Add Makefile for easy building
- Create install script
- Add version flag (-v, --version)
- Add help flag (-h, --help)
- Consider adding to package managers (brew, apt, etc)

## Nice-to-Have Features

### Search
- Add '/' for search
- Add 'n' and 'N' for next/previous match
- Highlight search results

### Configuration
- Add config file support (~/.vxrc)
- Allow custom color schemes
- Configurable tab width
- Configurable key bindings

### Advanced Editing
- Add clipboard support (yank/paste)
- Add multiple file support (buffers)
- Add split windows
- Add macros (record/replay)

### Terminal Integration
- Better handling of mouse events (optional)
- Support for 24-bit true color terminals
- Handle different terminal sizes gracefully

## Priority Order

1. **Must Fix Before Release**
   - Undo/redo functionality
   - Large file handling
   - Data safety (backups, warnings)
   - Basic vi commands (a, o, x, dd, 0, $, gg, G)

2. **Should Have**
   - Line numbers
   - Search functionality
   - Better error messages
   - README and documentation

3. **Nice to Have**
   - Configuration file
   - Multiple buffers
   - Visual mode
   - Clipboard integration

4. **Future Enhancements**
   - Split windows
   - Macros
   - Plugin system
   - LSP integration
