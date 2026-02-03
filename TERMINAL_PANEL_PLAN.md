# Bottom Terminal Panel Plan

## Phase 0: Scope and Behavior
- Define terminal panel height and default state (closed/open).
- Decide input focus toggle behavior between editor panes and terminal panel.
- Decide how terminal output scrollback is managed.

## Phase 1: Data Model and State
- Add terminal panel state to editor (open/closed, height, scroll offset).
- Define a terminal buffer model to store output and input line.
- Add focus tracking (active pane vs terminal panel).

## Phase 2: Rendering
- Reserve bottom rows when terminal panel is open.
- Render terminal panel with a distinct style and divider line.
- Render terminal content with scrollback and input prompt.

## Phase 3: Input Handling
- Add `:t` command to toggle terminal panel.
- When terminal is focused, route keystrokes to the input buffer.
- Add keybindings to switch focus back to editor panes.

## Phase 4: Command Execution
- Execute shell commands from the terminal input.
- Append stdout/stderr to terminal buffer with timestamps or markers.
- Handle long-running commands and cancellation.

## Phase 5: Mouse Support
- Mouse click to focus terminal panel.
- Scroll wheel to move through terminal scrollback.
- Optional: resize terminal panel height by dragging divider.

## Phase 6: Edge Cases and Testing
- Ensure terminal panel does not break pane layout or status bar.
- Handle very small terminal window sizes.
- Verify output rendering with long lines and ANSI stripping (or ignore ANSI).
