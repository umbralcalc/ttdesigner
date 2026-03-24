package gamedata

import (
	"testing"

	"slices"
)

func TestDefault1889Config(t *testing.T) {
	cfg := Default1889Config()

	t.Run("config_validates", func(t *testing.T) {
		if err := cfg.Validate(); err != nil {
			t.Fatalf("default config failed validation: %v", err)
		}
	})

	t.Run("company_count", func(t *testing.T) {
		if len(cfg.Companies) != 7 {
			t.Errorf("expected 7 companies, got %d", len(cfg.Companies))
		}
	})

	t.Run("private_count", func(t *testing.T) {
		if len(cfg.Privates) != 7 {
			t.Errorf("expected 7 privates, got %d", len(cfg.Privates))
		}
	})

	t.Run("train_types", func(t *testing.T) {
		if len(cfg.Trains) != 6 {
			t.Errorf("expected 6 train types, got %d", len(cfg.Trains))
		}
	})

	t.Run("phase_count", func(t *testing.T) {
		if len(cfg.Phases) != 6 {
			t.Errorf("expected 6 phases, got %d", len(cfg.Phases))
		}
	})

	t.Run("bank_size", func(t *testing.T) {
		if cfg.BankSize != 7000 {
			t.Errorf("expected bank size 7000, got %d", cfg.BankSize)
		}
	})

	t.Run("starting_cash", func(t *testing.T) {
		if cfg.StartingCash[4] != 420 {
			t.Errorf("expected 420 starting cash for 4 players, got %d", cfg.StartingCash[4])
		}
		if cfg.StartingCash[6] != 390 {
			t.Errorf("expected 390 starting cash for 6 players, got %d", cfg.StartingCash[6])
		}
	})

	t.Run("cert_limits", func(t *testing.T) {
		if cfg.CertLimits[4] != 14 {
			t.Errorf("expected cert limit 14 for 4 players, got %d", cfg.CertLimits[4])
		}
	})
}

func TestDefault1889Companies(t *testing.T) {
	companies := Default1889Companies()

	t.Run("company_homes", func(t *testing.T) {
		homes := map[string]string{
			"AR": "K8", "IR": "E2", "SR": "I2", "KO": "K4",
			"TR": "F9", "KU": "C10", "UR": "B7",
		}
		for _, c := range companies {
			expected, ok := homes[c.Sym]
			if !ok {
				t.Errorf("unexpected company symbol %q", c.Sym)
				continue
			}
			if c.HomeHex != expected {
				t.Errorf("company %s home: expected %s, got %s", c.Sym, expected, c.HomeHex)
			}
		}
	})

	t.Run("token_costs", func(t *testing.T) {
		// TR and UR have 3 tokens (0, 40, 40); KU has 1 (0); rest have 2 (0, 40)
		for _, c := range companies {
			switch c.Sym {
			case "TR", "UR":
				if len(c.TokenCosts) != 3 {
					t.Errorf("%s: expected 3 tokens, got %d", c.Sym, len(c.TokenCosts))
				}
			case "KU":
				if len(c.TokenCosts) != 1 {
					t.Errorf("%s: expected 1 token, got %d", c.Sym, len(c.TokenCosts))
				}
			default:
				if len(c.TokenCosts) != 2 {
					t.Errorf("%s: expected 2 tokens, got %d", c.Sym, len(c.TokenCosts))
				}
			}
		}
	})
}

func TestDefault1889Privates(t *testing.T) {
	privates := Default1889Privates()

	t.Run("values_ascending", func(t *testing.T) {
		for i := 1; i < len(privates); i++ {
			if privates[i].Value < privates[i-1].Value {
				t.Errorf("privates not in ascending value order at index %d: %d < %d",
					i, privates[i].Value, privates[i-1].Value)
			}
		}
	})

	t.Run("min_players_filter", func(t *testing.T) {
		cfg := Default1889Config()

		twoPlayer := cfg.PrivatesForPlayerCount(2)
		fourPlayer := cfg.PrivatesForPlayerCount(4)
		sixPlayer := cfg.PrivatesForPlayerCount(6)

		if len(twoPlayer) != 5 {
			t.Errorf("expected 5 privates for 2 players, got %d", len(twoPlayer))
		}
		if len(fourPlayer) != 7 {
			t.Errorf("expected 7 privates for 4 players, got %d", len(fourPlayer))
		}
		if len(sixPlayer) != 7 {
			t.Errorf("expected 7 privates for 6 players, got %d", len(sixPlayer))
		}
	})
}

