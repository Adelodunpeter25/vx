# Git Status Indicators Implementation Guide

## Goal

Add git status indicators to the file browser so files and folders visibly show repository state such as modified, added, deleted, untracked, ignored, conflicted, and clean. The indicators should be lightweight, readable in a terminal, and integrated into the current `vx` architecture without turning the file browser into a separate subsystem.

## Scope

This feature should primarily affect:

- File browser rows for files and folders
- Optional status data in the editor/file browser state
- Refresh behavior when the working tree changes
- Display conventions for terminal rendering

This should not require changing core editing behavior unless the file browser needs new callbacks or refresh hooks.

## Existing Structure To Reuse

Relevant code paths in this repository:

- [`internal/file-browser/browser.go`](../internal/file-browser/browser.go)
- [`internal/editor/file_browser.go`](../internal/editor/file_browser.go)
- [`internal/editor/render.go`](../internal/editor/render.go)
- [`internal/editor/render_pane.go`](../internal/editor/render_pane.go)
- [`internal/utils/file-icons.go`](../internal/utils/file-icons.go)

Current file browser behavior already includes:

- Directory tree building
- Expand/collapse state
- File icons
- Hidden file filtering
- Selection and preview support

That means the git feature should be added as metadata layered on top of the existing tree, not as a replacement for it.

## Recommended Display Model

Use short text markers that work well in a narrow terminal UI.

Suggested status symbols:

- `M` modified
- `A` added
- `D` deleted
- `?` untracked
- `!` ignored
- `U` conflicted/unmerged
- ` ` clean

Suggested presentation:

- Prefix or suffix the icon with a small badge
- Keep directory badges aggregated and conservative
- Prefer one-byte or two-byte markers over long labels
- Use color only as a secondary cue, not the only cue

Example row formats:

- `  src/`
- ` M main.go`
- ` ? new.txt`
- ` M internal/`

## Phase Plan

### Phase 1: Define the repository metadata model

Purpose:

- Decide how git state is represented inside the app
- Keep the data model small enough that rendering remains simple

Tasks:

- Add a status enum or small set of constants in a new package or in the file browser package
- Define a per-node status field for files
- Define an optional aggregate status field for directories
- Decide whether a node can hold multiple flags or only one primary status
- Decide how to represent "dirty directory" state when any child is modified

Implementation notes:

- A node should probably store one primary status plus optional child summary counts
- Directory status should usually be derived from descendants, not guessed from directory mtime
- Keep the model independent from any specific git command output format

Exit criteria:

- A node can represent `clean`, `modified`, `untracked`, `ignored`, `deleted`, `added`, or `conflicted`
- Directory nodes can derive a summary state from child nodes

### Phase 2: Build a git status collector

Purpose:

- Read repository state from git without blocking normal editor use

Tasks:

- Detect whether the current root is inside a git repository
- Collect status for paths under the repo root
- Map git status output to internal node statuses
- Handle subdirectories and nested paths consistently
- Ignore paths that are not part of the file browser tree

Recommended implementation approach:

- Use `git status --porcelain=v1 -z` or a similar stable porcelain format
- Parse results into a path-to-status map
- Normalize paths relative to the repository root
- Cache results until refresh is needed

Important constraints:

- Do not shell out on every render
- Avoid expensive scans on every cursor move or mouse event
- Fallback gracefully when git is unavailable or the folder is not a repo

Exit criteria:

- The app can produce a repository-relative map of statuses
- Non-git folders still render normally

### Phase 3: Wire git metadata into file browser state

Purpose:

- Make file nodes aware of git status without coupling rendering to parsing logic

Tasks:

- Extend file browser state with git status data
- Add helper methods to resolve status for a node path
- Add aggregate status support for parent directories
- Ensure root refresh logic updates both tree data and git metadata

Suggested approach:

- Load git status once when the browser root changes
- Recompute on explicit refresh and file system change events
- Use path lookup at render time, but keep the lookup O(1) with a map

