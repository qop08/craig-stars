package cs

import (
	"fmt"
	"math"
)

type MineFieldType string

const (
	MineFieldTypeStandard  MineFieldType = "Standard"
	MineFieldTypeHeavy     MineFieldType = "Heavy"
	MineFieldTypeSpeedBump MineFieldType = "SpeedBump"
)

func (t MineFieldType) String() string {
	switch t {
	case MineFieldTypeSpeedBump:
		return "Speed Bump"
	default:
		return string(t)
	}
}

func (t MineFieldType) CanDetonate() bool {
	switch t {
	case MineFieldTypeStandard:
		return true
	default:
		return false
	}
}

type MineField struct {
	MapObject
	MineFieldOrders
	MineFieldType MineFieldType `json:"mineFieldType"`
	NumMines      int           `json:"numMines"`
	Spec          MineFieldSpec `json:"spec"`
}

type MineFieldOrders struct {
	Detonate bool `json:"detonate,omitempty"`
}

type MineFieldSpec struct {
	Radius      float64 `json:"radius"`
	DecayRate   int     `json:"decayRate"`
	CanDetonate bool    `json:"canDetonate"`
}

type MineFieldStats struct {
	MinDamagePerFleetRS int     `json:"minDamagePerFleetRS"`
	DamagePerEngineRS   int     `json:"damagePerEngineRS"`
	MaxSpeed            int     `json:"maxSpeed"`
	ChanceOfHit         float64 `json:"chanceOfHit"`
	MinDamagePerFleet   int     `json:"minDamagePerFleet"`
	DamagePerEngine     int     `json:"damagePerEngine"`
	SweepFactor         float64 `json:"sweepFactor"`
	MinDecay            int     `json:"minDecay"`
	CanDetonate         bool    `json:"canDetonate"`
}

type MineFieldDamage struct {
	Damage         int  `json:"damage,omitempty"`
	ShipsDestroyed int  `json:"shipsDestroyed,omitempty"`
	FleetDestroyed bool `json:"fleetDestroyed,omitempty"`
}

// The radius of a minefield is the sqrt of its mines
func (mf *MineField) Radius() float64 {
	return math.Sqrt(float64(mf.NumMines))
}

func computeMinefieldSpec(rules *Rules, player *Player, mineField *MineField, numPlanets int) MineFieldSpec {
	spec := MineFieldSpec{}
	spec.Radius = mineField.Radius()
	spec.DecayRate = mineField.getDecayRate(rules, player, numPlanets)
	spec.CanDetonate = mineField.MineFieldType.CanDetonate()

	return spec
}

func newMineField(player *Player, mineFieldType MineFieldType, numMines int, num int, position Vector) *MineField {
	return &MineField{
		MapObject: MapObject{
			Type:      MapObjectTypeMineField,
			PlayerNum: player.Num,
			Num:       num,
			Name:      fmt.Sprintf("%s %s Mine Field #%d", player.Race.PluralName, mineFieldType.String(), num),
			Position:  position,
		},
		MineFieldType: mineFieldType,
		NumMines:      numMines,
	}
}

func (mineField *MineField) withOrders(orders MineFieldOrders) *MineField {
	mineField.MineFieldOrders = orders
	return mineField
}

// get the number of mines that will decay this year
// * The base rate for minefield decay is 2% per year.
// * Minefields will decay an additional 4% per planet that is within the field, or 1% per planet for SD races.
// * A detonating SD minefield has an additional 25% decay each year.
// * Normal and Heavy Minefields have a minimum total decay rate of 10 mines per year
// * Speed Bump Minefields have a minimum total decay rate of 2 mines per year
// * There is a maximum total decay rate of 50% per year.
func (mineField *MineField) getDecayRate(rules *Rules, player *Player, numPlanets int) int {
	if !mineField.Owned() {
		// we can't determine decay rate for minefields we don't own
		return -1
	}

	decayRate := player.Race.Spec.MineFieldBaseDecayRate
	decayRate += player.Race.Spec.MineFieldPlanetDecayRate * float64(numPlanets)
	if mineField.Detonate {
		decayRate += player.Race.Spec.MineFieldDetonateDecayRate
	}

	// Space Demolition mines decay slower
	decayFactor := player.Race.Spec.MineFieldMinDecayFactor
	decayRate *= decayFactor
	decayRate = math.Min(decayRate, player.Race.Spec.MineFieldMaxDecayRate)

	// we decay at least 10 mines a year for normal and standard mines
	decayedMines := MaxInt(rules.MineFieldStatsByType[mineField.MineFieldType].MinDecay, int(float64(mineField.NumMines)*decayRate+0.5))
	return decayedMines
}

