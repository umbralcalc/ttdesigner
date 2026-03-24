package gamedata

import "fmt"

// TerrainType describes terrain features that affect building costs.
type TerrainType int

const (
	TerrainNone     TerrainType = iota
	TerrainMountain             // cost 80
	TerrainWater                // cost 80
	TerrainMountainWater        // both mountain and water, cost 80
)

// HexType describes what kind of hex this is.
type HexType int

const (
	HexEmpty    HexType = iota // buildable, no stop
	HexCity                    // large station with token slots
	HexTown                    // small station (halt)
	HexOffBoard                // red off-board location
	HexGray                    // permanent, cannot be upgraded
)

// OffBoardRevenue defines revenue that scales with game phase.
type OffBoardRevenue struct {
	Yellow int
	Brown  int
	Diesel int
}

// HexDef defines a hex on the game map.
type HexDef struct {
	ID       string // coordinate, e.g. "K8"
	Name     string // location name, e.g. "Tokushima"
	Type     HexType
	Terrain  TerrainType
	Cost     int    // extra build cost from terrain
	Slots    int    // token slots (cities only)
	Revenue  int    // base revenue (0 for empty)
	Label    string // special label: "K" (Kouchi), "T" (Takamatsu), "H" (Kotohira)

	// PrePrintedTile: tile ID pre-printed on the board (-1 = none).
	PrePrintedTile int

	// OffBoard revenue (only for HexOffBoard).
	OffBoard *OffBoardRevenue

	// Exits lists which hex edges (0-5) connect to adjacent hexes.
	// For off-board hexes, this constrains which edges have track.
	Exits []int
}

// HexAdjacency maps a hex ID to its neighbor hex IDs indexed by edge (0-5).
// A nil/empty string means no neighbor on that edge (map boundary).
type HexAdjacency map[string][6]string

// Default1889Map returns all hex definitions for the 1889 Shikoku map.
func Default1889Map() []HexDef {
	return []HexDef{
		// === OFF-BOARD (Red) ===
		{ID: "F1", Name: "Imabari", Type: HexOffBoard, OffBoard: &OffBoardRevenue{30, 60, 100}, Exits: []int{0, 1}},
		{ID: "J1", Name: "Sakaide & Okayama", Type: HexOffBoard, OffBoard: &OffBoardRevenue{20, 40, 80}, Exits: []int{0, 1}},
		{ID: "L7", Name: "Naruto & Awaji", Type: HexOffBoard, OffBoard: &OffBoardRevenue{20, 40, 80}, Exits: []int{1, 2}},

		// === GRAY (permanent, no upgrade) ===
		{ID: "B3", Name: "Yawatahama", Type: HexGray, Revenue: 20},
		{ID: "B7", Name: "Uwajima", Type: HexGray, Slots: 2, Revenue: 40},
		{ID: "G14", Name: "Muroto", Type: HexGray, Revenue: 20},
		{ID: "J7", Name: "", Type: HexGray}, // pass-through, no stop

		// === PRE-PRINTED YELLOW ===
		{ID: "C4", Name: "Ohzu", Type: HexCity, Slots: 1, Revenue: 20, PrePrintedTile: -1}, // starts yellow
		{ID: "K4", Name: "Takamatsu", Type: HexCity, Slots: 1, Revenue: 30, Label: "T", PrePrintedTile: -1},

		// === PRE-PRINTED GREEN ===
		{ID: "F9", Name: "Kouchi", Type: HexCity, Slots: 2, Revenue: 30, Label: "K", Cost: 80, PrePrintedTile: -1},

		// === CITIES (white, buildable) ===
		{ID: "E2", Name: "Matsuyama", Type: HexCity, Slots: 1},
		{ID: "F3", Name: "Saijou", Type: HexCity, Slots: 1},
		{ID: "G4", Name: "Niihama", Type: HexCity, Slots: 1},
		{ID: "H7", Name: "Ikeda", Type: HexCity, Slots: 1},
		{ID: "I2", Name: "Marugame", Type: HexCity, Slots: 1},
		{ID: "K8", Name: "Tokushima", Type: HexCity, Slots: 1},
		{ID: "A10", Name: "Sukumo", Type: HexCity, Slots: 1},
		{ID: "C10", Name: "Kubokawa", Type: HexCity, Slots: 1},
		{ID: "J11", Name: "Anan", Type: HexCity, Slots: 1},
		{ID: "G12", Name: "Nahari", Type: HexCity, Slots: 1},

		// === TOWNS ===
		{ID: "J5", Name: "Ritsurin Kouen", Type: HexTown},
		{ID: "B11", Name: "Nakamura", Type: HexTown},   // port
		{ID: "G10", Name: "Nangoku", Type: HexTown},     // port
		{ID: "I12", Name: "Muki", Type: HexTown},        // port
		{ID: "J9", Name: "Komatsujima", Type: HexTown},  // port

		// === KOTOHIRA (city with label H and upgrade cost) ===
		{ID: "I4", Name: "Kotohira", Type: HexCity, Slots: 1, Label: "H", Cost: 80},

		// === MOUNTAIN TERRAIN ===
		{ID: "E4", Name: "", Type: HexEmpty, Terrain: TerrainMountain, Cost: 80},
		{ID: "D5", Name: "", Type: HexEmpty, Terrain: TerrainMountain, Cost: 80},
		{ID: "F5", Name: "", Type: HexEmpty, Terrain: TerrainMountain, Cost: 80},
		{ID: "C6", Name: "", Type: HexEmpty, Terrain: TerrainMountain, Cost: 80},
		{ID: "E6", Name: "", Type: HexEmpty, Terrain: TerrainMountain, Cost: 80},
		{ID: "G6", Name: "", Type: HexEmpty, Terrain: TerrainMountain, Cost: 80},
		{ID: "D7", Name: "", Type: HexEmpty, Terrain: TerrainMountain, Cost: 80},
		{ID: "F7", Name: "", Type: HexEmpty, Terrain: TerrainMountain, Cost: 80},
		{ID: "A8", Name: "", Type: HexEmpty, Terrain: TerrainMountain, Cost: 80},
		{ID: "G8", Name: "", Type: HexEmpty, Terrain: TerrainMountain, Cost: 80},
		{ID: "B9", Name: "", Type: HexEmpty, Terrain: TerrainMountain, Cost: 80},
		{ID: "H9", Name: "", Type: HexEmpty, Terrain: TerrainMountain, Cost: 80},
		{ID: "H11", Name: "", Type: HexEmpty, Terrain: TerrainMountain, Cost: 80},
		{ID: "H13", Name: "", Type: HexEmpty, Terrain: TerrainMountain, Cost: 80},

		// === WATER TERRAIN ===
		{ID: "K6", Name: "", Type: HexEmpty, Terrain: TerrainWater, Cost: 80},

		// === MOUNTAIN + WATER TERRAIN ===
		{ID: "H5", Name: "", Type: HexEmpty, Terrain: TerrainMountainWater, Cost: 80},
		{ID: "I6", Name: "", Type: HexEmpty, Terrain: TerrainMountainWater, Cost: 80},

		// === PLAIN EMPTY ===
		{ID: "D3", Name: "", Type: HexEmpty},
		{ID: "H3", Name: "", Type: HexEmpty},
		{ID: "J3", Name: "", Type: HexEmpty},
		{ID: "B5", Name: "", Type: HexEmpty},
		{ID: "C8", Name: "", Type: HexEmpty},
		{ID: "E8", Name: "", Type: HexEmpty},
		{ID: "I8", Name: "", Type: HexEmpty},
		{ID: "D9", Name: "", Type: HexEmpty},
		{ID: "I10", Name: "", Type: HexEmpty},
	}
}