Exit criteria:

- A tree node can ask for its status in constant time
- Directory rows can derive a summarized status from children

### Phase 4: Update rendering

Purpose:

- Show git state clearly in the existing tree layout

Tasks:

- Add a badge column or inline marker in file browser row rendering
- Add color styling for status types
- Keep the current icon and indentation layout intact
- Make sure long file names still clip cleanly
- Preserve selection highlight behavior

Rendering rules:

- Status badge should not break alignment
- Directory status should be visible even when collapsed
- For directories, prefer a summary badge if descendants have mixed states
- For files, show the exact file state

Exit criteria:

- File browser rows visibly show status
- The tree remains readable at narrow widths

### Phase 5: Handle refresh and invalidation

Purpose:

- Keep status indicators accurate as the working tree changes

Tasks:

- Refresh git status when the watcher detects file changes
- Refresh when the file browser root changes
- Refresh when `:cd` changes directories
- Consider a manual refresh command later if needed

Notes:

- The repo already has a file watcher in [`internal/editor/fswatch.go`](../internal/editor/fswatch.go)
- That watcher can trigger tree refreshes, but git status should still be cached separately
- If refresh becomes expensive, debounce it

Exit criteria:

- Modified files update after save or external edit
- Newly created/deleted files update after watcher refresh

### Phase 6: Add tests or verification checks

Purpose:

- Reduce regression risk in parsing and directory aggregation

Tasks:

- Add unit tests for git status parsing
- Add unit tests for directory aggregation rules
- Add rendering tests only if there is already a test harness for terminal output
- Verify behavior in a real git repository and a non-git directory

Recommended test cases:

- Clean repo returns clean statuses
- Untracked file returns `?`
- Modified file returns `M`
- Deleted file returns `D`
- Nested directory aggregates child state
- Non-git directory returns empty status data

Exit criteria:

- Parsing and aggregation are covered
- Manual verification is performed in a real terminal session

## Implementation Tasks

### Task 1: Introduce a status type

- Create a small status type or enum
- Add helpers for label, color, and badge rendering
- Keep it simple enough to reuse for files and directories

### Task 2: Add a collector package or module

- Parse git porcelain output
- Build a map from relative path to status
- Expose lookup helpers to the file browser

### Task 3: Extend file browser nodes

- Add a status field to the node model
- Add directory summary support
- Preserve current icon and expansion behavior

### Task 4: Refresh status at the right times

- On browser root change
- On editor watcher events
- On explicit root changes from `:cd`
- On file open/reload if needed

### Task 5: Render badges

- Add a small fixed-width badge area
- Apply status colors and preserve selection styles
- Keep clipping behavior stable

### Task 6: Verify behavior

- Run `go test ./...`
- Run `make build`
- Open a repo with modified and untracked files
- Confirm hidden-file and directory expansion behavior still works

## Suggested File Layout

If the implementation is split into new files, a practical layout would be:

- `internal/gitstatus/collector.go`
- `internal/gitstatus/status.go`
- `internal/file-browser/status.go`
- `internal/file-browser/browser.go`

This is only a suggestion. The exact structure should follow the existing package boundaries and avoid creating circular dependencies.

## Risks And Tradeoffs

- Running git commands too often can make the browser feel sluggish
- Directory aggregation can become confusing if mixed states are overrepresented
- Status badges can crowd narrow terminal layouts if they are too verbose
- Shelling out to git is simpler than using a library, but it needs careful caching
- The file browser currently handles tree loading and rendering together, so status integration should stay lightweight

## Definition Of Done

The feature is done when:

- Files show accurate git status indicators
- Folders show useful aggregated status indicators
- The browser remains responsive
- Non-git folders behave normally
- `go test ./...` passes
- `make build` passes

## Next Step Recommendation

Implement the collector and data model first. Rendering should come second. If the status map is correct, the UI work will be straightforward.
