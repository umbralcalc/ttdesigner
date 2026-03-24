package gamedata

// TileColor represents the phase a tile belongs to.
type TileColor int

const (
	TileColorYellow TileColor = iota
	TileColorGreen
	TileColorBrown
	TileColorGray
)

func (c TileColor) String() string {
	switch c {
	case TileColorYellow:
		return "yellow"
	case TileColorGreen:
		return "green"
	case TileColorBrown:
		return "brown"
	case TileColorGray:
		return "gray"
	default:
		return "unknown"
	}
}

// StopType describes what kind of stop (if any) a tile has.
type StopType int

const (
	StopNone StopType = iota
	StopTown           // small station, adds fixed revenue
	StopCity           // large station, has token slots
)

// TrackSegment represents a connection between two endpoints on a tile.
// Endpoints are hex edges (0-5) or a city/town node (represented as -1).
type TrackSegment struct {
	From int // 0-5 for hex edges, -1 for city/town
	To   int
}

// TileDef defines a tile's geometry and properties.
//
// All path/edge data is sourced from the 18xx.games project
// (https://github.com/tobymao/18xx) lib/engine/config/tile.rb.
//
// Edge numbering convention (flat-top hex):
//
//	0 = top, 1 = upper-right, 2 = lower-right,
//	3 = bottom, 4 = lower-left, 5 = upper-left
//
// Segments use -1 to represent the city/town node. A segment {0, -1}
// means "track from edge 0 to the city/town". A segment {0, 3} means
// "track from edge 0 to edge 3" with no stop.
//
// Upgrade rules (from 18xx.games engine):
//   - Color must advance exactly one step: yellow->green->brown->gray.
//   - The new tile must contain all paths of the old tile (subset check).
//   - City count and town count must match.
//   - Labels must match (e.g. "T" only upgrades to "T").
//   - Rotation (0-5) is chosen at placement to satisfy the subset check.
type TileDef struct {
	ID       int
	Color    TileColor
	Segments []TrackSegment // track connections
	Stop     StopType
	Slots    int    // token slots (only for cities)
	Revenue  int    // base revenue (for cities/towns)
	Label    string // special label (e.g. "K", "T", "H")

	// UpgradesTo lists tile IDs this tile can legally be replaced by.
	// Determined by color progression + path subset + city/town/label matching.
	UpgradesTo []int
}

// TileManifestEntry records how many copies of a tile are available.
type TileManifestEntry struct {
	TileID int
	Count  int
}

// Default1889TileManifest returns the available tile counts for 1889.
func Default1889TileManifest() []TileManifestEntry {
	return []TileManifestEntry{
		// Yellow
		{3, 2}, {5, 2}, {6, 2}, {7, 2}, {8, 5}, {9, 5}, {57, 2}, {58, 3},
		{437, 1}, {438, 1},
		// Green
		{12, 1}, {13, 1}, {14, 1}, {15, 3}, {16, 1}, {19, 1}, {20, 1},
		{23, 2}, {24, 2}, {25, 1}, {26, 1}, {27, 1}, {28, 1}, {29, 1},
		{205, 1}, {206, 1}, {439, 1}, {440, 1},
		// Brown
		{39, 1}, {40, 1}, {41, 1}, {42, 1}, {45, 1}, {46, 1}, {47, 1},
		{448, 4}, {465, 1}, {466, 1}, {492, 1}, {611, 2},
	}
}

