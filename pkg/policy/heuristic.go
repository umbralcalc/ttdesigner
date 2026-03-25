package policy

import (
	"github.com/umbralcalc/stochadex/pkg/simulator"
	"github.com/umbralcalc/ttdesigner/pkg/gamedata"
	"github.com/umbralcalc/ttdesigner/pkg/engine"
)

// HeuristicAgent implements a simple rule-based agent for 1889.
//
// Stock Round strategy:
//   - If cash > 2x cheapest par price and holds < 2 companies, par one.
//   - If a floated company has IPO shares and we can afford it, buy.
//   - Otherwise pass.
//
// Operating Round strategy:
//   - Tile lay: pick the first legal tile placement (if any).
//   - Token: pass (stub).
//   - Routes: pass/withhold (stub until route-finding is implemented).
//   - Buy train: pass (stub).
type HeuristicAgent struct{}

func (h *HeuristicAgent) ChooseAction(ctx *engine.GameContext) []float64 {
	turnState := ctx.TurnState
	roundType := turnState[engine.TurnRoundType]

	switch roundType {
	case engine.RoundStockRound:
		return h.chooseStockRoundAction(ctx)
	case engine.RoundOperatingRound:
		return h.chooseOperatingRoundAction(ctx)
	default:
		return passAction()
	}
}

func (h *HeuristicAgent) chooseStockRoundAction(ctx *engine.GameContext) []float64 {
	playerIndex := int(ctx.TurnState[engine.TurnActiveID])

	srCtx := engine.ExtractSRContext(
		playerIndex,
		ctx.StateHistories,
		ctx.Config,
		ctx.MarketGrid,
		ctx.NumPlayers,
		ctx.Layout,
	)

	actions := engine.LegalStockRoundActions(srCtx)
	if len(actions) <= 1 {
		return passAction()
	}

	// Priority 1: Par a company (cheapest par) if holding < 2 companies.
	var bestPar *engine.Action
	bestParCost := float64(999999)
	for i := range actions {
		if actions[i].Values[engine.ActionType] == engine.ActionParCompany {
			cost := actions[i].Values[engine.ActionArg0+1] * actions[i].Values[engine.ActionArg0+2]
			if cost < bestParCost {
				bestParCost = cost
				a := actions[i]
				bestPar = &a
			}
		}
	}

	companiesHeld := 0
	for _, shares := range srCtx.PlayerShares {
		if shares > 0 {
			companiesHeld++
		}
	}
	if bestPar != nil && companiesHeld < 2 {
		return bestPar.Values[:]
	}

	// Priority 2: Buy cheapest IPO share.
	var bestBuy *engine.Action
	bestBuyCost := float64(999999)
	for i := range actions {
		if actions[i].Values[engine.ActionType] == engine.ActionBuyShare {
			cost := actions[i].Values[engine.ActionArg0+1]
			if cost < bestBuyCost {
				bestBuyCost = cost
				a := actions[i]
				bestBuy = &a
			}
		}
	}
	if bestBuy != nil {
		return bestBuy.Values[:]
	}

	return passAction()
}

func (h *HeuristicAgent) chooseOperatingRoundAction(ctx *engine.GameContext) []float64 {
	orStep := ctx.TurnState[engine.TurnActionStep]
	companyIndex := int(ctx.TurnState[engine.TurnActiveID])

	// Check if company is floated; unfloated companies pass all steps.
	compState := ctx.StateHistories[ctx.Layout.CompanyPartitions[companyIndex]].Values.RawRowView(0)
	if compState[engine.CompFloated] == 0 {
		return passAction()
	}

	switch orStep {
	case engine.ORStepTileLay:
		return h.chooseTileLay(ctx, companyIndex)
	case engine.ORStepToken:
		return h.chooseToken(ctx, companyIndex)
	case engine.ORStepRoutes:
		return h.chooseRouteAction(ctx, companyIndex)
	case engine.ORStepBuyTrain:
		return h.chooseBuyTrain(ctx, companyIndex)
	default:
		return passAction()
	}
}

func (h *HeuristicAgent) chooseTileLay(ctx *engine.GameContext, companyIndex int) []float64 {
	tileCtx := extractTileLayContextFromGameCtx(ctx, companyIndex)
	actions := engine.LegalTileLayActions(tileCtx)

	if len(actions) == 0 {
		return passAction()
	}

	// Pick the first legal tile placement (simple heuristic).
	// Prefer placements near the company's home hex — but for now, just pick first.
	return actions[0].Values[:]
}

func (h *HeuristicAgent) chooseToken(ctx *engine.GameContext, companyIndex int) []float64 {
	tokenCtx := engine.ExtractTokenContext(
		companyIndex,
		ctx.StateHistories,
		ctx.Config,
		gamedata.Default1889Map(),
		ctx.Layout,
	)
	actions := engine.LegalTokenActions(tokenCtx)
	if len(actions) == 0 {
		return passAction()
	}
	// Place cheapest token available.
	return actions[0].Values[:]
}

