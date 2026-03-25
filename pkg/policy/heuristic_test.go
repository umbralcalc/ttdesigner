package policy

import (
	"testing"

	"github.com/umbralcalc/stochadex/pkg/simulator"
	"github.com/umbralcalc/ttdesigner/pkg/engine"
)

func TestHeuristicAgent(t *testing.T) {
	t.Run("runs_50_steps_without_panic", func(t *testing.T) {
		builder := engine.NewGameBuilder(4, &HeuristicAgent{})
		settings, implementations := builder.Build()
		implementations.TerminationCondition = &simulator.NumberOfStepsTerminationCondition{
			MaxNumberOfSteps: 50,
		}

		coordinator := simulator.NewPartitionCoordinator(settings, implementations)
		coordinator.Run()

		// Verify the simulation ran and the turn state is valid.
		turnState := coordinator.Shared.StateHistories[0].Values.RawRowView(0)
		roundType := turnState[engine.TurnRoundType]
		if roundType < 0 || roundType > 2 {
			t.Errorf("invalid round type: %v", roundType)
		}
	})

	t.Run("harness_4_player", func(t *testing.T) {
		builder := engine.NewGameBuilder(4, &HeuristicAgent{})
		settings, implementations := builder.Build()
		implementations.OutputCondition = &simulator.EveryStepOutputCondition{}
		implementations.TerminationCondition = &simulator.NumberOfStepsTerminationCondition{
			MaxNumberOfSteps: 50,
		}

		if err := simulator.RunWithHarnesses(settings, implementations); err != nil {
			t.Errorf("harness failed: %v", err)
		}
	})

	t.Run("full_game_to_bank_break", func(t *testing.T) {
		builder := engine.NewGameBuilder(4, &HeuristicAgent{})
		settings, implementations := builder.Build()
		layout := builder.Layout()

		// Use bank-broken OR max steps, whichever comes first.
		implementations.TerminationCondition = &engine.OrTerminationCondition{
			Conditions: []simulator.TerminationCondition{
				&engine.BankBrokenTerminationCondition{BankPartitionIndex: layout.BankPartition},
				&simulator.NumberOfStepsTerminationCondition{MaxNumberOfSteps: 5000},
			},
		}

		coordinator := simulator.NewPartitionCoordinator(settings, implementations)
		coordinator.Run()

		bankState := coordinator.Shared.StateHistories[layout.BankPartition].Values.RawRowView(0)
		steps := coordinator.Shared.TimestepsHistory.CurrentStepNumber

		t.Logf("game ended after %d steps, bank cash: %.0f", steps, bankState[engine.BankCash])

		// Game should terminate in a reasonable number of steps.
		if steps >= 5000 {
			t.Errorf("game did not terminate within 5000 steps (likely stuck)")
		}
		if steps < 50 {
			t.Errorf("game ended suspiciously fast (%d steps)", steps)
		}
	})

	t.Run("pars_a_company_in_sr", func(t *testing.T) {
		builder := engine.NewGameBuilder(4, &HeuristicAgent{})
		settings, implementations := builder.Build()
		implementations.TerminationCondition = &simulator.NumberOfStepsTerminationCondition{
			MaxNumberOfSteps: 30,
		}

		coordinator := simulator.NewPartitionCoordinator(settings, implementations)
		coordinator.Run()

		layout := builder.Layout()

		// After 30 steps with heuristic agent, at least one company should be parred.
		// Check market partition: at least one company should have row >= 0.
		mktState := coordinator.Shared.StateHistories[layout.MarketPartition].Values.RawRowView(0)
		parred := false
		for i := 0; i < len(builder.Config.Companies); i++ {
			if mktState[engine.MarketRowIdx(i)] >= 0 {
				parred = true
				break
			}
		}
		if !parred {
			t.Error("expected at least one company to be parred after 30 steps")
		}
	})
}
