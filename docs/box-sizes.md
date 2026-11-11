# Box Sizes and Grid Layout

## Overview

The dashboard uses a CSS Grid layout with tetris-style packing to efficiently display boxes of varying sizes
without vertical gaps. All box dimensions are calculated to tile perfectly on a 13px base grid with 6px gaps.

## Why Grid Instead of Flexbox?

**Original Problem (Flexbox):**

- Mixed box heights created vertical gaps
- Example: A `dsmall` (138px tall) next to a `small` (66px tall) left empty space
- Boxes couldn't "stack" to fill gaps efficiently
- Horizontal alignment worked, but vertical packing was poor

**Solution (CSS Grid with Dense Packing):**

- `grid-auto-flow: dense` enables tetris-style packing
- Smaller boxes automatically fill gaps left by larger boxes
- Perfect vertical and horizontal alignment
- Example: A `medium` box fits 4 `small` boxes in a 2×2 grid with no gaps

## Size Calculation Formula

All box sizes follow this formula to align perfectly to the 13px grid:

```text
Box size = (N × 13px) + ((N-1) × 6px)
```

Where:

- `N` = number of grid cells (columns or rows)
- `13px` = base grid cell size
- `6px` = gap between boxes

### Examples

```text
dot:    (1 × 13) + (0 × 6) = 13px
micro:  (2 × 13) + (1 × 6) = 32px
small:  (4 × 13) + (3 × 6) = 70px
medium: (8 × 13) + (7 × 6) = 146px
large:  (16 × 13) + (15 × 6) = 298px
xlarge: (32 × 13) + (31 × 6) = 602px
```

## Box Sizes Reference

### Single (Square) Boxes

| Size    | Cells | Width×Height | Grid Span   |
|---------|-------|--------------|-------------|
| dot     | 1×1   | 13×13px      | 1 col × 1 row |
| micro   | 2×2   | 32×32px      | 2 col × 2 row |
| small   | 4×4   | 70×70px      | 4 col × 4 row |
| medium  | 8×8   | 146×146px    | 8 col × 8 row |
| large   | 16×16 | 298×298px    | 16 col × 16 row |
| xlarge  | 32×32 | 602×602px    | 32 col × 32 row |

A `dot` is a solid block of its status colour with no name or stat; four dots
fit in the space of one `micro`. Hovering shows the name in the status bar and
clicking opens the details as for any other box.

### Double-Width Boxes

Double boxes span 2× width but have single height:

| Size     | Cells  | Width×Height | Grid Span     |
|----------|--------|--------------|---------------|
| dmicro   | 4×2    | 70×32px      | 4 col × 2 row |
| dsmall   | 8×4    | 146×70px     | 8 col × 4 row |
| dmedium  | 16×8   | 298×146px    | 16 col × 8 row |
| dlarge   | 32×16  | 602×298px    | 32 col × 16 row |

## Perfect Tetris Packing Examples

### Example 1: Medium Box Equivalence

```text
1 medium box (8×8 cells) = 4 small boxes (4×4 each) in a 2×2 grid
1 medium box = 2 small boxes side-by-side (horizontally)
1 medium box = 2 small boxes stacked (vertically)
```

### Example 2: Stacking with Doubles

```text
1 medium box (146px wide, 146px tall) occupies same space as:
  - 1 dsmall (146px wide, 70px tall) on top
  - 1 small (70px wide, 70px tall) below
  - Empty space filled by dense grid packing
```

### Example 3: Large Box Grid

```text
1 large box (16×16 cells) = 4 medium boxes (8×8 each) in 2×2 grid
1 large box = 64 micro boxes (2×2 each) in 8×8 grid
1 micro box (2×2 cells) = 4 dot boxes (1×1 each) in 2×2 grid
```

## Grid Configuration

### CSS Grid Setup

```css
.boxes-grid {
    display: grid;
    gap: 6px;
    grid-auto-flow: dense;           /* Enable tetris packing */
    grid-auto-rows: 13px;            /* Row height = base cell */
    grid-template-columns: repeat(auto-fill, 13px 13px); /* Columns in pairs */
}
```

### Dense Packing Behaviour

The `grid-auto-flow: dense` property:

- Fills gaps with smaller boxes that fit
- Reorders boxes visually (not in DOM) for optimal packing
- Earlier boxes in DOM get priority placement
- Later boxes backfill available gaps

Columns are added in pairs so the column count is always even. Every size
except `dot` spans an even number of cells, so those boxes always start on an
even column and dots fill the remaining micro sized gaps four at a time, rather
than trickling down a single odd column at the right hand edge.

## Viewport Calculations

### Maximum Boxes Per Row

For a viewport width `W`:

```text
Available width = W - drawer_width - 16px (padding)
Grid columns = floor((Available_width + 6) / 19), rounded down to even
               (19 = 13px cell + 6px gap)
```

Example with 1900px viewport (no drawer):

```text
Available: 1900 - 0 - 16 = 1884px
Columns: floor(1890 / 19) = 99, rounded down to 98 columns

Max per row:
- xlarge: 98 ÷ 32 = 3 boxes
- large: 98 ÷ 16 = 6 boxes
- medium: 98 ÷ 8 = 12 boxes
- small: 98 ÷ 4 = 24 boxes
- micro: 98 ÷ 2 = 49 boxes
- dot: 98 ÷ 1 = 98 boxes
```

