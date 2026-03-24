package engine

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"github.com/umbralcalc/ttdesigner/pkg/gamedata"
)

// PlayerIteration tracks a single player's cash, shares, privates, and cert count.
//
// State layout: see PlayerCash..PlayerStateWidth constants in state.go.
//
// Receives action_values and turn_values via params_from_upstream.
type PlayerIteration struct {
	PlayerIndex int
	Config      *gamedata.GameConfig
}

func (p *PlayerIteration) Configure(partitionIndex int, settings *simulator.Settings) {}

func (p *PlayerIteration) Iterate(
	params *simulator.Params,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	prev := stateHistories[partitionIndex].Values.RawRowView(0)
	state := make([]float64, PlayerStateWidth)
	copy(state, prev)

	turnValues := params.Get("turn_values")
	actionValues := params.Get("action_values")
	actionType := actionValues[ActionType]

	activeType := turnValues[TurnActiveType]
	activeID := int(turnValues[TurnActiveID])

	// Only process if this player is the active entity.
	isActivePlayer := activeType == ActivePlayer && activeID == p.PlayerIndex
	if !isActivePlayer {
		// Still need to handle actions that affect non-active players
		// (e.g., dividends paying all shareholders). Handled in later steps.
		return state
	}

	switch actionType {
	case ActionBuyShare:
		p.handleBuyShare(state, actionValues)
	case ActionSellShares:
		p.handleSellShares(state, actionValues)
	case ActionParCompany:
		p.handleParCompany(state, actionValues)
	case ActionPass:
		state[PlayerPassed] = 1.0
	}

	return state
}

func (p *PlayerIteration) handleBuyShare(state []float64, action []float64) {
	companyID := int(action[ActionArg0])
	cost := action[ActionArg0+1]

	if state[PlayerCash] < cost {
		return
	}

	state[PlayerCash] -= cost
	state[PlayerShareIdx(companyID)] += 1
	state[PlayerCertCount] += 1
}

func (p *PlayerIteration) handleSellShares(state []float64, action []float64) {
	companyID := int(action[ActionArg0])
	numShares := action[ActionArg0+1]
	revenue := action[ActionArg0+2]

	state[PlayerShareIdx(companyID)] -= numShares
	state[PlayerCash] += revenue
	state[PlayerCertCount] -= numShares
}

func (p *PlayerIteration) handleParCompany(state []float64, action []float64) {
	companyID := int(action[ActionArg0])
	parPrice := action[ActionArg0+1]
	numShares := action[ActionArg0+2] // shares bought (usually 2 = 20%)

	cost := parPrice * numShares
	if state[PlayerCash] < cost {
		return
	}

	state[PlayerCash] -= cost
	state[PlayerShareIdx(companyID)] += numShares
	state[PlayerCertCount] += numShares
}

// InitPlayerState returns the initial state for a player partition.
func InitPlayerState(cfg *gamedata.GameConfig, numPlayers int) []float64 {
	state := make([]float64, PlayerStateWidth)
	state[PlayerCash] = float64(cfg.StartingCash[numPlayers])
	return state
}
