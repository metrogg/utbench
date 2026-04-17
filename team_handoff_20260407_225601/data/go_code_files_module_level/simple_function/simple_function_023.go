func (t *Table) Draw(screen tcell.Screen) {
	t.Box.Draw(screen)

	// What's our available screen space?
	x, y, width, height := t.GetInnerRect()
	if t.borders {
		t.visibleRows = height / 2
	} else {
		t.visibleRows = height
	}

	// Return the cell at the specified position (nil if it doesn't exist).
	getCell := func(row, column int) *TableCell {
		if row < 0 || column < 0 || row >= len(t.cells) || column >= len(t.cells[row]) {
			return nil
		}
		return t.cells[row][column]
	}

	// If this cell is not selectable, find the next one.
	if t.rowsSelectable || t.columnsSelectable {
		if t.selectedColumn < 0 {
			t.selectedColumn = 0
		}
		if t.selectedRow < 0 {
			t.selectedRow = 0
		}
		for t.selectedRow < len(t.cells) {
			cell := getCell(t.selectedRow, t.selectedColumn)
			if cell == nil || !cell.NotSelectable {
				break
			}
			t.selectedColumn++
			if t.selectedColumn > t.lastColumn {
				t.selectedColumn = 0
				t.selectedRow++
			}
		}
	}

	// Clamp row offsets.
	if t.rowsSelectable {
		if t.selectedRow >= t.fixedRows && t.selectedRow < t.fixedRows+t.rowOffset {
			t.rowOffset = t.selectedRow - t.fixedRows
			t.trackEnd = false
		}
		if t.borders {
			if 2*(t.selectedRow+1-t.rowOffset) >= height {
				t.rowOffset = t.selectedRow + 1 - height/2
				t.trackEnd = false
			}
		} else {
			if t.selectedRow+1-t.rowOffset >= height {
				t.rowOffset = t.selectedRow + 1 - height
				t.trackEnd = false
			}
		}
	}
	if t.borders {
		if 2*(len(t.cells)-t.rowOffset) < height {
			t.trackEnd = true
		}
	} else {
		if len(t.cells)-t.rowOffset < height {
			t.trackEnd = true
		}
	}
	if t.trackEnd {
		if t.borders {
			t.rowOffset = len(t.cells) - height/2
		} else {
			t.rowOffset = len(t.cells) - height
		}
	}
	if t.rowOffset < 0 {
		t.rowOffset = 0
	}

	// Clamp column offset. (Only left side here. The right side is more
	// difficult and we'll do it below.)
	if t.columnsSelectable && t.selectedColumn >= t.fixedColumns && t.selectedColumn < t.fixedColumns+t.columnOffset {
		t.columnOffset = t.selectedColumn - t.fixedColumns
	}
	if t.columnOffset < 0 {
		t.columnOffset = 0
	}
	if t.selectedColumn < 0 {
		t.selectedColumn = 0
	}

	// Determine the indices and widths of the columns and rows which fit on the
	// screen.
	var (
		columns, rows, widths   []int
		tableHeight, tableWidth int
	)
	rowStep := 1
	if t.borders {
		rowStep = 2    // With borders, every table row takes two screen rows.
		tableWidth = 1 // We start at the second character because of the left table border.
	}
	indexRow := func(row int) bool { // Determine if this row is visible, store its index.
		if tableHeight >= height {
			return false
		}
		rows = append(rows, row)
		tableHeight += rowStep
		return true
	}
	for row := 0; row < t.fixedRows && row < len(t.cells); row++ { // Do the fixed rows first.
		if !indexRow(row) {
			break
		}
	}
	for row := t.fixedRows + t.rowOffset; row < len(t.cells); row++ { // Then the remaining rows.
		if !indexRow(row) {
			break
		}
	}
	var (
		skipped, lastTableWidth, expansionTotal int
		expansions                              []int
	)
ColumnLoop:
	for column := 0; ; column++ {
		// If we've moved beyond the right border, we stop or skip a column.
		for tableWidth-1 >= width { // -1 because we include one extra column if the separator falls on the right end of the box.
			// We've moved beyond the available space.
			if column < t.fixedColumns {
				break ColumnLoop // We're in the fixed area. We're done.
			}
			if !t.columnsSelectable && skipped >= t.columnOffset {
				break ColumnLoop // There is no selection and we've already reached the offset.
			}
			if t.columnsSelectable && t.selectedColumn-skipped == t.fixedColumns {
				break ColumnLoop // The selected column reached the leftmost point before disappearing.
			}
			if t.columnsSelectable && skipped >= t.columnOffset &&
				(t.selectedColumn < column && lastTableWidth < width-1 && tableWidth < width-1 || t.selectedColumn < column-1) {
				break ColumnLoop // We've skipped as many as requested and the selection is visible.
			}
			if len(columns) <= t.fixedColumns {
				break // Nothing to skip.
			}

			// We need to skip a column.
			skipped++
			lastTableWidth -= widths[t.fixedColumns] + 1
			tableWidth -= widths[t.fixedColumns] + 1
			columns = append(columns[:t.fixedColumns], columns[t.fixedColumns+1:]...)
			widths = append(widths[:t.fixedColumns], widths[t.fixedColumns+1:]...)
			expansions = append(expansions[:t.fixedColumns], expansions[t.fixedColumns+1:]...)
		}

		// What's this column's width (without expansion)?
		maxWidth := -1
		expansion := 0
		for _, row := range rows {
			if cell := getCell(row, column); cell != nil {
				_, _, _, _, _, _, cellWidth := decomposeString(cell.Text, true, false)
				if cell.MaxWidth > 0 && cell.MaxWidth < cellWidth {
					cellWidth = cell.MaxWidth
				}
				if cellWidth > maxWidth {
					maxWidth = cellWidth
				}
				if cell.Expansion > expansion {
					expansion = cell.Expansion
				}
			}
		}
		if maxWidth < 0 {
			break // No more cells found in this column.
		}

		// Store new column info at the end.
		columns = append(columns, column)
		widths = append(widths, maxWidth)
		lastTableWidth = tableWidth
		tableWidth += maxWidth + 1
		expansions = append(expansions, expansion)
		expansionTotal += expansion
	}
	t.columnOffset = skipped

	// If we have space left, distribute it.
	if tableWidth < width {
		toDistribute := width - tableWidth
		for index, expansion := range expansions {
			if expansionTotal <= 0 {
				break
			}
			expWidth := toDistribute * expansion / expansionTotal
			widths[index] += expWidth
			toDistribute -= expWidth
			expansionTotal -= expansion
		}
	}

	// Helper function which draws border runes.
	borderStyle := tcell.StyleDefault.Background(t.backgroundColor).Foreground(t.bordersColor)
	drawBorder := func(colX, rowY int, ch rune) {
		screen.SetContent(x+colX, y+rowY, ch, nil, borderStyle)
	}

	// Draw the cells (and borders).
	var columnX int
	if !t.borders {
		columnX--
	}
	for columnIndex, column := range columns {
		columnWidth := widths[columnIndex]
		for rowY, row := range rows {
			if t.borders {
				// Draw borders.
				rowY *= 2
				for pos := 0; pos < columnWidth && columnX+1+pos < width; pos++ {
					drawBorder(columnX+pos+1, rowY, Borders.Horizontal)
				}
				ch := Borders.Cross
				if columnIndex == 0 {
					if rowY == 0 {
						ch = Borders.TopLeft
					} else {
						ch = Borders.LeftT
					}
				} else if rowY == 0 {
					ch = Borders.TopT
				}
				drawBorder(columnX, rowY, ch)
				rowY++
				if rowY >= height {
					break // No space for the text anymore.
				}
				drawBorder(columnX, rowY, Borders.Vertical)
			} else if columnIndex > 0 {
				// Draw separator.
				drawBorder(columnX, rowY, t.separator)
			}

			// Get the cell.
			cell := getCell(row, column)
			if cell == nil {
				continue
			}

			// Draw text.
			finalWidth := columnWidth
			if columnX+1+columnWidth >= width {
				finalWidth = width - columnX - 1
			}
			cell.x, cell.y, cell.width = x+columnX+1, y+rowY, finalWidth
			_, printed := printWithStyle(screen, cell.Text, x+columnX+1, y+rowY, finalWidth, cell.Align, tcell.StyleDefault.Foreground(cell.Color)|tcell.Style(cell.Attributes))
			if TaggedStringWidth(cell.Text)-printed > 0 && printed > 0 {
				_, _, style, _ := screen.GetContent(x+columnX+1+finalWidth-1, y+rowY)
				printWithStyle(screen, string(SemigraphicsHorizontalEllipsis), x+columnX+1+finalWidth-1, y+rowY, 1, AlignLeft, style)
			}
		}

		// Draw bottom border.
		if rowY := 2 * len(rows); t.borders && rowY < height {
			for pos := 0; pos < columnWidth && columnX+1+pos < width; pos++ {
				drawBorder(columnX+pos+1, rowY, Borders.Horizontal)
			}
			ch := Borders.BottomT
			if columnIndex == 0 {
				ch = Borders.BottomLeft
			}
			drawBorder(columnX, rowY, ch)
		}

		columnX += columnWidth + 1
	}

	// Draw right border.
	if t.borders && len(t.cells) > 0 && columnX < width {
		for rowY := range rows {
			rowY *= 2
			if rowY+1 < height {
				drawBorder(columnX, rowY+1, Borders.Vertical)
			}
			ch := Borders.RightT
			if rowY == 0 {
				ch = Borders.TopRight
			}
			drawBorder(columnX, rowY, ch)
		}
		if rowY := 2 * len(rows); rowY < height {
			drawBorder(columnX, rowY, Borders.BottomRight)
		}
	}

	// Helper function which colors the background of a box.
	// backgroundColor == tcell.ColorDefault => Don't color the background.
	// textColor == tcell.ColorDefault => Don't change the text color.
	// attr == 0 => Don't change attributes.
	// invert == true => Ignore attr, set text to backgroundColor or t.backgroundColor;
	//                   set background to textColor.
	colorBackground := func(fromX, fromY, w, h int, backgroundColor, textColor tcell.Color, attr tcell.AttrMask, invert bool) {
		for by := 0; by < h && fromY+by < y+height; by++ {
			for bx := 0; bx < w && fromX+bx < x+width; bx++ {
				m, c, style, _ := screen.GetContent(fromX+bx, fromY+by)
				fg, bg, a := style.Decompose()
				if invert {
					if fg == textColor || fg == t.bordersColor {
						fg = backgroundColor
					}
					if fg == tcell.ColorDefault {
						fg = t.backgroundColor
					}
					style = style.Background(textColor).Foreground(fg)
				} else {
					if backgroundColor != tcell.ColorDefault {
						bg = backgroundColor
					}
					if textColor != tcell.ColorDefault {
						fg = textColor
					}
					if attr != 0 {
						a = attr
					}
					style = style.Background(bg).Foreground(fg) | tcell.Style(a)
				}
				screen.SetContent(fromX+bx, fromY+by, m, c, style)
			}
		}
	}

	// Color the cell backgrounds. To avoid undesirable artefacts, we combine
	// the drawing of a cell by background color, selected cells last.
	type cellInfo struct {
		x, y, w, h int
		text       tcell.Color
		selected   bool
	}
	cellsByBackgroundColor := make(map[tcell.Color][]*cellInfo)
	var backgroundColors []tcell.Color
	for rowY, row := range rows {
		columnX := 0
		rowSelected := t.rowsSelectable && !t.columnsSelectable && row == t.selectedRow
		for columnIndex, column := range columns {
			columnWidth := widths[columnIndex]
			cell := getCell(row, column)
			if cell == nil {
				continue
			}
			bx, by, bw, bh := x+columnX, y+rowY, columnWidth+1, 1
			if t.borders {
				by = y + rowY*2
				bw++
				bh = 3
			}
			columnSelected := t.columnsSelectable && !t.rowsSelectable && column == t.selectedColumn
			cellSelected := !cell.NotSelectable && (columnSelected || rowSelected || t.rowsSelectable && t.columnsSelectable && column == t.selectedColumn && row == t.selectedRow)
			entries, ok := cellsByBackgroundColor[cell.BackgroundColor]
			cellsByBackgroundColor[cell.BackgroundColor] = append(entries, &cellInfo{
				x:        bx,
				y:        by,
				w:        bw,
				h:        bh,
				text:     cell.Color,
				selected: cellSelected,
			})
			if !ok {
				backgroundColors = append(backgroundColors, cell.BackgroundColor)
			}
			columnX += columnWidth + 1
		}
	}
	sort.Slice(backgroundColors, func(i int, j int) bool {
		// Draw brightest colors last (i.e. on top).
		r, g, b := backgroundColors[i].RGB()
		c := colorful.Color{R: float64(r) / 255, G: float64(g) / 255, B: float64(b) / 255}
		_, _, li := c.Hcl()
		r, g, b = backgroundColors[j].RGB()
		c = colorful.Color{R: float64(r) / 255, G: float64(g) / 255, B: float64(b) / 255}
		_, _, lj := c.Hcl()
		return li < lj
	})
	selFg, selBg, selAttr := t.selectedStyle.Decompose()
	for _, bgColor := range backgroundColors {
		entries := cellsByBackgroundColor[bgColor]
		for _, cell := range entries {
			if cell.selected {
				if t.selectedStyle != 0 {
					defer colorBackground(cell.x, cell.y, cell.w, cell.h, selBg, selFg, selAttr, false)
				} else {
					defer colorBackground(cell.x, cell.y, cell.w, cell.h, bgColor, cell.text, 0, true)
				}
			} else {
				colorBackground(cell.x, cell.y, cell.w, cell.h, bgColor, tcell.ColorDefault, 0, false)
			}
		}
	}
}