func TestDefault1889Market(t *testing.T) {
	market := Default1889Market()

	t.Run("dimensions", func(t *testing.T) {
		if market.Rows != 11 {
			t.Errorf("expected 11 rows, got %d", market.Rows)
		}
		if len(market.Cells[0]) != 15 {
			t.Errorf("expected 15 cols in row 0, got %d", len(market.Cells[0]))
		}
	})

	t.Run("par_values", func(t *testing.T) {
		pars := market.ParValues()
		expectedPrices := []int{100, 90, 80, 75, 70, 65}
		var gotPrices []int
		for _, p := range pars {
			gotPrices = append(gotPrices, p.Price)
		}
		if len(gotPrices) != len(expectedPrices) {
			t.Fatalf("expected %d par values, got %d", len(expectedPrices), len(gotPrices))
		}
		for _, expected := range expectedPrices {
			if !slices.Contains(gotPrices, expected) {
				t.Errorf("missing par value %d", expected)
			}
		}
	})

	t.Run("movement", func(t *testing.T) {
		// Start at par 100 (row 0, col 3)
		r, c := 0, 3
		if market.Price(r, c) != 100 {
			t.Fatalf("expected price 100 at (0,3), got %d", market.Price(r, c))
		}

		// Pay dividends → move right to 110
		r, c = market.MoveRight(r, c)
		if market.Price(r, c) != 110 {
			t.Errorf("after move right from 100: expected 110, got %d", market.Price(r, c))
		}

		// Withhold → move left back to 100
		r, c = market.MoveLeft(r, c)
		if market.Price(r, c) != 100 {
			t.Errorf("after move left from 110: expected 100, got %d", market.Price(r, c))
		}

		// Sell shares → move down to row 1 col 3 = 90
		r, c = market.MoveDown(r, c)
		if market.Price(r, c) != 90 {
			t.Errorf("after move down from (0,3): expected 90, got %d", market.Price(r, c))
		}
	})

	t.Run("orange_zone_closes", func(t *testing.T) {
		// Bottom-left cells should be orange
		zone := market.Zone(10, 0)
		if zone != MarketZoneOrange {
			t.Errorf("expected orange zone at (10,0), got %v", zone)
		}
	})

	t.Run("clamp_right_edge", func(t *testing.T) {
		// Row 0 has 15 cells (indices 0-14); moving right from col 14 should stay at 14
		r, c := market.MoveRight(0, 14)
		if c != 14 || r != 0 {
			t.Errorf("expected clamped at (0,14), got (%d,%d)", r, c)
		}
	})

	t.Run("clamp_bottom_row_short", func(t *testing.T) {
		// Row 10 has 4 cells; moving down from row 9 col 3 should clamp
		r, c := market.MoveDown(9, 3)
		if r != 10 {
			t.Errorf("expected row 10, got %d", r)
		}
		if c != 3 {
			t.Errorf("expected col 3, got %d", c)
		}
	})
}

func TestDefault1889Tiles(t *testing.T) {
	tiles := Default1889Tiles()
	manifest := Default1889TileManifest()

	t.Run("manifest_tiles_exist", func(t *testing.T) {
		for _, m := range manifest {
			if _, ok := tiles[m.TileID]; !ok {
				t.Errorf("manifest references tile %d which has no definition", m.TileID)
			}
		}
	})

	t.Run("yellow_city_tiles_have_upgrades", func(t *testing.T) {
		for id, tile := range tiles {
			if tile.Color == TileColorYellow && tile.Stop == StopCity && len(tile.UpgradesTo) == 0 {
				t.Errorf("yellow city tile %d has no upgrades", id)
			}
		}
	})

	t.Run("city_tiles_have_slots", func(t *testing.T) {
		for id, tile := range tiles {
			if tile.Stop == StopCity && tile.Slots == 0 {
				t.Errorf("city tile %d has 0 slots", id)
			}
		}
	})

	t.Run("city_tiles_have_revenue", func(t *testing.T) {
		for id, tile := range tiles {
			if tile.Stop == StopCity && tile.Revenue == 0 {
				t.Errorf("city tile %d has 0 revenue", id)
			}
		}
	})
}