func (h *HeuristicAgent) chooseRouteAction(ctx *engine.GameContext, companyIndex int) []float64 {
	compState := ctx.StateHistories[ctx.Layout.CompanyPartitions[companyIndex]].Values.RawRowView(0)
	mapState := ctx.StateHistories[ctx.Layout.MapPartition].Values.RawRowView(0)
	bankState := ctx.StateHistories[ctx.Layout.BankPartition].Values.RawRowView(0)
	gamePhase := int(bankState[engine.BankTrainPhase])

	hexes := gamedata.Default1889Map()
	tileDefs := gamedata.Default1889Tiles()
	adjacency := gamedata.Default1889Adjacency()

	graph := engine.BuildTrackGraph(mapState, hexes, tileDefs, adjacency)

	// Collect trains held by this company.
	var trains []int
	var distances []int
	for i, tr := range ctx.Config.Trains {
		count := int(compState[engine.CompTrainsBase+i])
		for j := 0; j < count; j++ {
			trains = append(trains, i)
			distances = append(distances, tr.Distance)
		}
	}

	if len(trains) == 0 {
		return passAction()
	}

	_, totalRevenue := engine.OptimalRouteAssignment(graph, companyIndex, trains, distances, gamePhase)

	if totalRevenue == 0 {
		return passAction()
	}

	// Heuristic: pay dividends if revenue > 2x cheapest needed train cost, else withhold.
	// Simple version: always pay dividends if we have revenue.
	action := make([]float64, engine.ActionStateWidth)
	needsTrain := len(trains) == 0
	treasury := compState[engine.CompTreasury]

	if needsTrain && treasury < float64(totalRevenue)*2 {
		// Withhold to build treasury.
		action[engine.ActionType] = engine.ActionWithhold
		action[engine.ActionArg0] = float64(companyIndex)
		action[engine.ActionArg0+1] = float64(totalRevenue)
	} else {
		// Pay dividends.
		action[engine.ActionType] = engine.ActionPayDividends
		action[engine.ActionArg0] = float64(companyIndex)
		action[engine.ActionArg0+1] = float64(totalRevenue)
	}

	return action
}

func (h *HeuristicAgent) chooseBuyTrain(ctx *engine.GameContext, companyIndex int) []float64 {
	compState := ctx.StateHistories[ctx.Layout.CompanyPartitions[companyIndex]].Values.RawRowView(0)
	bankState := ctx.StateHistories[ctx.Layout.BankPartition].Values.RawRowView(0)
	treasury := compState[engine.CompTreasury]

	// Count trains held.
	totalTrains := 0
	for i := range ctx.Config.Trains {
		totalTrains += int(compState[engine.CompTrainsBase+i])
	}

	// Check train limit for current phase.
	gamePhase := int(bankState[engine.BankTrainPhase])
	trainLimit := ctx.Config.Phases[gamePhase].TrainLimit

	// Must own at least one train. Buy cheapest available.
	// Also buy if under train limit (simple heuristic).
	if totalTrains >= trainLimit {
		return passAction()
	}
	if totalTrains > 0 {
		return passAction() // already has a train, don't buy more for now
	}

	for i, tr := range ctx.Config.Trains {
		avail := bankState[engine.BankTrainsBase+i]
		if avail <= 0 {
			continue
		}
		cost := float64(tr.Price)
		if treasury < cost {
			continue
		}
		action := make([]float64, engine.ActionStateWidth)
		action[engine.ActionType] = engine.ActionBuyTrain
		action[engine.ActionArg0] = float64(i)
		action[engine.ActionArg0+1] = cost
		return action
	}

	return passAction()
}

// extractTileLayContextFromGameCtx bridges GameContext → TileLayContext.
func extractTileLayContextFromGameCtx(ctx *engine.GameContext, companyIndex int) *engine.TileLayContext {
	return engine.ExtractTileLayContext(
		companyIndex,
		ctx.StateHistories,
		ctx.Config,
		gamedata.Default1889Map(),
		ctx.Layout,
	)
}

// Ensure HeuristicAgent satisfies the Agent interface.
var _ engine.Agent = (*HeuristicAgent)(nil)

// stateHistoryValue is a helper to read a value from a partition's state.
func stateHistoryValue(
	stateHistories []*simulator.StateHistory,
	partitionIdx int,
	stateIdx int,
) float64 {
	return stateHistories[partitionIdx].Values.At(0, stateIdx)
}

func passAction() []float64 {
	action := make([]float64, engine.ActionStateWidth)
	action[engine.ActionType] = engine.ActionPass
	return action
}
