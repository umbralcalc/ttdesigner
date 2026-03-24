package gamedata

// MarketZone classifies a cell's special behavior.
type MarketZone int

const (
	MarketZoneNormal MarketZone = iota
	MarketZonePar               // valid par price
	MarketZoneYellow            // must sell: shares don't count toward cert limit
	MarketZoneOrange            // company closes / bankrupt
)

// MarketCell represents one cell in the stock market grid.
type MarketCell struct {
	Price int
	Zone  MarketZone
}

// MarketGrid is the 2D stock market. Row 0 is the top (highest prices).
type MarketGrid struct {
	Cells [][]MarketCell
	Rows  int
	Cols  int
}

// Price returns the share price at the given (row, col), or 0 if out of bounds.
func (g *MarketGrid) Price(row, col int) int {
	if row < 0 || row >= g.Rows || col < 0 || col >= len(g.Cells[row]) {
		return 0
	}
	return g.Cells[row][col].Price
}

// Zone returns the zone at the given (row, col).
func (g *MarketGrid) Zone(row, col int) MarketZone {
	if row < 0 || row >= g.Rows || col < 0 || col >= len(g.Cells[row]) {
		return MarketZoneNormal
	}
	return g.Cells[row][col].Zone
}

// MoveRight returns the new (row, col) after paying dividends.
// Moves one column right; clamps to rightmost column in the row.
func (g *MarketGrid) MoveRight(row, col int) (int, int) {
	newCol := col + 1
	if newCol >= len(g.Cells[row]) {
		newCol = len(g.Cells[row]) - 1
	}
	return row, newCol
}

// MoveLeft returns the new (row, col) after withholding dividends.
// Moves one column left; clamps to column 0.
func (g *MarketGrid) MoveLeft(row, col int) (int, int) {
	newCol := col - 1
	if newCol < 0 {
		newCol = 0
	}
	return row, newCol
}

// MoveDown returns the new (row, col) after shares are sold.
// Moves one row down; if the new row is shorter, clamps to rightmost column.
// Clamps to bottom row.
func (g *MarketGrid) MoveDown(row, col int) (int, int) {
	newRow := row + 1
	if newRow >= g.Rows {
		newRow = g.Rows - 1
	}
	if col >= len(g.Cells[newRow]) {
		col = len(g.Cells[newRow]) - 1
	}
	return newRow, col
}

// ParValues returns all valid par prices (cells in the par zone) as (row, col, price) triples.
func (g *MarketGrid) ParValues() []struct {
	Row, Col, Price int
} {
	var pars []struct{ Row, Col, Price int }
	for r, row := range g.Cells {
		for c, cell := range row {
			if cell.Zone == MarketZonePar {
				pars = append(pars, struct{ Row, Col, Price int }{r, c, cell.Price})
			}
		}
	}
	return pars
}

// cell is a shorthand constructor.
func cell(price int, zone MarketZone) MarketCell {
	return MarketCell{Price: price, Zone: zone}
}

// Default1889Market returns the stock market grid for 1889.
func Default1889Market() *MarketGrid {
	n := MarketZoneNormal
	p := MarketZonePar
	y := MarketZoneYellow
	o := MarketZoneOrange

	cells := [][]MarketCell{
		// Row 0
		{cell(75, n), cell(80, n), cell(90, n), cell(100, p), cell(110, n), cell(125, n), cell(140, n), cell(155, n), cell(175, n), cell(200, n), cell(225, n), cell(255, n), cell(285, n), cell(315, n), cell(350, n)},
		// Row 1
		{cell(70, n), cell(75, n), cell(80, n), cell(90, p), cell(100, n), cell(110, n), cell(125, n), cell(140, n), cell(155, n), cell(175, n), cell(200, n), cell(225, n), cell(255, n), cell(285, n), cell(315, n)},
		// Row 2
		{cell(65, n), cell(70, n), cell(75, n), cell(80, p), cell(90, n), cell(100, n), cell(110, n), cell(125, n), cell(140, n), cell(155, n), cell(175, n), cell(200, n)},
		// Row 3
		{cell(60, n), cell(65, n), cell(70, n), cell(75, p), cell(80, n), cell(90, n), cell(100, n), cell(110, n), cell(125, n), cell(140, n)},
		// Row 4
		{cell(55, n), cell(60, n), cell(65, n), cell(70, p), cell(75, n), cell(80, n), cell(90, n), cell(100, n)},
		// Row 5
		{cell(50, y), cell(55, n), cell(60, n), cell(65, p), cell(70, n), cell(75, n), cell(80, n)},
		// Row 6
		{cell(45, y), cell(50, y), cell(55, n), cell(60, n), cell(65, n), cell(70, n)},
		// Row 7
		{cell(40, y), cell(45, y), cell(50, y), cell(55, n), cell(60, n)},
		// Row 8
		{cell(30, o), cell(40, y), cell(45, y), cell(50, y)},
		// Row 9
		{cell(20, o), cell(30, o), cell(40, y), cell(45, y)},
		// Row 10
		{cell(10, o), cell(20, o), cell(30, o), cell(40, y)},
	}

	return &MarketGrid{
		Cells: cells,
		Rows:  len(cells),
		Cols:  15, // max width (row 0 and 1)
	}
}
