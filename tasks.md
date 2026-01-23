# VX Editor Tasks

## Completed ✓
- ~~Linux clipboard not working~~ - Switched to atotto/clipboard, install.sh now prompts for xclip
- ~~Directory crash~~ - Added directory detection with clear error message
- ~~Search highlight persistence~~ - Clear highlights when replace completes
- ~~Auto-expiring messages~~ - Implemented transient/persistent/error message system
- ~~Mouse selection issues~~ - Fixed drag-only detection
- ~~Cursor positioning with wrapping~~ - Fixed insert mode and scroll calculations

## TODO
- [ ] Test clipboard on Linux (needs xclip/xsel installed)
- [ ] Verify Windows build works (untested)

## Future Features (Post v0.1.0)
- [ ] Split windows (horizontal/vertical)
- [ ] Configuration file (keybindings, colors, tab width)
- [ ] Line number toggle
- [ ] Search wrap-around option
- [ ] Better paste behavior (paste on new line option)
- [ ] Delete to register (dd should allow pasting deleted lines)

## Nice to Have
- [ ] LSP integration (basic go-to-definition, diagnostics)
- [ ] Terminal integration
- [ ] Tab pages