// parseHexCoord converts a hex ID like "K8" to (col, row) where col is 0-based
// from the letter (A=0, B=1, ...) and row is the number.
func parseHexCoord(id string) (col int, row int) {
	col = int(id[0] - 'A')
	for _, c := range id[1:] {
		row = row*10 + int(c-'0')
	}
	return
}

// hexCoordToID converts (col, row) back to a hex ID string.
func hexCoordToID(col, row int) string {
	return string(rune('A'+col)) + fmt.Sprintf("%d", row)
}

// hexNeighbor returns the (col, row) of the neighbor on the given edge.
// Edge convention for flat-top hexes: 0=N, 1=NE, 2=SE, 3=S, 4=SW, 5=NW.
//
// 18xx coordinate system: letter=column, number=row.
// Even columns have even rows, odd columns have odd rows.
// Neighbors:
//
//	0 (N):  same col, row-2
//	1 (NE): col+1,    row-1
//	2 (SE): col+1,    row+1
//	3 (S):  same col, row+2
//	4 (SW): col-1,    row+1
//	5 (NW): col-1,    row-1
func hexNeighbor(col, row, edge int) (int, int) {
	switch edge {
	case 0:
		return col, row - 2
	case 1:
		return col + 1, row - 1
	case 2:
		return col + 1, row + 1
	case 3:
		return col, row + 2
	case 4:
		return col - 1, row + 1
	case 5:
		return col - 1, row - 1
	default:
		return col, row
	}
}

// Default1889Adjacency computes the hex adjacency graph for 1889's flat-top hex grid.
// Edge convention for flat-top hexes:
//
//	0 = N, 1 = NE, 2 = SE, 3 = S, 4 = SW, 5 = NW
//
// Adjacency is computed algorithmically from hex coordinates.
func Default1889Adjacency() HexAdjacency {
	hexes := Default1889Map()

	// Build set of valid hex IDs.
	hexSet := make(map[string]bool, len(hexes))
	for _, h := range hexes {
		hexSet[h.ID] = true
	}

	adj := make(HexAdjacency, len(hexes))

	for _, h := range hexes {
		col, row := parseHexCoord(h.ID)
		var neighbors [6]string
		for edge := 0; edge < 6; edge++ {
			nc, nr := hexNeighbor(col, row, edge)
			if nc < 0 {
				continue
			}
			nID := hexCoordToID(nc, nr)
			if hexSet[nID] {
				neighbors[edge] = nID
			}
		}
		adj[h.ID] = neighbors
	}

	return adj
}