// damage a fleet that hit this minefield
// https://wiki.starsautohost.org/wiki/Guts_of_Minefields
func (mineField *MineField) damageFleet(fleet *Fleet, fleetPlayer *Player, stats MineFieldStats) MineFieldDamage {
	hasRamScoop := false
	for _, token := range fleet.Tokens {
		if token.design.Spec.Engine.FreeSpeed > 1 {
			hasRamScoop = true
			break
		}
	}

	minDamage := stats.MinDamagePerFleet
	damagePerEngine := stats.DamagePerEngine
	if hasRamScoop {
		minDamage = stats.MinDamagePerFleetRS
		damagePerEngine = stats.DamagePerEngineRS
	}

	totalDamage := 0
	shipsDestroyed := 0

	if minDamage > 0 {
		if fleet.Spec.TotalShips <= 5 {
			firstDesignNumEngines := 0
			for i := range fleet.Tokens {
				token := &fleet.Tokens[i]
				if mineField.Detonate && token.design.Spec.ImmuneToOwnDetonation && mineField.OwnedBy(fleetPlayer.Num) {
					continue
				}

				design := token.design
				if firstDesignNumEngines == 0 {
					firstDesignNumEngines = design.Spec.NumEngines
					tokenDamage := firstDesignNumEngines * minDamage
					totalDamage += tokenDamage
					result := token.applyMineDamage(tokenDamage)
					shipsDestroyed += result.shipsDestroyed
				} else if design.Spec.NumEngines > firstDesignNumEngines {
					tokenDamage := damagePerEngine * (design.Spec.NumEngines - firstDesignNumEngines) * token.Quantity
					totalDamage += tokenDamage
					result := token.applyMineDamage(tokenDamage)
					shipsDestroyed += result.shipsDestroyed
				}
			}
		} else {
			for i := range fleet.Tokens {
				token := &fleet.Tokens[i]
				if mineField.Detonate && token.design.Spec.ImmuneToOwnDetonation && mineField.OwnedBy(fleetPlayer.Num) {
					continue
				}

				design := token.design
				tokenDamage := damagePerEngine * design.Spec.NumEngines * token.Quantity
				totalDamage += tokenDamage
				result := token.applyMineDamage(tokenDamage)
				shipsDestroyed += result.shipsDestroyed
			}
		}
	}

	return MineFieldDamage{
		Damage:         totalDamage,
		ShipsDestroyed: shipsDestroyed,
		FleetDestroyed: fleet.Spec.TotalShips <= shipsDestroyed,
	}
}

// When a minefield is collided with, reduce its number of mines
func (mineField *MineField) reduceMineFieldOnImpact() {
	numMines := mineField.NumMines
	if numMines <= 10 {
		numMines = 0
	} else if numMines <= 200 {
		numMines -= 10
	} else if numMines <= 1000 {
		numMines = int(float64(numMines) * 0.95)
	} else if numMines <= 5000 {
		numMines -= 50
	} else {
		numMines = int(float64(numMines) * 0.95)
	}
	mineField.NumMines = numMines
}

func (mineField *MineField) sweep(rules *Rules, fleetPosition Vector, mineSweep int) int {

	// we can only sweep up to our position in the minefield, so figure out how far we are from the center
	// and subtract that from the radius to determine the edge amount
	//		***
	// 	   *****
	// 	  *******
	// 	 *F**C**** // fleet is 1 from the edge, 3 from the center
	// 	  *******
	// 	   *****
	// 	    ***
	//
	radius := mineField.Radius()
	distFromCenter := fleetPosition.DistanceTo(mineField.Position)
	distFromEdge := mineField.Radius() - distFromCenter

	// radius of a minefield is sqrt(numMines) so we can sweet our dist^2 in mines
	sweepableMines := mineField.NumMines - int(math.Ceil((radius-distFromEdge)*(radius-distFromEdge)))

	old := mineField.NumMines
	mineField.NumMines -= MinInt(sweepableMines, int(float64(mineSweep)*rules.MineFieldStatsByType[mineField.MineFieldType].SweepFactor))
	mineField.NumMines = MaxInt(mineField.NumMines, 0)

	numSwept := old - mineField.NumMines
	return numSwept
}

