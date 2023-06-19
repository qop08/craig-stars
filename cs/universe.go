package cs

import (
	"sort"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type Universe struct {
	Planets              []*Planet                            `json:"planets,omitempty"`
	Fleets               []*Fleet                             `json:"fleets,omitempty"`
	Starbases            []*Fleet                             `json:"starbases,omitempty"`
	Wormholes            []*Wormhole                          `json:"wormholes,omitempty"`
	MineralPackets       []*MineralPacket                     `json:"mineralPackets,omitempty"`
	MineFields           []*MineField                         `json:"mineFields,omitempty"`
	Salvages             []*Salvage                           `json:"salvage,omitempty"`
	rules                *Rules                               `json:"-"`
	mapObjectsByPosition map[Vector][]interface{}             `json:"-"`
	fleetsByPosition     map[Vector]*Fleet                    `json:"-"`
	fleetsByNum          map[playerFleetNum]*Fleet            `json:"-"`
	designsByUUID        map[uuid.UUID]*ShipDesign            `json:"-"`
	battlePlansByName    map[playerBattlePlanName]*BattlePlan `json:"-"`
	salvagesByNum        map[int]*Salvage                     `json:"-"`
	wormholesByNum       map[int]*Wormhole                    `json:"-"`
}

func NewUniverse(rules *Rules) Universe {
	return Universe{
		rules:                rules,
		mapObjectsByPosition: make(map[Vector][]interface{}),
		fleetsByPosition:     make(map[Vector]*Fleet),
		fleetsByNum:          make(map[playerFleetNum]*Fleet),
		designsByUUID:        make(map[uuid.UUID]*ShipDesign),
		battlePlansByName:    make(map[playerBattlePlanName]*BattlePlan),
		salvagesByNum:        make(map[int]*Salvage),
		wormholesByNum:       make(map[int]*Wormhole),
	}
}

type mapObjectGetter interface {
	getShipDesign(uuid uuid.UUID) *ShipDesign
	getPlanet(num int) *Planet
	getFleet(playerNum int, num int) *Fleet
	getWormhole(num int) *Wormhole
	getSalvage(num int) *Salvage
	getCargoHolder(mapObjectType MapObjectType, num int, playerNum int) cargoHolder
	getMapObjectsAtPosition(position Vector) []interface{}
	isPositionValid(pos Vector, occupiedLocations *[]Vector, minDistance float64) bool
	updateMapObjectAtPosition(mo interface{}, originalPosition, newPosition Vector)
}

type playerFleetNum struct {
	PlayerNum int
	Num       int
}

type playerBattlePlanName struct {
	PlayerNum int
	Name      string
}

// build the maps used for the Get functions
func (u *Universe) buildMaps(players []*Player) {

	// make a big map to hold all of our universe objects by position
	u.mapObjectsByPosition = make(map[Vector][]interface{}, len(u.Planets))

	// build a map of designs by uuid
	// so we can inject the design into each token
	numDesigns := 0
	numBattlePlans := 0
	for _, p := range players {
		numDesigns += len(p.Designs)
		numBattlePlans += len(p.BattlePlans)
	}
	u.designsByUUID = make(map[uuid.UUID]*ShipDesign, numDesigns)
	u.battlePlansByName = make(map[playerBattlePlanName]*BattlePlan, numBattlePlans)

	for _, p := range players {
		for i := range p.Designs {
			design := &p.Designs[i]
			u.designsByUUID[design.UUID] = design
		}

		for i := range p.BattlePlans {
			plan := &p.BattlePlans[i]
			u.battlePlansByName[playerBattlePlanName{PlayerNum: p.Num, Name: plan.Name}] = plan
		}
	}

	u.fleetsByPosition = make(map[Vector]*Fleet, len(u.Fleets))
	u.fleetsByNum = make(map[playerFleetNum]*Fleet, len(u.Fleets))
	for _, fleet := range u.Fleets {
		u.addMapObjectByPosition(fleet, fleet.Position)
		u.fleetsByPosition[fleet.Position] = fleet
		u.fleetsByNum[playerFleetNum{fleet.PlayerNum, fleet.Num}] = fleet

		fleet.battlePlan = u.battlePlansByName[playerBattlePlanName{fleet.PlayerNum, fleet.BattlePlanName}]

		fleet.InjectDesigns(u.designsByUUID)
	}

	for _, starbase := range u.Starbases {
		u.Planets[starbase.PlanetNum-1].starbase = starbase
	}

	for _, planet := range u.Planets {
		u.addMapObjectByPosition(planet, planet.Position)
	}
	for _, wormhole := range u.Wormholes {
		u.addMapObjectByPosition(wormhole, wormhole.Position)
	}
	for _, mineralPacket := range u.MineralPackets {
		u.addMapObjectByPosition(mineralPacket, mineralPacket.Position)
	}
	for _, mineField := range u.MineFields {
		u.addMapObjectByPosition(mineField, mineField.Position)
	}

	u.salvagesByNum = make(map[int]*Salvage, len(u.Salvages))
	for _, salvage := range u.Salvages {
		u.salvagesByNum[salvage.Num] = salvage
		u.addMapObjectByPosition(salvage, salvage.Position)
	}

	u.wormholesByNum = make(map[int]*Wormhole, len(u.Wormholes))
	for _, wormhole := range u.Wormholes {
		u.wormholesByNum[wormhole.Num] = wormhole
		u.addMapObjectByPosition(wormhole, wormhole.Position)
	}

}

func (u *Universe) addMapObjectByPosition(mo interface{}, position Vector) {
	mos, found := u.mapObjectsByPosition[position]
	if !found {
		mos = []interface{}{}
		u.mapObjectsByPosition[position] = mos
	}
	mos = append(mos, mo)
	u.mapObjectsByPosition[position] = mos
}

// Check if a position vector is a mininum distance away from all other points
func (u *Universe) isPositionValid(pos Vector, occupiedLocations *[]Vector, minDistance float64) bool {
	minDistanceSquared := minDistance * minDistance

	for _, to := range *occupiedLocations {
		if pos.DistanceSquaredTo(to) <= minDistanceSquared {
			return false
		}
	}
	return true
}

// get all commandable map objects for a player
func (u *Universe) GetPlayerMapObjects(playerNum int) PlayerMapObjects {
	pmo := PlayerMapObjects{}

	pmo.Fleets = u.getFleets(playerNum)
	pmo.Planets = u.getPlanets(playerNum)
	pmo.MineFields = u.getMineFields(playerNum)

	return pmo
}

// get a ship design by uuid
func (u *Universe) getShipDesign(uuid uuid.UUID) *ShipDesign {
	return u.designsByUUID[uuid]
}

// Get a planet by num
func (u *Universe) getPlanet(num int) *Planet {
	return u.Planets[num-1]
}

// Get a fleet by player num and fleet num
func (u *Universe) getFleet(playerNum int, num int) *Fleet {
	return u.fleetsByNum[playerFleetNum{playerNum, num}]
}

// Get a planet by num
func (u *Universe) getWormhole(num int) *Wormhole {
	return u.wormholesByNum[num]
}

// Get a salvage by num
func (u *Universe) getSalvage(num int) *Salvage {
	return u.Salvages[num]
}

// Get a mineralpacket by num
func (u *Universe) getMineralPacket(num int) *MineralPacket {
	return u.MineralPackets[num]
}

// get a cargo holder by natural key (num, playerNum, etc)
func (u *Universe) getCargoHolder(mapObjectType MapObjectType, num int, playerNum int) cargoHolder {
	switch mapObjectType {
	case MapObjectTypePlanet:
		return u.getPlanet(num)
	case MapObjectTypeFleet:
		return u.getFleet(playerNum, num)
	}
	return nil
}

// mark a fleet as deleted and remove it from the universe
func (u *Universe) deleteFleet(fleet *Fleet) {
	fleet.Delete = true
	
	index := slices.Index(u.Fleets, fleet)
	slices.Delete(u.Fleets, index, index)

	delete(u.fleetsByNum, playerFleetNum{fleet.PlayerNum, fleet.Num})
	delete(u.fleetsByPosition, fleet.Position)
	u.removeMapObjectAtPosition(fleet, fleet.Position)
}

// move a fleet from one position to another
func (u *Universe) moveFleet(fleet *Fleet, originalPosition Vector) {
	fleet.Dirty = true
	delete(u.fleetsByPosition, originalPosition)
	u.fleetsByPosition[originalPosition] = fleet

	// upadte mapobjects position
	u.updateMapObjectAtPosition(fleet, originalPosition, fleet.Position)
}

// move a wormhole from one position to another
func (u *Universe) moveWormhole(wormhole *Wormhole, originalPosition Vector) {
	wormhole.Dirty = true
	u.updateMapObjectAtPosition(wormhole, originalPosition, wormhole.Position)
}

// delete a wormhole from the universe
func (u *Universe) deleteWormhole(wormhole *Wormhole) {
	wormhole.Delete = true

	index := slices.Index(u.Wormholes, wormhole)
	slices.Delete(u.Wormholes, index, index)

	delete(u.wormholesByNum, wormhole.Num)
	u.removeMapObjectAtPosition(wormhole, wormhole.Position)
}

// create a new wormhole in the universe
func (u *Universe) createWormhole(position Vector, stability WormholeStability, companion *Wormhole) *Wormhole {
	num := 1
	if len(u.Wormholes) > 0 {
		num = u.Wormholes[len(u.Wormholes)-1].Num + 1
	}

	wormhole := newWormhole(position, num, stability)

	if companion != nil {
		companion.DestinationNum = wormhole.Num
		wormhole.DestinationNum = companion.Num
	}

	// compute the spec for this wormhole
	wormhole.Spec = computeWormholeSpec(wormhole, u.rules)

	u.Wormholes = append(u.Wormholes, wormhole)
	u.wormholesByNum[wormhole.Num] = wormhole
	u.addMapObjectByPosition(wormhole, wormhole.Position)

	return wormhole
}

// delete a wormhole from the universe
func (u *Universe) createSalvage(position Vector, playerNum int, cargo Cargo) *Salvage {
	num := 1
	if len(u.Salvages) > 0 {
		num = u.Salvages[len(u.Salvages)-1].Num + 1
	}
	salvage := newSalvage(position, num, playerNum, cargo)
	u.Salvages = append(u.Salvages, salvage)
	u.salvagesByNum[num] = salvage
	u.addMapObjectByPosition(salvage, salvage.Position)

	return salvage
}

// delete a salvage from the universe
func (u *Universe) deleteSalvage(salvage *Salvage) {
	salvage.Delete = true

	index := slices.Index(u.Salvages, salvage)
	slices.Delete(u.Salvages, index, index)

	delete(u.salvagesByNum, salvage.Num)
	u.removeMapObjectAtPosition(salvage, salvage.Position)
}

func (u *Universe) getPlanets(playerNum int) []*Planet {
	planets := []*Planet{}
	for _, planet := range u.Planets {
		if planet.PlayerNum == playerNum {
			planets = append(planets, planet)
		}
	}
	return planets
}

func (u *Universe) getFleets(playerNum int) []*Fleet {
	fleets := []*Fleet{}
	for _, fleet := range u.Fleets {
		if fleet.PlayerNum == playerNum {
			fleets = append(fleets, fleet)
		}
	}
	return fleets
}

func (u *Universe) getMineFields(playerNum int) []*MineField {
	mineFields := []*MineField{}
	for _, mineField := range u.MineFields {
		if mineField.PlayerNum == playerNum {
			mineFields = append(mineFields, mineField)
		}
	}
	return mineFields
}

func (u *Universe) getNextFleetNum(playerNum int) int {
	num := 1

	playerFleets := u.getFleets(playerNum)
	orderedFleets := make([]*Fleet, len(playerFleets))
	copy(orderedFleets, playerFleets)
	sort.Slice(orderedFleets, func(i, j int) bool { return orderedFleets[i].Num < orderedFleets[j].Num })

	for i := 0; i < len(orderedFleets); i++ {
		fleet := orderedFleets[i]
		if i > 0 {
			// if we are past fleet #1 and we skipped a number, used the skipped number
			if fleet.Num > 1 && fleet.Num != orderedFleets[i-1].Num+1 {
				return orderedFleets[i-1].Num + 1
			}
		}
		// we are the next num...
		num = fleet.Num + 1
	}

	return num
}

// get a slice of mapobjects at a position, or nil if none
func (u *Universe) getMapObjectsAtPosition(position Vector) []interface{} {
	return u.mapObjectsByPosition[position]
}

// get a slice of mapobjects at a position, or nil if none
func (u *Universe) updateMapObjectAtPosition(mo interface{}, originalPosition, newPosition Vector) {
	mos := u.mapObjectsByPosition[originalPosition]
	if mos != nil {
		index := slices.IndexFunc(mos, func(item interface{}) bool { return item == mo })
		if index >= 0 && index < len(mos) {
			slices.Delete(mos, index, index)
		} else {
			log.Warn().Msgf("tried to update position of %s from %v to %v, but index %d of original position out of range", mo, originalPosition, newPosition, index)
		}
	} else {
		log.Warn().Msgf("tried to update position of %s from %v to %v, no mapobjects were found at %v", mo, originalPosition, newPosition, originalPosition)
	}

	// add the new object to the list
	u.addMapObjectByPosition(mo, newPosition)
}

// get a slice of mapobjects at a position, or nil if none
func (u *Universe) removeMapObjectAtPosition(mo interface{}, position Vector) {
	mos := u.mapObjectsByPosition[position]
	if mos != nil {
		index := slices.IndexFunc(mos, func(item interface{}) bool { return item == mo })
		if index >= 0 && index < len(mos) {
			slices.Delete(mos, index, index)
		} else {
			log.Warn().Msgf("tried to remove mapobject %s at position %v but index %d of position out of range", mo, position, index)
		}
	} else {
		log.Warn().Msgf("tried to to remove mapobject %s at position %v, no mapobjects were found at %v", mo, position, position)
	}
}
