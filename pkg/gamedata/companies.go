package gamedata

// CompanyDef defines a public corporation.
type CompanyDef struct {
	ID         int
	Sym        string
	Name       string
	HomeHex    string // hex coordinate, e.g. "K8"
	TokenCosts []int  // cost for each token beyond the free home token
	Color      string
	FloatPct   int // percentage of shares that must be sold to float (e.g. 50)
}

// PrivateDef defines a private company.
type PrivateDef struct {
	ID         int
	Sym        string
	Name       string
	Value      int
	Revenue    int
	Abilities  string // human-readable description
	BlocksHex  string // hex blocked while player-owned ("" if none)
	MinPlayers int    // 0 means available at all player counts
}

// TrainDef defines a train type.
type TrainDef struct {
	Name      string
	Distance  int // number of stops (-1 for unlimited/diesel)
	Price     int
	Quantity  int // -1 for unlimited
	RustsOn   string // train type that causes this to rust ("" if never)
	Discount  int    // trade-in discount when buying this train
	Events    []string
}

// Default1889Companies returns the 7 public corporations for 1889.
func Default1889Companies() []CompanyDef {
	return []CompanyDef{
		{ID: 0, Sym: "AR", Name: "Awa Railroad", HomeHex: "K8", TokenCosts: []int{0, 40}, Color: "#37383a", FloatPct: 50},
		{ID: 1, Sym: "IR", Name: "Iyo Railway", HomeHex: "E2", TokenCosts: []int{0, 40}, Color: "#f48221", FloatPct: 50},
		{ID: 2, Sym: "SR", Name: "Sanuki Railway", HomeHex: "I2", TokenCosts: []int{0, 40}, Color: "#76a042", FloatPct: 50},
		{ID: 3, Sym: "KO", Name: "Takamatsu & Kotohira Electric Railway", HomeHex: "K4", TokenCosts: []int{0, 40}, Color: "#d81e3e", FloatPct: 50},
		{ID: 4, Sym: "TR", Name: "Tosa Electric Railway", HomeHex: "F9", TokenCosts: []int{0, 40, 40}, Color: "#00a993", FloatPct: 50},
		{ID: 5, Sym: "KU", Name: "Tosa Kuroshio Railway", HomeHex: "C10", TokenCosts: []int{0}, Color: "#0189d1", FloatPct: 50},
		{ID: 6, Sym: "UR", Name: "Uwajima Railway", HomeHex: "B7", TokenCosts: []int{0, 40, 40}, Color: "#6f533e", FloatPct: 50},
	}
}

// Default1889Privates returns the 7 private companies for 1889.
func Default1889Privates() []PrivateDef {
	return []PrivateDef{
		{
			ID: 0, Sym: "TER", Name: "Takamatsu E-Railroad",
			Value: 20, Revenue: 5,
			Abilities: "Blocks Takamatsu (K4) while owned by a player.",
			BlocksHex: "K4",
		},
		{
			ID: 1, Sym: "MF", Name: "Mitsubishi Ferry",
			Value: 30, Revenue: 5,
			Abilities: "Player owner may place port tile (437) on a coastal town (B11, G10, I12, or J9) without a tile already, outside of another player's OR.",
		},
		{
			ID: 2, Sym: "ER", Name: "Ehime Railway",
			Value: 40, Revenue: 10,
			Abilities: "When sold to a corporation, the selling player may immediately place a green tile on Ohzu (C4). Blocks C4 while owned by a player.",
			BlocksHex: "C4",
		},
		{
			ID: 3, Sym: "SMR", Name: "Sumitomo Mines Railway",
			Value: 50, Revenue: 15,
			Abilities: "Owning corporation may ignore building cost for mountain hexes (discount 80).",
		},
		{
			ID: 4, Sym: "DR", Name: "Dougo Railway",
			Value: 60, Revenue: 15,
			Abilities: "Owning player may exchange for a 10% share of Iyo Railway (IR) from the IPO.",
		},
		{
			ID: 5, Sym: "SIR", Name: "South Iyo Railway",
			Value: 80, Revenue: 20,
			Abilities:  "",
			MinPlayers: 3,
		},
		{
			ID: 6, Sym: "UTF", Name: "Uno-Takamatsu Ferry",
			Value: 150, Revenue: 30,
			Abilities:  "Does not close while owned by a player. Revenue increases to 50 when first 5-train purchased.",
			MinPlayers: 4,
		},
	}
}

// Default1889Trains returns the train types for 1889.
func Default1889Trains() []TrainDef {
	return []TrainDef{
		{Name: "2", Distance: 2, Price: 80, Quantity: 6, RustsOn: "4"},
		{Name: "3", Distance: 3, Price: 180, Quantity: 5, RustsOn: "6"},
		{Name: "4", Distance: 4, Price: 300, Quantity: 4, RustsOn: "D"},
		{Name: "5", Distance: 5, Price: 450, Quantity: 3, Events: []string{"close_companies"}},
		{Name: "6", Distance: 6, Price: 630, Quantity: 2},
		{Name: "D", Distance: -1, Price: 1100, Quantity: -1, Discount: 300},
	}
}
