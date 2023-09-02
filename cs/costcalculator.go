package cs

import (
	"fmt"
)

type CostCalculator interface {
	StarbaseUpgradeCost(design, newDesign *ShipDesign) Cost
	CostOfOne(player *Player, item ProductionQueueItem) (Cost, error)
}

func NewCostCalculator() CostCalculator {
	return &costCalculate{}
}

type costCalculate struct {
}

// get the upgrade cost for a starbase
// TODO: do this
func (p *costCalculate) StarbaseUpgradeCost(design, newDesign *ShipDesign) Cost {
	return newDesign.Spec.Cost.Minus(design.Spec.Cost).MinZero()
}

// Get the cost of one item in a production queue, for a player
func (p *costCalculate) CostOfOne(player *Player, item ProductionQueueItem) (Cost, error) {
	cost := player.Race.Spec.Costs[item.Type]
	if item.Type == QueueItemTypeStarbase || item.Type == QueueItemTypeShipToken {
		if item.design != nil {
			cost = item.design.Spec.Cost
		} else {
			return Cost{}, fmt.Errorf("design %d not populated in queue item", item.DesignNum)
		}
	}
	return cost, nil
}
