# File Browser Plan

Phase 0: Scope and state
- Define file browser state (open/closed, width, focus, root path, expanded nodes, selected item)
- Add minimal tree model to represent directories/files and lazy-load children
- Add render slot in layout (left sidebar) and reserve width when open

Phase 1: Rendering and navigation
- Render sidebar list with folders/files, indentation, and icons/markers
- Keyboard navigation: up/down, left/right to collapse/expand, Enter to open
- Mouse navigation: click to focus, click to open file, click folder to toggle

Phase 2: Open and preview behavior
- Open file into active pane on Enter/click
- Keep selection in browser when file opens (no focus loss unless requested)
- Add optional preview on hover or single click (configurable later)

Phase 3: Commands and focus
- Implement :f to toggle sidebar open/closed
- Implement focus switching between sidebar and panes
- Add status hints for navigation keys

Phase 4: Performance and polish
- Cache directory reads and refresh on demand
- Handle large folders with lazy expansion
- Add simple search/filter within the sidebar

Phase 5: Future extensions
- Create/rename/delete actions
- Git status decorations
- Favorites/recents
