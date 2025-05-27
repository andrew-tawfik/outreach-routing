## Step-by-Step Process for Vehicle UI Refactor
### 1. Establish the UI Hierarchy

- Create a custom vehicleGrid widget (similar to boardContainer in chess)
- This grid will contain all vehicle cards and handle global drag operations
- Each vehicle card will have a dynamic grid of tiles (g + 2 tiles)
- Guests will be draggable widgets that can move between any tile

### 2. Create a Global Drag Layer

- Add an overlay container at the top level (like the chess over image)
- This overlay will display the guest being dragged
- The overlay must be outside all vehicle cards to allow cross-vehicle dragging
- Structure: container.NewMax(vehicleGrid, container.NewWithoutLayout(dragOverlay))

### 3. Implement Custom Tile Widget

Create a guestTile widget that can:

- Display a guest (if occupied)
- Accept dropped guests
- Highlight when a guest is hovering over it
- Each tile knows its position (vehicleIndex, tileIndex)

### 4. Track Global State

Maintain a central state manager that tracks:

- Which guest is being dragged
- Original position (vehicle, tile)
- Current hover position
- Valid drop targets


This prevents the scoping issues you're experiencing

### 5. Implement Drag Mechanics
- DragStart:
  - Hide guest from original tile
  - Show guest in overlay at mouse position
  - Record original position in global state

- Dragging:
  - Update overlay position to follow mouse
  - Calculate which tile is under mouse
  - Highlight valid drop targets

- DragEnd:
  - Determine final tile position
  - Validate the move
  - Update data model
  - Refresh affected vehicle grids
  - Hide overlay
### 6. Handle Dynamic Grid Resizing

When a guest moves between vehicles:

- Source vehicle: Reduce grid to (remaining guests + 2)
- Target vehicle: Increase grid to (new guest count + 2)
- Redistribute remaining guests to fill gaps
- Trigger UI refresh for both vehicles



### 7. Implement Position Calculation

Create helper functions to:

- Convert screen coordinates to (vehicleIndex, tileIndex)
- Calculate tile boundaries globally (not relative to vehicle card)
- Account for scrolling if vehicle grid is scrollable



### 8. Add Visual Feedback

- Highlight valid drop zones during drag
- Show invalid drop indicators (like chess red highlight)
- Animate guest movement (like chess piece animation)
- Add hover effects on tiles

### 9. Handle Edge Cases

- Prevent dropping multiple guests on same tile
- Handle dragging outside valid areas
- Manage rapid drag operations
- Handle window resizing during drag

### 10. Implement Submit/Reset

- Track initial state when view loads
- Submit: Apply changes to underlying RouteManager
- Reset: Restore to initial state and refresh UI
- Consider adding undo/redo functionality

### Key Differences from Current Approach:

- Global drag handling instead of per-vehicle handling
- Overlay layer for dragging across boundaries
- Central state management instead of distributed state
- Dynamic grid sizing based on guest count
- Global coordinate system for drag calculations