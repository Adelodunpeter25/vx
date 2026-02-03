# Split Panels Plan (Side-by-Side First)

## Phase 0: Baseline and Scoping
- Confirm target behaviors: side-by-side split only, independent cursors, shared buffers optional, preview behavior in splits
- Decide on message/status bar scope: global vs per-pane
- Decide on search/replace scope: per-pane vs global
- Define minimal UX for pane focus, mouse focus, and resize feedback

## Phase 1: Core Data Model Refactor
- Introduce a Pane/View model that separates buffer content from view state
- Move cursor/scroll/selection state into Pane
- Keep Editor as the manager of panes and active pane index
- Add pane-aware helper methods for cursor clamp, scroll adjustment, and selection handling
- Ensure single-pane behavior matches current behavior exactly

## Phase 2: Rendering Pipeline for Multiple Panes
- Create a layout system that defines pane rectangles for side-by-side split
- Refactor rendering to draw each pane within its rectangle
- Update gutter, line wrapping, selection, and search highlighting to respect pane bounds
- Decide status bar layout: single global bar or per-pane status segments

## Phase 3: Input Routing and Focus
- Route keyboard input to active pane only
- Implement mouse hit-testing to select pane focus on click
- Ensure selection and cursor movement operate within the focused pane
- Add visual focus indicator (e.g., active pane border or status highlight)

## Phase 4: Mouse Support for Split and Resize
- Implement mouse-driven focus: click inside pane focuses it
- Add a draggable splitter bar between panes
- Implement resize logic that updates pane widths and reflows content
- Add resize constraints (min pane width, gutter width, status bar safety)
- Handle mouse capture state during drag to prevent selection while resizing

## Phase 5: Buffer and Pane Interactions
- Decide if panes can show the same buffer or different buffers
- If shared buffers are allowed, ensure per-pane cursor/selection remains isolated
- Update buffer switch commands to operate on active pane only
- Update buffer close behavior to handle panes that reference the same buffer

## Phase 6: Search/Replace and Preview Behavior
- If per-pane search: keep search state in Pane and scope highlights to that pane
- If global search: define how results are rendered in each pane
- Decide preview behavior in split mode (per-pane preview or full-screen override)

## Phase 7: Edge Cases and UX Polish
- Define behavior for very narrow panes (line numbers, wrapping, status clipping)
- Ensure mouse selection respects pane bounds and does not cross panes
- Handle resize while selection active (preserve selection or clear)
- Ensure pane focus persists across buffer switches

## Phase 8: Testing and Verification
- Manual tests: focus switching, resize, scrolling, selection, search, replace
- Stress tests: large files, wrapped lines, rapid resize, multi-pane buffer changes
- Regression checks: single-pane mode works exactly as before