func TestDefault1889Map(t *testing.T) {
	hexes := Default1889Map()

	t.Run("hex_count", func(t *testing.T) {
		// Count total hexes in the map
		if len(hexes) < 40 {
			t.Errorf("expected at least 40 hexes, got %d", len(hexes))
		}
	})

	t.Run("off_board_count", func(t *testing.T) {
		count := 0
		for _, h := range hexes {
			if h.Type == HexOffBoard {
				count++
			}
		}
		if count != 3 {
			t.Errorf("expected 3 off-board hexes, got %d", count)
		}
	})

	t.Run("gray_count", func(t *testing.T) {
		count := 0
		for _, h := range hexes {
			if h.Type == HexGray {
				count++
			}
		}
		if count != 4 {
			t.Errorf("expected 4 gray hexes, got %d", count)
		}
	})

	t.Run("company_homes_exist", func(t *testing.T) {
		hexMap := make(map[string]*HexDef)
		for i := range hexes {
			hexMap[hexes[i].ID] = &hexes[i]
		}
		for _, c := range Default1889Companies() {
			h, ok := hexMap[c.HomeHex]
			if !ok {
				t.Errorf("company %s home hex %s not found in map", c.Sym, c.HomeHex)
				continue
			}
			// Home hex should be a city (or gray with slots for Uwajima)
			if h.Type != HexCity && h.Type != HexGray {
				t.Errorf("company %s home hex %s is type %d, expected City or Gray", c.Sym, c.HomeHex, h.Type)
			}
		}
	})
}

func TestDefault1889Adjacency(t *testing.T) {
	adj := Default1889Adjacency()

	t.Run("bidirectional", func(t *testing.T) {
		opposite := func(edge int) int { return (edge + 3) % 6 }
		for hex, neighbors := range adj {
			for edge, neighbor := range neighbors {
				if neighbor == "" {
					continue
				}
				nAdj, ok := adj[neighbor]
				if !ok {
					t.Errorf("hex %s edge %d → %s, but %s has no adjacency entry", hex, edge, neighbor, neighbor)
					continue
				}
				opp := opposite(edge)
				if nAdj[opp] != hex {
					t.Errorf("hex %s edge %d → %s, but %s edge %d → %q (expected %s)",
						hex, edge, neighbor, neighbor, opp, nAdj[opp], hex)
				}
			}
		}
	})

	t.Run("key_connections", func(t *testing.T) {
		// Takamatsu (K4) should connect to Ritsurin Kouen (J5)
		k4 := adj["K4"]
		found := false
		for _, n := range k4 {
			if n == "J5" {
				found = true
				break
			}
		}
		if !found {
			t.Error("K4 (Takamatsu) should be adjacent to J5 (Ritsurin Kouen)")
		}
	})
}

func TestConfigValidation(t *testing.T) {
	t.Run("valid_config", func(t *testing.T) {
		cfg := Default1889Config()
		if err := cfg.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("zero_bank", func(t *testing.T) {
		cfg := Default1889Config()
		cfg.BankSize = 0
		if err := cfg.Validate(); err == nil {
			t.Error("expected error for zero bank size")
		}
	})

	t.Run("no_companies", func(t *testing.T) {
		cfg := Default1889Config()
		cfg.Companies = nil
		if err := cfg.Validate(); err == nil {
			t.Error("expected error for no companies")
		}
	})

	t.Run("duplicate_symbol", func(t *testing.T) {
		cfg := Default1889Config()
		cfg.Companies[1].Sym = cfg.Companies[0].Sym
		if err := cfg.Validate(); err == nil {
			t.Error("expected error for duplicate company symbol")
		}
	})
}