// Default1889Tiles returns all tile definitions used in 1889.
//
// Source: 18xx.games lib/engine/config/tile.rb tile string notation:
//
//	path=a:E1,b:E2       — track from edge E1 to edge E2
//	path=a:E,b:_0        — track from edge E to node 0 (city/town)
//	city=revenue:R,slots:S — city with revenue R and S token slots
//	town=revenue:R       — town with revenue R
func Default1889Tiles() map[int]TileDef {
	tiles := map[int]TileDef{

		// =================================================================
		// YELLOW TILES
		// =================================================================

		// Tile 3: town, edges 0-1 through town
		// town=revenue:10;path=a:0,b:_0;path=a:_0,b:1
		3: {ID: 3, Color: TileColorYellow, Stop: StopTown, Revenue: 10,
			Segments:   []TrackSegment{{0, -1}, {-1, 1}},
			UpgradesTo: nil, // towns do not upgrade to cities; no standard green town tiles in 1889
		},

		// Tile 5: city 1 slot, edges 0 and 1
		// city=revenue:20;path=a:0,b:_0;path=a:1,b:_0
		5: {ID: 5, Color: TileColorYellow, Stop: StopCity, Slots: 1, Revenue: 20,
			Segments:   []TrackSegment{{0, -1}, {1, -1}},
			UpgradesTo: []int{12, 13, 14, 15, 205, 206},
		},

		// Tile 6: city 1 slot, edges 0 and 2
		// city=revenue:20;path=a:0,b:_0;path=a:2,b:_0
		6: {ID: 6, Color: TileColorYellow, Stop: StopCity, Slots: 1, Revenue: 20,
			Segments:   []TrackSegment{{0, -1}, {2, -1}},
			UpgradesTo: []int{12, 13, 14, 15, 205, 206},
		},

		// Tile 7: track only, edge 0 to edge 1 (tight curve)
		// path=a:0,b:1
		7: {ID: 7, Color: TileColorYellow, Stop: StopNone,
			Segments:   []TrackSegment{{0, 1}},
			UpgradesTo: []int{16, 19, 20, 23, 24, 26, 27},
		},

		// Tile 8: track only, edge 0 to edge 2 (gentle curve)
		// path=a:0,b:2
		8: {ID: 8, Color: TileColorYellow, Stop: StopNone,
			Segments:   []TrackSegment{{0, 2}},
			UpgradesTo: []int{16, 19, 23, 24, 25, 28, 29},
		},

		// Tile 9: track only, edge 0 to edge 3 (straight)
		// path=a:0,b:3
		9: {ID: 9, Color: TileColorYellow, Stop: StopNone,
			Segments:   []TrackSegment{{0, 3}},
			UpgradesTo: []int{19, 20, 23, 24, 26, 27},
		},

		// Tile 57: city 1 slot, edges 0 and 3 (straight through city)
		// city=revenue:20;path=a:0,b:_0;path=a:_0,b:3
		57: {ID: 57, Color: TileColorYellow, Stop: StopCity, Slots: 1, Revenue: 20,
			Segments:   []TrackSegment{{0, -1}, {-1, 3}},
			UpgradesTo: []int{12, 13, 14, 15, 205, 206},
		},

		// Tile 58: town, edges 0 and 2 through town (gentle curve)
		// town=revenue:10;path=a:0,b:_0;path=a:_0,b:2
		58: {ID: 58, Color: TileColorYellow, Stop: StopTown, Revenue: 10,
			Segments:   []TrackSegment{{0, -1}, {-1, 2}},
			UpgradesTo: nil, // towns do not upgrade to cities
		},

		// Tile 437: town with port, edges 0 and 2 (1889-specific)
		// town=revenue:30;path=a:0,b:_0;path=a:_0,b:2;icon=image:port
		437: {ID: 437, Color: TileColorYellow, Stop: StopTown, Revenue: 30,
			Segments:   []TrackSegment{{0, -1}, {-1, 2}},
			UpgradesTo: nil, // port town; no standard upgrade
		},

		// Tile 438: city 1 slot, label H, edges 0 and 2 (1889 Kotohira yellow)
		// city=revenue:40;path=a:0,b:_0;path=a:2,b:_0;label=H
		438: {ID: 438, Color: TileColorYellow, Stop: StopCity, Slots: 1, Revenue: 40, Label: "H",
			Segments:   []TrackSegment{{0, -1}, {2, -1}},
			UpgradesTo: []int{439},
		},

		// =================================================================
		// GREEN TILES
		// =================================================================

		// Tile 12: city 1 slot, edges 0-1-2
		// city=revenue:30;path=a:0,b:_0;path=a:1,b:_0;path=a:2,b:_0
		12: {ID: 12, Color: TileColorGreen, Stop: StopCity, Slots: 1, Revenue: 30,
			Segments:   []TrackSegment{{0, -1}, {1, -1}, {2, -1}},
			UpgradesTo: []int{448, 611},
		},

		// Tile 13: city 1 slot, edges 0-2-4
		// city=revenue:30;path=a:0,b:_0;path=a:2,b:_0;path=a:4,b:_0
		13: {ID: 13, Color: TileColorGreen, Stop: StopCity, Slots: 1, Revenue: 30,
			Segments:   []TrackSegment{{0, -1}, {2, -1}, {4, -1}},
			UpgradesTo: []int{448, 611},
		},

		// Tile 14: city 2 slots, edges 0-1-3-4
		// city=revenue:30,slots:2;path=a:0,b:_0;path=a:1,b:_0;path=a:3,b:_0;path=a:4,b:_0
		14: {ID: 14, Color: TileColorGreen, Stop: StopCity, Slots: 2, Revenue: 30,
			Segments:   []TrackSegment{{0, -1}, {1, -1}, {3, -1}, {4, -1}},
			UpgradesTo: []int{448, 611},
		},

		// Tile 15: city 2 slots, edges 0-1-2-3
		// city=revenue:30,slots:2;path=a:0,b:_0;path=a:1,b:_0;path=a:2,b:_0;path=a:3,b:_0
		15: {ID: 15, Color: TileColorGreen, Stop: StopCity, Slots: 2, Revenue: 30,
			Segments:   []TrackSegment{{0, -1}, {1, -1}, {2, -1}, {3, -1}},
			UpgradesTo: []int{448, 611},
		},

		// Tile 16: two parallel tracks, no stop (edges 0-2, 1-3)
		// path=a:0,b:2;path=a:1,b:3
		16: {ID: 16, Color: TileColorGreen, Stop: StopNone,
			Segments:   []TrackSegment{{0, 2}, {1, 3}},
			UpgradesTo: []int{45, 46},
		},

		// Tile 19: two tracks crossing, no stop (edges 0-3, 2-4)
		// path=a:0,b:3;path=a:2,b:4
		19: {ID: 19, Color: TileColorGreen, Stop: StopNone,
			Segments:   []TrackSegment{{0, 3}, {2, 4}},
			UpgradesTo: []int{45, 46},
		},

		// Tile 20: two tracks, no stop (edges 0-3, 1-4)
		// path=a:0,b:3;path=a:1,b:4
		20: {ID: 20, Color: TileColorGreen, Stop: StopNone,
			Segments:   []TrackSegment{{0, 3}, {1, 4}},
			UpgradesTo: []int{47},
		},

		// Tile 23: Y-junction, no stop (edges 0-3, 0-4)
		// path=a:0,b:3;path=a:0,b:4
		23: {ID: 23, Color: TileColorGreen, Stop: StopNone,
			Segments:   []TrackSegment{{0, 3}, {0, 4}},
			UpgradesTo: []int{41, 42, 45, 46, 47},
		},

		// Tile 24: Y-junction, no stop (edges 0-3, 0-2)
		// path=a:0,b:3;path=a:0,b:2
		24: {ID: 24, Color: TileColorGreen, Stop: StopNone,
			Segments:   []TrackSegment{{0, 3}, {0, 2}},
			UpgradesTo: []int{41, 42, 45, 46, 47},
		},

		// Tile 25: Y-junction, no stop (edges 0-2, 0-4)
		// path=a:0,b:2;path=a:0,b:4
		25: {ID: 25, Color: TileColorGreen, Stop: StopNone,
			Segments:   []TrackSegment{{0, 2}, {0, 4}},
			UpgradesTo: []int{40, 45, 46},
		},

		// Tile 26: Y-junction, no stop (edges 0-3, 0-5)
		// path=a:0,b:3;path=a:0,b:5
		26: {ID: 26, Color: TileColorGreen, Stop: StopNone,
			Segments:   []TrackSegment{{0, 3}, {0, 5}},
			UpgradesTo: []int{42, 45, 46, 47},
		},

		// Tile 27: Y-junction, no stop (edges 0-3, 0-1)
		// path=a:0,b:3;path=a:0,b:1
		27: {ID: 27, Color: TileColorGreen, Stop: StopNone,
			Segments:   []TrackSegment{{0, 3}, {0, 1}},
			UpgradesTo: []int{41, 45, 46, 47},
		},

		// Tile 28: Y-junction, no stop (edges 0-4, 0-5)
		// path=a:0,b:4;path=a:0,b:5
		28: {ID: 28, Color: TileColorGreen, Stop: StopNone,
			Segments:   []TrackSegment{{0, 4}, {0, 5}},
			UpgradesTo: []int{39, 40, 45, 46},
		},

		// Tile 29: Y-junction, no stop (edges 0-2, 0-1)
		// path=a:0,b:2;path=a:0,b:1
		29: {ID: 29, Color: TileColorGreen, Stop: StopNone,
			Segments:   []TrackSegment{{0, 2}, {0, 1}},
			UpgradesTo: []int{39, 40, 45, 46},
		},

		// Tile 205: city 1 slot, edges 0-1-3
		// city=revenue:30;path=a:0,b:_0;path=a:1,b:_0;path=a:3,b:_0
		205: {ID: 205, Color: TileColorGreen, Stop: StopCity, Slots: 1, Revenue: 30,
			Segments:   []TrackSegment{{0, -1}, {1, -1}, {3, -1}},
			UpgradesTo: []int{448, 611},
		},

		// Tile 206: city 1 slot, edges 0-3-5
		// city=revenue:30;path=a:0,b:_0;path=a:5,b:_0;path=a:3,b:_0
		206: {ID: 206, Color: TileColorGreen, Stop: StopCity, Slots: 1, Revenue: 30,
			Segments:   []TrackSegment{{0, -1}, {5, -1}, {3, -1}},
			UpgradesTo: []int{448, 611},
		},

		// Tile 439: city 2 slots, label H, edges 0-2-4 (1889 Kotohira green)
		// city=revenue:60,slots:2;path=a:0,b:_0;path=a:2,b:_0;path=a:4,b:_0;label=H
		439: {ID: 439, Color: TileColorGreen, Stop: StopCity, Slots: 2, Revenue: 60, Label: "H",
			Segments:   []TrackSegment{{0, -1}, {2, -1}, {4, -1}},
			UpgradesTo: []int{492},
		},

		// Tile 440: city 2 slots, label T, edges 0-1-2 (1889 Takamatsu green)
		// city=revenue:40,slots:2;path=a:0,b:_0;path=a:1,b:_0;path=a:2,b:_0;label=T
		440: {ID: 440, Color: TileColorGreen, Stop: StopCity, Slots: 2, Revenue: 40, Label: "T",
			Segments:   []TrackSegment{{0, -1}, {1, -1}, {2, -1}},
			UpgradesTo: []int{466},
		},

		// =================================================================
		// BROWN TILES
		// =================================================================

		// Tile 39: triangle track, no stop (edges 0-2, 0-1, 1-2)
		// path=a:0,b:2;path=a:0,b:1;path=a:1,b:2
		39: {ID: 39, Color: TileColorBrown, Stop: StopNone,
			Segments: []TrackSegment{{0, 2}, {0, 1}, {1, 2}},
		},

		// Tile 40: wide triangle track, no stop (edges 0-2, 2-4, 0-4)
		// path=a:0,b:2;path=a:2,b:4;path=a:0,b:4
		40: {ID: 40, Color: TileColorBrown, Stop: StopNone,
			Segments: []TrackSegment{{0, 2}, {2, 4}, {0, 4}},
		},

		// Tile 41: triangle track, no stop (edges 0-3, 0-1, 1-3)
		// path=a:0,b:3;path=a:0,b:1;path=a:1,b:3
		41: {ID: 41, Color: TileColorBrown, Stop: StopNone,
			Segments: []TrackSegment{{0, 3}, {0, 1}, {1, 3}},
		},

		// Tile 42: triangle track, no stop (edges 0-3, 3-5, 0-5)
		// path=a:0,b:3;path=a:3,b:5;path=a:0,b:5
		42: {ID: 42, Color: TileColorBrown, Stop: StopNone,
			Segments: []TrackSegment{{0, 3}, {3, 5}, {0, 5}},
		},

		// Tile 45: 4-way track, no stop (edges 0-3, 2-4, 0-4, 2-3)
		// path=a:0,b:3;path=a:2,b:4;path=a:0,b:4;path=a:2,b:3
		45: {ID: 45, Color: TileColorBrown, Stop: StopNone,
			Segments: []TrackSegment{{0, 3}, {2, 4}, {0, 4}, {2, 3}},
		},

		// Tile 46: 4-way track, no stop (edges 0-3, 2-4, 3-4, 0-2)
		// path=a:0,b:3;path=a:2,b:4;path=a:3,b:4;path=a:0,b:2
		46: {ID: 46, Color: TileColorBrown, Stop: StopNone,
			Segments: []TrackSegment{{0, 3}, {2, 4}, {3, 4}, {0, 2}},
		},

		// Tile 47: 4-way track, no stop (edges 0-3, 1-4, 1-3, 0-4)
		// path=a:0,b:3;path=a:1,b:4;path=a:1,b:3;path=a:0,b:4
		47: {ID: 47, Color: TileColorBrown, Stop: StopNone,
			Segments: []TrackSegment{{0, 3}, {1, 4}, {1, 3}, {0, 4}},
		},

		// Tile 448: city 2 slots, edges 0-1-2-3
		// city=revenue:40,slots:2;path=a:0,b:_0;path=a:1,b:_0;path=a:2,b:_0;path=a:3,b:_0
		448: {ID: 448, Color: TileColorBrown, Stop: StopCity, Slots: 2, Revenue: 40,
			Segments: []TrackSegment{{0, -1}, {1, -1}, {2, -1}, {3, -1}},
		},

		// Tile 465: city 3 slots, label K, edges 0-1-2-3 (1889 Kouchi brown)
		// city=revenue:60,slots:3;path=a:0,b:_0;path=a:1,b:_0;path=a:2,b:_0;path=a:3,b:_0;label=K
		465: {ID: 465, Color: TileColorBrown, Stop: StopCity, Slots: 3, Revenue: 60, Label: "K",
			Segments: []TrackSegment{{0, -1}, {1, -1}, {2, -1}, {3, -1}},
		},

		// Tile 466: city 2 slots, label T, edges 0-1-2 (1889 Takamatsu brown)
		// city=revenue:60,slots:2;path=a:0,b:_0;path=a:1,b:_0;path=a:2,b:_0;label=T
		466: {ID: 466, Color: TileColorBrown, Stop: StopCity, Slots: 2, Revenue: 60, Label: "T",
			Segments: []TrackSegment{{0, -1}, {1, -1}, {2, -1}},
		},

		// Tile 492: city 3 slots, label H, all 6 edges (1889 Kotohira brown)
		// city=revenue:80,slots:3;path=a:0,b:_0;...;path=a:5,b:_0;label=H
		492: {ID: 492, Color: TileColorBrown, Stop: StopCity, Slots: 3, Revenue: 80, Label: "H",
			Segments: []TrackSegment{{0, -1}, {1, -1}, {2, -1}, {3, -1}, {4, -1}, {5, -1}},
		},

		// Tile 611: city 2 slots, edges 0-1-2-3-4
		// city=revenue:40,slots:2;path=a:0,b:_0;path=a:1,b:_0;path=a:2,b:_0;path=a:3,b:_0;path=a:4,b:_0
		611: {ID: 611, Color: TileColorBrown, Stop: StopCity, Slots: 2, Revenue: 40,
			Segments: []TrackSegment{{0, -1}, {1, -1}, {2, -1}, {3, -1}, {4, -1}},
		},
	}
	return tiles
}
