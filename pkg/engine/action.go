package engine

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
)

// Agent is the interface that AI policies implement to choose actions.
type Agent interface {
	// ChooseAction selects an action given the current game state.
	// turnState is the turn partition's current state.
	// stateHistories gives access to all partition states.
	// Returns the action as a float64 slice of length ActionStateWidth.
	ChooseAction(
		turnState []float64,
		stateHistories []*simulator.StateHistory,
		timestepsHistory *simulator.CumulativeTimestepsHistory,
	) []float64
}

// PassAgent always passes. Used for skeleton testing.
type PassAgent struct{}

func (p *PassAgent) ChooseAction(
	turnState []float64,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	action := make([]float64, ActionStateWidth)
	action[ActionType] = ActionPass
	return action
}

// ActionIteration reads the current turn state and delegates to an Agent
// to choose the action for this step.
//
// It receives turn_values via params_from_upstream from the turn partition.
// Its output state is the chosen action vector.
type ActionIteration struct {
	Agent          Agent
	TurnPartition  int // partition index of the turn controller
}

func (a *ActionIteration) Configure(partitionIndex int, settings *simulator.Settings) {
	// No-op: Agent is set at construction time.
}

func (a *ActionIteration) Iterate(
	params *simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	// Read the turn state from the turn partition's current output.
	turnState := stateHistories[a.TurnPartition].Values.RawRowView(0)

	return a.Agent.ChooseAction(turnState, stateHistories, timestepsHistory)
}
