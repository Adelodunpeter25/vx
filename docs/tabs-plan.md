# Tabs Plan — LazyVim-style for vx

> Click a file in the file-browser / palette → new tab if not already open, otherwise switch. Top bar shows open files like `bufferline.nvim` in LazyVim.

## 1) Goal and Non-Goals

- **Goal:** Tab bar as in LazyVim `bufferline.nvim` — top row shows open files as tabs, clicking a file in browser/palette adds new tab if not already open, otherwise switches. Not vim window-splits (`:sp`/`:vsp` remain as `panes`).
- **Non-goals (MVP):** Not replacing splits; not vim-style tabpages that each hold a layout. Keep `panes []*Pane` as is.

## 2) Current Architecture (read-only snapshot)

- `internal/editor/editor.go:19` `Editor{ panes []*Pane, activePane int, fileBrowser }` — splits already; no concept of buffers/tabs
- `internal/editor/pane.go:13` `Pane{ buffer, syntax, cursorX/Y, offsetX/Y, visualOffsetY, selection, search }` — per-buffer view state; reusable for tab
- `internal/editor/file_browser.go:14` `openFileInActivePane` / `previewFileInActivePane` reuses active pane — this is where "click should add tab" will branch
- `internal/editor/palette.go:30` `OnSelect → openFileInActivePane` — second entry point
- `internal/editor/render.go:9` `render()` draws browser left + splits side-by-side + status line; tab bar would be `y=0, h=1` above content
- `internal/terminal/terminal.go` mouse `Button1` already routed per-pane in `handleMouseEventForPane:266` — same needed for tab bar hit test

## 3) Data Model (MVP, keeps splits)

```go
// internal/editor/tab.go (new file)
type Tab struct {
    Pane    *Pane  // owns buffer/syntax/cursor/scroll/selection
    Title   string // display: filepath.Base + modified dot
    Dirty   bool
    Pinned  bool
    Preview bool // single reusable preview tab (VS Code/LazyVim)
}
type Editor struct {
    tabs      []*Tab
    activeTab int
    tabScroll int
    panes     []*Pane // keep as is for splits; activeTab's Pane is shown in panes[activePane]
}
```

**MVP choice:** tabs are global buffers shown in bar; `activeTab`'s `Pane` is displayed in `panes[activePane]`. When splits exist (>1 pane), tabs switch the *active pane's* buffer, not all panes — matches LazyVim bufferline (buffer-centric). Alternative `tabs hold split layout` (vim tabpages) deferred to Phase 4.

**Preview tab:** single preview tab reused for single-clicks in browser (`HandleMouse:493` `PreviewPath`), double-click / `Enter` pins it → converts preview to real tab. Same as VS Code / LazyVim.

## 4) UI — render.go / new render_tab.go

- Bar at `y=0, h=1, x= browserWidth+1 ... width-1`. Each tab: `icon + basename + modified • + close x`. Active tab: `Bg #45475a Bold`, inactive: `Bg #1e1e2e`. Overflow → horizontal `tabScroll`.
- Reuse `drawClippedText` from `file-browser/browser.go:544` and `utils.FileIconFor`.
- Height budget: `contentHeight = e.height -1(status) -1(tabs) - cdPromptRows`. When `len(tabs)==0` fallback to no bar.
- Mouse: top row hit-test `if ev.MouseY==0` → `tabAtX(ev.MouseX)` → `activeTab=idx`, update `panes[activePane]` buffer to tab's buffer. `WheelUp/Down` over bar scrolls bar. Right-click / middle-click closes.

## 5) Interactions — file open paths to add/switch tab

- `file_browser.go:openFileInActivePane(path)`:
  ```
  if idx := e.findTabByPath(path); idx>=0 { e.setActiveTab(idx); return }
  newPane := NewPane(newBuf, path); tab := &Tab{Pane:newPane}
  e.tabs = append(e.tabs, tab); e.activeTab = len-1
  e.panes[e.activePane] = newPane  // reuse active split slot, keep other splits
  ```
- `previewFileInActivePane`: if preview tab exists, reuse its `Pane.buffer`; else create preview tab at end. Next `open` promotes preview.
- `command_mode.go:55` `:e file` same as open (adds tab). `:b file` retains old meaning: open in new pane (`addPaneWithBuffer` also adds tab).
- `palette.go:30` same as browser.
- Close: `:q` / `Ctrl+W` / mouse `x` → `closeTab(idx)` removes from `tabs` and if that tab's `Pane` is displayed in any `panes[i]`, replace with neighbor tab's Pane. If last tab, keep empty `[No Name]` tab.

## 6) Keys / Commands (additive, not breaking)

- `gt`/`gT` or `Ctrl+Tab` / `Ctrl+Shift+Tab` → `nextTab`/`prevTab` (like `nextPane:169` but for `tabs`)
- `Mouse Button1` on tab → switch, `Button2/Middle` or `x` click → close, `Shift+click` pin
- `:tabnew`, `:tabclose` aliases to pane ops for compat

## 7) Phases

- **P1 — Model:** new `internal/editor/tab.go` (`Tab`, `findTabByPath`, `setActiveTab`, `closeTab`, `renderTabs`). Extend `Editor` with `tabs, activeTab, tabScroll`. Init `New`/`NewWithFile` creates initial tab from initial pane.
- **P2 — Render:** `render.go:9` reserve bar, call `renderTabs` before splits; `handleMouseEventForPane` early return for `y==0` tab bar hits; `handleResize` invalidates tabs.
- **P3 — Switch:** patch `file_browser.go:14` + `previewFileInActivePane:38` + `command_mode.go:55` + `palette.go:30` to use tab helpers; add `resetViewport` already fixed in `pane.go:60`.
- **P4 — Advanced (deferred):** pin, dirty indicator, overflow scrolling, session persist `.vx/tabs.json`, optional "tabs hold layout" mode.

## 8) Files Touched

- `internal/editor/tab.go` (new), `editor.go`, `pane.go`, `render.go`, `render_tab.go` (new), `file_browser.go`, `command_mode.go`, `palette.go`, `terminal/event.go` (if new tab events), `docs/tabs-plan.md` (this file)

## 9) Risks

- Tab bar steals one row from `contentHeight` — must clamp `visualOffsetY` already fixed in `cursor.go:35`.
- Preview reuse may surprise users expecting always new tab — make preview opt configurable (`showPreview`).
- Too many tabs overflow — implement scroll, not wrap.

## 10) Verification

- `go build -o build/vx cmd/vx/*.go` + `go test ./tests -count=1`
- Manual: open dir, single-click file → preview tab updates, double-click → pinned tab, palette `Ctrl+P` → adds, `gt` cycles, close keeps buffer.
