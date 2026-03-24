package gamedata

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// PhaseDef defines a game phase triggered by purchasing a train type.
type PhaseDef struct {
	Name           string   // e.g. "2", "3", "4", "5", "6", "D"
	TriggerTrain   string   // train type that triggers this phase ("" for starting phase)
	TilesAvailable []string // tile colors available: "yellow", "green", "brown"
	TrainLimit     int      // max trains per company
	ORsPerSR       int      // operating rounds per stock round
	CanBuyPrivates bool
}

// GameConfig holds all tweakable game parameters for a single 18xx title.
// Designed to be loaded from YAML for designer experimentation.
type GameConfig struct {
	Title       string `yaml:"title"`
	BankSize    int    `yaml:"bank_size"`
	CurrencyFmt string `yaml:"currency_fmt"` // e.g. "¥%d"

	// Player counts and per-count settings.
	MinPlayers    int            `yaml:"min_players"`
	MaxPlayers    int            `yaml:"max_players"`
	StartingCash  map[int]int    `yaml:"starting_cash"`   // players → cash
	CertLimits    map[int]int    `yaml:"cert_limits"`     // players → cert limit

	// Entity definitions.
	Companies []CompanyDef `yaml:"companies"`
	Privates  []PrivateDef `yaml:"privates"`
	Trains    []TrainDef   `yaml:"trains"`
	Phases    []PhaseDef   `yaml:"phases"`

	// Market is not YAML-loaded — use Default1889Market().
	// Tiles are not YAML-loaded — use Default1889Tiles().
	// Map is not YAML-loaded — use Default1889Map().

	// Rules
	MustSellInBlocks bool `yaml:"must_sell_in_blocks"`
}

// Default1889Config returns the full default configuration for 1889.
func Default1889Config() *GameConfig {
	return &GameConfig{
		Title:       "1889: History of Shikoku Railways",
		BankSize:    7000,
		CurrencyFmt: "¥%d",
		MinPlayers:  2,
		MaxPlayers:  6,
		StartingCash: map[int]int{
			2: 420, 3: 420, 4: 420, 5: 390, 6: 390,
		},
		CertLimits: map[int]int{
			2: 25, 3: 19, 4: 14, 5: 12, 6: 11,
		},
		Companies: Default1889Companies(),
		Privates:  Default1889Privates(),
		Trains:    Default1889Trains(),
		Phases:    Default1889Phases(),
		MustSellInBlocks: true,
	}
}

// Default1889Phases returns the phase definitions for 1889.
func Default1889Phases() []PhaseDef {
	return []PhaseDef{
		{Name: "2", TriggerTrain: "", TilesAvailable: []string{"yellow"}, TrainLimit: 4, ORsPerSR: 1, CanBuyPrivates: false},
		{Name: "3", TriggerTrain: "3", TilesAvailable: []string{"yellow", "green"}, TrainLimit: 4, ORsPerSR: 2, CanBuyPrivates: true},
		{Name: "4", TriggerTrain: "4", TilesAvailable: []string{"yellow", "green"}, TrainLimit: 3, ORsPerSR: 2, CanBuyPrivates: true},
		{Name: "5", TriggerTrain: "5", TilesAvailable: []string{"yellow", "green", "brown"}, TrainLimit: 2, ORsPerSR: 3, CanBuyPrivates: false},
		{Name: "6", TriggerTrain: "6", TilesAvailable: []string{"yellow", "green", "brown"}, TrainLimit: 2, ORsPerSR: 3, CanBuyPrivates: false},
		{Name: "D", TriggerTrain: "D", TilesAvailable: []string{"yellow", "green", "brown"}, TrainLimit: 2, ORsPerSR: 3, CanBuyPrivates: false},
	}
}

// LoadConfig reads a GameConfig from a YAML file.
func LoadConfig(path string) (*GameConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	var cfg GameConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}
	return &cfg, nil
}

// Validate checks that the config is internally consistent.
func (c *GameConfig) Validate() error {
	if c.BankSize <= 0 {
		return fmt.Errorf("bank_size must be positive, got %d", c.BankSize)
	}
	if c.MinPlayers < 2 || c.MaxPlayers > 6 || c.MinPlayers > c.MaxPlayers {
		return fmt.Errorf("invalid player range: %d-%d", c.MinPlayers, c.MaxPlayers)
	}
	for p := c.MinPlayers; p <= c.MaxPlayers; p++ {
		if _, ok := c.StartingCash[p]; !ok {
			return fmt.Errorf("missing starting_cash for %d players", p)
		}
		if _, ok := c.CertLimits[p]; !ok {
			return fmt.Errorf("missing cert_limits for %d players", p)
		}
	}
	if len(c.Companies) == 0 {
		return fmt.Errorf("no companies defined")
	}
	if len(c.Trains) == 0 {
		return fmt.Errorf("no trains defined")
	}
	if len(c.Phases) == 0 {
		return fmt.Errorf("no phases defined")
	}

	// Check company IDs are sequential.
	for i, co := range c.Companies {
		if co.ID != i {
			return fmt.Errorf("company %q has ID %d, expected %d", co.Sym, co.ID, i)
		}
	}

	// Check unique company symbols.
	syms := make(map[string]bool)
	for _, co := range c.Companies {
		if syms[co.Sym] {
			return fmt.Errorf("duplicate company symbol %q", co.Sym)
		}
		syms[co.Sym] = true
	}

	return nil
}

// PrivatesForPlayerCount returns privates available at the given player count.
func (c *GameConfig) PrivatesForPlayerCount(players int) []PrivateDef {
	var result []PrivateDef
	for _, p := range c.Privates {
		if p.MinPlayers == 0 || players >= p.MinPlayers {
			result = append(result, p)
		}
	}
	return result
}
