package policy

import (
	"math"
	"testing"

	"github.com/umbralcalc/stochadex/pkg/simulator"
	"github.com/umbralcalc/18xxdesigner/pkg/engine"
)

func TestMCTSAgent(t *testing.T) {
	t.Run("chooses_valid_action", func(t *testing.T) {
		// Run a game for 30 steps to reach SR, then use MCTS for one decision.
		builder := engine.NewGameBuilder(4, &HeuristicAgent{})
		settings, impls := builder.Build()
		impls.TerminationCondition = &simulator.NumberOfStepsTerminationCondition{
			MaxNumberOfSteps: 30,
		}
		coordinator := simulator.NewPartitionCoordinator(settings, impls)
		coordinator.Run()

		layout := builder.Layout()
		turnState := coordinator.Shared.StateHistories[layout.TurnPartition].Values.RawRowView(0)
		ctx := &engine.GameContext{
			TurnState:        turnState,
			StateHistories:   coordinator.Shared.StateHistories,
			TimestepsHistory: coordinator.Shared.TimestepsHistory,
			Layout:           layout,
			Config:           builder.Config,
			MarketGrid:       builder.Market,
			NumPlayers:       4,
		}

		mcts := NewMCTSAgent(0, 5)
		mcts.MaxPlayoutSteps = 1000

		action := mcts.ChooseAction(ctx)
		if len(action) != engine.ActionStateWidth {
			t.Errorf("expected action width %d, got %d",
				engine.ActionStateWidth, len(action))
		}
	})

	t.Run("full_game_completes", func(t *testing.T) {
		agent := NewMCTSAgent(0, 3)
		agent.MaxPlayoutSteps = 500

		builder := engine.NewGameBuilder(4, agent)
		settings, impls := builder.Build()
		layout := builder.Layout()

		impls.TerminationCondition = &engine.OrTerminationCondition{
			Conditions: []simulator.TerminationCondition{
				&engine.BankBrokenTerminationCondition{
					BankPartitionIndex: layout.BankPartition,
				},
				&simulator.NumberOfStepsTerminationCondition{
					MaxNumberOfSteps: 5000,
				},
			},
		}

		coordinator := simulator.NewPartitionCoordinator(settings, impls)
		coordinator.Run()

		steps := coordinator.Shared.TimestepsHistory.CurrentStepNumber
		t.Logf("game ended after %d steps", steps)
		if steps >= 5000 {
			t.Error("game did not terminate within 5000 steps")
		}
	})
}

func TestMCTSPlayoutPartitions(t *testing.T) {
	t.Run("runs_without_nan", func(t *testing.T) {
		// Get a mid-game snapshot by running heuristic for 30 steps.
		builder := engine.NewGameBuilder(4, &HeuristicAgent{})
		settings, impls := builder.Build()
		impls.TerminationCondition = &simulator.NumberOfStepsTerminationCondition{
			MaxNumberOfSteps: 30,
		}
		coordinator := simulator.NewPartitionCoordinator(settings, impls)
		coordinator.Run()

		snapshot := make([][]float64, len(coordinator.Shared.StateHistories))
		for i, sh := range coordinator.Shared.StateHistories {
			row := sh.Values.RawRowView(0)
			cp := make([]float64, len(row))
			copy(cp, row)
			snapshot[i] = cp
		}

		candidates := [][]float64{passAction(), passAction()}
		partitions := NewMCTSPlayoutPartitions(AppliedMCTSOptimisation{
			Selector: MCTSSelector{
				CandidateActions: candidates,
				ExplorationC:     1.414,
			},
			Playout: MCTSPlayout{
				Snapshot:        snapshot,
				NumPlayers:      4,
				MaxPlayoutSteps: 1000,
				PlayerIndex:     0,
			},
		})

		gen := simulator.NewConfigGenerator()
		gen.SetSimulation(&simulator.SimulationConfig{
			OutputCondition: &simulator.NilOutputCondition{},
			OutputFunction:  &simulator.NilOutputFunction{},
			TerminationCondition: &simulator.NumberOfStepsTerminationCondition{
				MaxNumberOfSteps: 3,
			},
			TimestepFunction: &simulator.ConstantTimestepFunction{Stepsize: 1.0},
			InitTimeValue:    0.0,
		})
		for _, p := range partitions {
			gen.SetPartition(p)
		}

		outerSettings, outerImpls := gen.GenerateConfigs()
		outerCoordinator := simulator.NewPartitionCoordinator(
			outerSettings, outerImpls)
		outerCoordinator.Run()

		// Verify no NaN in statistics partition (index 2).
		statsState := outerCoordinator.Shared.StateHistories[2].Values.RawRowView(0)
		for i, v := range statsState {
			if math.IsNaN(v) {
				t.Errorf("NaN at statistics state index %d", i)
			}
		}

		// Verify some playouts were recorded.
		totalVisits := statsState[0] + statsState[2]
		if totalVisits < 1 {
			t.Error("expected at least one playout visit")
		}
		t.Logf("visits: action0=%.0f action1=%.0f", statsState[0], statsState[2])
	})
}

func TestEnumerateLegalMoves(t *testing.T) {
	t.Run("sr_has_multiple_options", func(t *testing.T) {
		builder := engine.NewGameBuilder(4, &HeuristicAgent{})
		settings, impls := builder.Build()
		impls.TerminationCondition = &simulator.NumberOfStepsTerminationCondition{
			MaxNumberOfSteps: 20,
		}
		coordinator := simulator.NewPartitionCoordinator(settings, impls)
		coordinator.Run()

		layout := builder.Layout()
		turnState := coordinator.Shared.StateHistories[layout.TurnPartition].Values.RawRowView(0)

		// Skip if not in SR.
		if turnState[engine.TurnRoundType] != engine.RoundStockRound {
			t.Skip("not in stock round after 20 steps")
		}

		ctx := &engine.GameContext{
			TurnState:        turnState,
			StateHistories:   coordinator.Shared.StateHistories,
			TimestepsHistory: coordinator.Shared.TimestepsHistory,
			Layout:           layout,
			Config:           builder.Config,
			MarketGrid:       builder.Market,
			NumPlayers:       4,
		}

		moves := enumerateLegalMoves(ctx)
		if len(moves) == 0 {
			t.Error("expected at least one legal move")
		}
		t.Logf("found %d legal moves in SR", len(moves))
	})
}