### With Drawer Open

When drawer opens (~400px wide):

```text
Available: 1900 - 400 - 16 = 1484px
Columns: floor(1490 / 19) = 78 columns

Max per row:
- xlarge: 78 ÷ 32 = 2 boxes
- small: 78 ÷ 4 = 19 boxes
```

## Alignment Guarantees

### Horizontal Alignment

All boxes align horizontally because all widths are multiples of (13px + 6px gap), less one gap.

### Vertical Alignment  

All boxes align vertically because all heights are multiples of (13px + 6px gap), less one gap.

### Status Bar Alignment

Status bar has `margin: 8px` to match boxes-grid `padding: 0 8px`, ensuring perfect alignment.

## Size Hierarchy

The doubling pattern ensures consistent scaling:

```text
dot (13)
  ↓ ×2 + gap
micro (32) = 2×13 + 6
  ↓ ×2 + gap
small (70) = 2×32 + 6
  ↓ ×2 + gap  
medium (146) = 2×70 + 6
  ↓ ×2 + gap
large (298) = 2×146 + 6
  ↓ ×2 + gap
xlarge (602) = 2×298 + 6
```

Each size is exactly twice the previous size plus one gap.

## Box Sizing Constraints

### Why These Specific Sizes?

1. **13px base**: Small enough for a `dot` while keeping `micro` at 32px (2×13 + 6)
2. **6px gap**: Visible but not excessive spacing between boxes
3. **Powers of 2**: Each size doubles the grid cells (1, 2, 4, 8, 16, 32)
4. **Perfect tiling**: Any combination of same-level boxes fills space exactly

### Modifying Sizes

If you need to change box sizes, follow these rules:

1. **Choose a base cell size** (e.g., 13px)
2. **Choose a gap size** (e.g., 6px)
3. **Use formula**: `Size = (N × base) + ((N-1) × gap)`
4. **Use powers of 2 for N**: 1, 2, 4, 8, 16, 32 for perfect doubling
5. **Update grid CSS**:

   ```css
   grid-auto-rows: [base]px;
   grid-template-columns: repeat(auto-fill, [base]px);
   gap: [gap]px;
   ```

### Example: 40px Base Grid

```text
micro:  (1 × 40) + (0 × 6) = 40px    → grid-column: span 1; grid-row: span 1;
small:  (2 × 40) + (1 × 6) = 86px    → grid-column: span 2; grid-row: span 2;
medium: (4 × 40) + (3 × 6) = 178px   → grid-column: span 4; grid-row: span 4;
large:  (8 × 40) + (7 × 6) = 362px   → grid-column: span 8; grid-row: span 8;
xlarge: (16 × 40) + (15 × 6) = 730px → grid-column: span 16; grid-row: span 16;
```

## Implementation Details

### Grid Span CSS

Each box size specifies its grid span:

```css
.small {
    grid-column: span 4;  /* Occupies 4 columns */
    grid-row: span 4;     /* Occupies 4 rows */
    height: 70px;
    width: 70px;
}
```

### Template Elements

Template elements containing box details are hidden:

```css
#boxes-grid > template {
    display: none;
}
```

### Sorting and Order

Boxes are sorted by:

1. Size (largest first: xlarge → micro)
2. Name (alphabetical within same size)

This ensures larger boxes place first, leaving gaps for smaller boxes to fill via dense packing.

## Troubleshooting

### Boxes Not Aligning

**Check:**

1. All box sizes follow the formula: `(N × 13) + ((N-1) × 6)`
2. Grid cell size matches: `grid-auto-rows: 13px`
3. Gap matches: `gap: 6px`
4. Grid spans match cell count: `grid-column: span N`

### Gaps Still Appearing

**Possible causes:**

1. Box dimensions don't match CSS grid spans
2. Template elements not hidden
3. Flexbox CSS still present (remove `float`, `flex` properties)
4. Browser doesn't support CSS Grid (check compatibility)

### Status Bar Misaligned

**Fix:**

- Status bar margin should equal boxes-grid padding
- Both should be 8px from viewport edge
- big-box width calculation should exclude scrollbar: use `document.documentElement.clientWidth`

## Performance Considerations

**CSS Grid vs Flexbox Performance:**

- Grid: Slightly more CPU for initial layout calculation
- Grid: Better rendering with many boxes (no reflow from wrapping)
- Dense packing: Minimal overhead for small to medium box counts (<500 boxes)
- Overall: Negligible performance difference for typical dashboard usage

## Future Extensions

### Adding New Box Sizes

To add intermediate sizes (e.g., "tiny" between micro and small):

1. Choose cell count: `N = 1.5` (non-integer) won't work perfectly
2. Use adjacent integer: `N = 1` (same as micro) - consider if needed
3. Or double existing: Use 2 micro boxes instead

**Recommendation**: Stick to powers-of-2 cell counts for perfect tiling.

### Responsive Breakpoints

Current implementation is fully responsive via `auto-fill` columns. To add breakpoints:

```css
@media (max-width: 1200px) {
    .boxes-grid {
        grid-template-columns: repeat(auto-fill, 13px 13px);
    }
    .xlarge {
        /* Optionally reduce xlarge to large at smaller screens */
        grid-column: span 16;
        grid-row: span 16;
    }
}
```