// / Check for mine field collisions. If we collide with one, do damage and stop the fleet
func checkForMineFieldCollision(rules *Rules, playerGetter playerGetter, mapObjectGetter mapObjectGetter, fleet *Fleet, dest Waypoint, distance float64) (mineField *MineField, distanceTravelled float64) {
	distanceTravelled = distance
	fleetPlayer := playerGetter.getPlayer(fleet.PlayerNum)
	safeWarpBonus := fleetPlayer.Race.Spec.MineFieldSafeWarpBonus

	// see if we are colliding with any of these minefields
	for _, mineField := range mapObjectGetter.getAllMineFields() {
		// we don't hit our own minefields
		if mineField.PlayerNum == fleet.PlayerNum {
			continue
		}

		// our allies don't hit our minefields
		mineFieldPlayer := playerGetter.getPlayer(mineField.PlayerNum)
		if mineFieldPlayer.IsFriend(fleetPlayer.Num) {
			continue
		}

		// we only check if we are going faster than allowed by the minefield.
		stats := rules.MineFieldStatsByType[mineField.MineFieldType]
		if dest.WarpSpeed > stats.MaxSpeed+safeWarpBonus {
			// this is not our minefield, and we are going fast, check if we intersect.
			from := fleet.Position
			to := (dest.Position.Subtract(fleet.Position).Normalized()).Scale(distance).Add(from)
			collision := segmentIntersectsCircle(from, to, mineField.Position, mineField.Spec.Radius)
			if collision == -1 {
				// miss! phew, that was close!
				continue
			} else {
				// we are travelling through this minefield, for each light year we go through, check for a hit
				// collision is 0 to 1, which is the percent of our travel segment that is NOT in the field.
				// figure out what that is in lightYears
				// if we are travelling 32 light years and 3/4 of it is through the minefield, we need to check
				// for collision 24 times
				lightYearsInField := int(math.Min(float64(mineField.Spec.Radius), math.Ceil(float64((1-collision)*distance))))
				lightYearsBeforeField := collision * distance

				// Each type of minefield has a chance to hit based on how fast
				// the fleet is travelling through the field. A normal mine has a .3% chance
				// of hitting a ship per extra warp over warp 4, so a warp 9 ship
				// has a 1.5% chance of hitting a mine per lightyear travelled
				unsafeWarp := dest.WarpSpeed - (stats.MaxSpeed + safeWarpBonus)
				chanceToHit := stats.ChanceOfHit * float64(unsafeWarp)
				for checkNum := 0; checkNum < lightYearsInField; checkNum++ {
					if chanceToHit >= rules.random.Float64() {
						// ouch, we hit a minefield!
						// we stop moving at the hit, so if we made it 8 checks out of 24 for our above example
						// we only travel 8 lightyears through the field (plus whatever distance we travelled to get to the field)
						fleet.struckMineField = true
						distanceTravelled = lightYearsBeforeField + float64(checkNum)
						return mineField, distanceTravelled
					}
				}
			}
		}

	}

	return mineField, distance
}

// Move this minefield closer to us (in case it's not in our location)
// This was taken from the FreeStars codebase (like many other things)
func (mineField *MineField) moveTowardsMineLayer(position Vector, minesLaid int) {
	totalDist := position.DistanceTo(mineField.Position)

	moveTowardsFactor := math.Min(1, float64(minesLaid)/float64(mineField.NumMines))
	heading := position.Subtract(mineField.Position).Normalized()

	// move the minefield towards the fleet
	mineField.Position = mineField.Position.Add(heading.Normalized().Scale(totalDist * moveTowardsFactor)).Round()
}
