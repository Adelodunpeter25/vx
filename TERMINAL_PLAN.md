# Terminal Design Plan

## Core choice
- PTY + `charmbracelet/x/vt` emulator
- Render onto tcell with a diffing renderer + frame throttle

## 1) Architecture
- Terminal pane owns:
  - PTY process
  - Emulator (`vt.SafeEmulator`)
  - Scrollback buffer
  - Render cache (last frame)
  - Input encoder
- Editor only orchestrates:
  - Layout + focus
  - Resize events
  - Toggle `:t`

## 2) Input handling
- Keyboard → VT sequences (arrows, function keys, ctrl combos)
- Mouse → optional (scroll, click, drag) behind a flag
- On input:
  - Write to PTY
  - Request redraw (non-blocking)

## 3) Output handling
- PTY reader goroutine:
  - Write into emulator
  - Mark terminal “dirty”
  - Coalesce redraw requests

## 4) Rendering performance
- Frame throttle: max 60 FPS, merge bursts
- Dirty diff: compare last frame cells vs current and only update changed cells
- Short-circuit: if no changes, no redraw

## 5) Scrollback
- Maintain ring buffer of screen snapshots or rows
- Wheel/pgup scrolls history without touching emulator state
- Cursor/typing snaps back to live view

## 6) Resize
- Resize PTY + emulator only on actual size change
- Preserve scrollback

## 7) Features checklist
- Alternate screen support (TUI apps)
- Truecolor + 256 color
- Unicode handling
- Cursor shape + blink (optional)

## 8) Testing & stability
- Manual test matrix:
  - `htop`, `top`, `vim`, `less`, `gh`/`claude`, `vx`
  - Resize while running
  - Large output spam (`yes`, `find /`)
- Performance test: sustained output + input without lag
