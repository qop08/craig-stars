package cs

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type PlayerMessage struct {
	Type            PlayerMessageType       `json:"type,omitempty"`
	Text            string                  `json:"text,omitempty"`
	TargetNum       int                     `json:"targetNum,omitempty"`
	TargetPlayerNum int                     `json:"targetPlayerNum,omitempty"`
	BattleNum       int                     `json:"battleNum,omitempty"`
	TargetType      PlayerMessageTargetType `json:"targetType,omitempty"`
}

type PlayerMessageTargetType string

const (
	TargetNone          PlayerMessageTargetType = ""
	TargetPlanet        PlayerMessageTargetType = "Planet"
	TargetFleet         PlayerMessageTargetType = "Fleet"
	TargetWormhole      PlayerMessageTargetType = "Wormhole"
	TargetMineField     PlayerMessageTargetType = "MineField"
	TargetMysteryTrader PlayerMessageTargetType = "MysteryTrader"
	TargetBattle        PlayerMessageTargetType = "Battle"
)

type PlayerMessageType string

const (
	PlayerMessageNone                               PlayerMessageType = ""
	PlayerMessageInfo                               PlayerMessageType = "Info"
	PlayerMessageError                              PlayerMessageType = "Error"
	PlayerMessageHomePlanet                         PlayerMessageType = "HomePlanet"
	PlayerMessagePlayerDiscovery                    PlayerMessageType = "PlayerDiscovery"
	PlayerMessagePlanetDiscovery                    PlayerMessageType = "PlanetDiscovery"
	PlayerMessagePlanetProductionQueueEmpty         PlayerMessageType = "PlanetProductionQueueEmpty"
	PlayerMessagePlanetProductionQueueComplete      PlayerMessageType = "PlanetProductionQueueComplete"
	PlayerMessageBuiltMine                          PlayerMessageType = "BuiltMine"
	PlayerMessageBuiltFactory                       PlayerMessageType = "BuiltFactory"
	PlayerMessageBuiltDefense                       PlayerMessageType = "BuiltDefense"
	PlayerMessageBuiltShip                          PlayerMessageType = "BuiltShip"
	PlayerMessageBuiltStarbase                      PlayerMessageType = "BuiltStarbase"
	PlayerMessageBuiltMineralPacket                 PlayerMessageType = "BuiltMineralPacket"
	PlayerMessageBuiltTerraform                     PlayerMessageType = "BuiltTerraform"
	PlayerMessageFleetOrdersComplete                PlayerMessageType = "FleetOrdersComplete"
	PlayerMessageFleetEngineFailure                 PlayerMessageType = "FleetEngineFailure"
	PlayerMessageFleetOutOfFuel                     PlayerMessageType = "FleetOutOfFuel"
	PlayerMessageFleetGeneratedFuel                 PlayerMessageType = "FleetGeneratedFuel"
	PlayerMessageFleetScrapped                      PlayerMessageType = "FleetScrapped"
	PlayerMessageFleetMerged                        PlayerMessageType = "FleetMerged"
	PlayerMessageFleetInvalidMergeNotFleet          PlayerMessageType = "FleetInvalidMergeNotFleet"
	PlayerMessageFleetInvalidMergeUnowned           PlayerMessageType = "FleetInvalidMergeUnowned"
	PlayerMessageFleetPatrolTargeted                PlayerMessageType = "FleetPatrolTargeted"
	PlayerMessageFleetInvalidRouteNotFriendlyPlanet PlayerMessageType = "FleetInvalidRouteNotFriendlyPlanet"
	PlayerMessageFleetInvalidRouteNotPlanet         PlayerMessageType = "FleetInvalidRouteNotPlanet"
	PlayerMessageFleetInvalidRouteNoRouteTarget     PlayerMessageType = "FleetInvalidRouteNoRouteTarget"
	PlayerMessageFleetInvalidTransport              PlayerMessageType = "FleetInvalidTransport"
	PlayerMessageFleetRoute                         PlayerMessageType = "FleetRoute"
	PlayerMessageInvalid                            PlayerMessageType = "Invalid"
	PlayerMessagePlanetColonized                    PlayerMessageType = "PlanetColonized"
	PlayerMessageGainTechLevel                      PlayerMessageType = "GainTechLevel"
	PlayerMessageMyPlanetBombed                     PlayerMessageType = "MyPlanetBombed"
	PlayerMessageMyPlanetRetroBombed                PlayerMessageType = "MyPlanetRetroBombed"
	PlayerMessageEnemyPlanetBombed                  PlayerMessageType = "EnemyPlanetBombed"
	PlayerMessageEnemyPlanetRetroBombed             PlayerMessageType = "EnemyPlanetRetroBombed"
	PlayerMessageMyPlanetInvaded                    PlayerMessageType = "MyPlanetInvaded"
	PlayerMessageEnemyPlanetInvaded                 PlayerMessageType = "EnemyPlanetInvaded"
	PlayerMessageBattle                             PlayerMessageType = "Battle"
	PlayerMessageCargoTransferred                   PlayerMessageType = "CargoTransferred"
	PlayerMessageMinesSwept                         PlayerMessageType = "MinesSwept"
	PlayerMessageMinesLaid                          PlayerMessageType = "MinesLaid"
	PlayerMessageMineFieldHit                       PlayerMessageType = "MineFieldHit"
	PlayerMessageFleetDumpedCargo                   PlayerMessageType = "FleetDumpedCargo"
	PlayerMessageFleetStargateDamaged               PlayerMessageType = "FleetStargateDamaged"
	PlayerMessageMineralPacketCaught                PlayerMessageType = "MineralPacketCaught"
	PlayerMessageMineralPacketDamage                PlayerMessageType = "MineralPacketDamage"
	PlayerMessageMineralPacketLanded                PlayerMessageType = "MineralPacketLanded"
	PlayerMessageVictor                             PlayerMessageType = "Victor"
	PlayerMessageFleetReproduce                     PlayerMessageType = "FleetReproduce"
	PlayerMessageRandomMineralDeposit               PlayerMessageType = "RandomMineralDeposit"
	PlayerMessagePermaform                          PlayerMessageType = "Permaform"
	PlayerMessageInstaform                          PlayerMessageType = "Instaform"
	PlayerMessagePacketTerraform                    PlayerMessageType = "PacketTerraform"
	PlayerMessagePacketPermaform                    PlayerMessageType = "PacketPermaform"
	PlayerMessageRemoteMined                        PlayerMessageType = "RemoteMined"
)

type Messager interface {
	homePlanet(player *Player, planet *Planet)
}

type messageClient struct {
}

var messager = messageClient{}

func (m *messageClient) error(player *Player, err error) {
	text := fmt.Sprintf("Something went wrong on the server. Please contact the administrator, %v", err)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageError, Text: text})
}

func (m *messageClient) homePlanet(player *Player, planet *Planet) {
	text := fmt.Sprintf("Your home planet is %s. Your people are ready to leave the nest and explore the universe.  Good luck.", planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageHomePlanet, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) longMessage(player *Player) {
	text := "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageHomePlanet, Text: text})
}

func (m *messageClient) minesBuilt(player *Player, planet *Planet, num int) {
	text := fmt.Sprintf("You have built %d mine(s) on %s.", num, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageBuiltMine, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) factoriesBuilt(player *Player, planet *Planet, num int) {
	text := fmt.Sprintf("You have built %d factory(s) on %s.", num, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageBuiltFactory, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) defensesBuilt(player *Player, planet *Planet, num int) {
	text := fmt.Sprintf("You have built %d defense(s) on %s.", num, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageBuiltFactory, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) fleetBuilt(player *Player, planet *Planet, fleet *Fleet, num int) {

	var text string

	if num == 1 {
		text = fmt.Sprintf("Your starbase at %s has built a new %s.", planet.Name, fleet.BaseName)
	} else {
		text = fmt.Sprintf("Your starbase at %s has built %d new %ss.", planet.Name, num, fleet.BaseName)
	}
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageBuiltShip, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum})
}

func (m *messageClient) fleetTransportedCargo(player *Player, fleet *Fleet, dest cargoHolder, cargoType CargoType, transferAmount int) {
	text := ""
	if cargoType == Colonists {
		if transferAmount < 0 {
			text = fmt.Sprintf("%s has beamed %d %s from %s", fleet.Name, -transferAmount*100, cargoType, dest.getMapObject().Name)
		} else {
			text = fmt.Sprintf("%s has beamed %d %s to %s", fleet.Name, transferAmount*100, cargoType, dest.getMapObject().Name)
		}
	} else {
		units := "kT"
		if cargoType == Fuel {
			units = "mg"
		}
		if transferAmount < 0 {
			text = fmt.Sprintf("%s has loaded %d%s of %s from %s", fleet.Name, -transferAmount, units, cargoType, dest.getMapObject().Name)
		} else {
			text = fmt.Sprintf("%s has unloaded %d%s of %s to %s", fleet.Name, transferAmount, units, cargoType, dest.getMapObject().Name)
		}
	}
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageCargoTransferred, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum})
}

func (m *messageClient) fleetInvalidLoadCargo(player *Player, fleet *Fleet, dest cargoHolder, cargoType CargoType, transferAmount int) {
	text := fmt.Sprintf("%s attempted to load %dkT of %v from %s, but you do not own %s", fleet.Name, transferAmount, cargoType, dest.getMapObject().Name, dest.getMapObject().Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetInvalidTransport, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum})

}

func (m *messageClient) fleetEngineFailure(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s was unable to engage it's engines due to balky equipment. Engineers think they have the problem fixed for the time being.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetEngineFailure, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum})
}

func (m *messageClient) fleetOutOfFuel(player *Player, fleet *Fleet, warpFactor int) {
	text := fmt.Sprintf("%s has run out of fuel. The fleet's speed has been decreased to Warp %d.", fleet.Name, warpFactor)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetOutOfFuel, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum})
}

func (m *messageClient) fleetGeneratedFuel(player *Player, fleet *Fleet, fuelGenerated int) {
	text := fmt.Sprintf("%s's ram scoops have produced %dmg of fuel from interstellar hydrogen.", fleet.Name, fuelGenerated)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetGeneratedFuel, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum})
}

func (m *messageClient) fleetScrapped(player *Player, fleet *Fleet, totalMinerals int, resources int, planet *Planet) {
	var text string
	if planet != nil {
		if planet.Spec.HasStarbase {
			text = fmt.Sprintf("%s has been dismantled for %dkT of minerals which have been deposited on %s.", fleet.Name, totalMinerals, planet.Name)
		} else {
			text = fmt.Sprintf("%s has been dismantled for %dkT of minerals at the starbase orbiting %s.", fleet.Name, totalMinerals, planet.Name)
		}
		if resources > 0 {
			text += fmt.Sprintf(" Ultimate recycling has also made %d resources available for immediate use (less if other ships were scrapped here this year).", resources)
		}
		player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetScrapped, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
	} else {
		text = fmt.Sprintf("%s has been dismantled. The scrap was left in deep space.", fleet.Name)
		player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetScrapped, Text: text})
	}
}

func (m *messageClient) fleetMerged(player *Player, fleet *Fleet, mergedInto *Fleet) {
	text := fmt.Sprintf("%s has been merged into %s.", fleet.Name, mergedInto.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetMerged, Text: text, TargetType: TargetFleet, TargetNum: mergedInto.Num, TargetPlayerNum: mergedInto.Num})
}

func (m *messageClient) fleetInvalidMergeNotFleet(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s was unable to complete it's merge orders as the waypoint destination wasn't a fleet.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetInvalidMergeNotFleet, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) fleetInvalidMergeNotOwned(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s was unable to complete it's merge orders as the destination fleet wasn't one of yours.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetInvalidMergeUnowned, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) fleetPatrolTargeted(player *Player, fleet *Fleet, target *Fleet) {
	text := fmt.Sprintf("Your patrolling %s has targeted %s for intercept.", fleet.Name, target.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetPatrolTargeted, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) fleetInvalidRouteNotPlanet(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s could not be routed because it is not at a planet.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetInvalidRouteNotPlanet, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) fleetInvalidRouteNotFriendlyPlanet(player *Player, fleet *Fleet, planet *Planet) {
	text := fmt.Sprintf("%s could not be routed because you are not friends with the owners of %s", fleet.Name, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetInvalidRouteNotFriendlyPlanet, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) fleetInvalidRouteNoRouteTarget(player *Player, fleet *Fleet, planet *Planet) {
	text := fmt.Sprintf("%s could not be routed because %s has no route set.", fleet.Name, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetInvalidRouteNoRouteTarget, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) fleetRouted(player *Player, fleet *Fleet, planet *Planet, target string) {
	text := fmt.Sprintf("%s has been routed by the citizens of %s to %s", fleet.Name, planet.Name, target)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetRoute, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) fleetStargateInvalidSource(player *Player, fleet *Fleet, wp0 Waypoint) {
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageInvalid,
		Text:            fmt.Sprintf("%s attempted to use a stargate at %s, but no stargate exists there.", fleet.Name, wp0.TargetName),
		TargetType:      TargetFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
	})
}

func (m *messageClient) fleetStargateInvalidSourceOwner(player *Player, fleet *Fleet, wp0, wp1 Waypoint) {
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageInvalid,
		Text:            fmt.Sprintf("%s attempted to use a stargate at %s, but could not because the starbase is not owned by you or a friend of yours.", fleet.Name, wp0.TargetName),
		TargetType:      TargetFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
	})
}

func (m *messageClient) fleetStargateInvalidDest(player *Player, fleet *Fleet, wp0, wp1 Waypoint) {
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageInvalid,
		Text:            fmt.Sprintf("%s attempted to use a stargate at %s to reach %s, but no stargate could be detected at the destination.", fleet.Name, wp0.TargetName, wp1.TargetName),
		TargetType:      TargetFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
	})
}

func (m *messageClient) fleetStargateInvalidDestOwner(player *Player, fleet *Fleet, wp0, wp1 Waypoint) {
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageInvalid,
		Text:            fmt.Sprintf("%s attempted to use a stargate at %s to reach %s, but could not because the destination starbase is not owned by you or a friend of yours.", fleet.Name, wp0.TargetName, wp1.TargetName),
		TargetType:      TargetFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
	})
}

func (m *messageClient) fleetStargateInvalidRange(player *Player, fleet *Fleet, wp0, wp1 Waypoint, totalDist float64) {
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageInvalid,
		Text:            fmt.Sprintf("%s attempted to use a stargate at %s to reach %s, but the distance of %.1f l.y. was outside the max range of the stargates.", fleet.Name, wp0.TargetName, wp1.TargetName, totalDist),
		TargetType:      TargetFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
	})
}

func (m *messageClient) fleetStargateInvalidMass(player *Player, fleet *Fleet, wp0, wp1 Waypoint) {
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageInvalid,
		Text:            fmt.Sprintf("%s attempted to use a stargate at %s to reach %s, but your ships are too massive.", fleet.Name, wp0.TargetName, wp1.TargetName),
		TargetType:      TargetFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
	})
}

func (m *messageClient) fleetStargateInvalidColonists(player *Player, fleet *Fleet, wp0 Waypoint, wp1 Waypoint) {
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageInvalid,
		Text:            fmt.Sprintf("%s attempted to use a stargate at %s to reach %s, but you are carrying colonists and can't drop them off as you don't own the planet.", fleet.Name, wp0.TargetName, wp1.TargetName),
		TargetType:      TargetFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
	})
}

func (m *messageClient) fleetStargateDumpedCargo(player *Player, fleet *Fleet, wp0 Waypoint, wp1 Waypoint, cargo Cargo) {
	var text string
	if cargo.HasColonists() && cargo.HasMinerals() {
		text = fmt.Sprintf("%s has unloaded %d colonists and %dkt of minerals in preparation for jumping through the stargate at %s to reach %s.", fleet.Name, cargo.Colonists*100, cargo.Total()-cargo.Colonists, wp0.TargetName, wp1.TargetName)
	} else if cargo.HasColonists() {
		text = fmt.Sprintf("%s has unloaded %d colonists in preparation for jumping through the stargate at %s to reach %s.", fleet.Name, cargo.Colonists*100, wp0.TargetName, wp1.TargetName)
	} else {
		text = fmt.Sprintf("%s has unloaded %dkt of minerals in preparation for jumping through the stargate at %s to reach %s.", fleet.Name, cargo.Total(), wp0.TargetName, wp1.TargetName)
	}
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageInvalid,
		Text:            text,
		TargetType:      TargetFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
	})
}

func (m *messageClient) fleetStargateDestroyed(player *Player, fleet *Fleet, wp0 Waypoint, wp1 Waypoint) {
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageFleetStargateDamaged,
		Text:            fmt.Sprintf("Heedless to the danger, %s attempted to use the stargate at %s to reach %s. The fleet never arrived. The distance or mass must have been too great.", fleet.Name, wp0.TargetName, wp1.TargetName),
		TargetType:      TargetFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
	})
}

func (m *messageClient) fleetStargateDamaged(player *Player, fleet *Fleet, wp0 Waypoint, wp1 Waypoint, damage int, startingShips int, shipsLostToDamage int, shipsLostToTheVoid int) {
	totalShipsLost := shipsLostToDamage + shipsLostToTheVoid
	var text string
	if totalShipsLost == 0 {
		text = fmt.Sprintf("%s used the stargate at %s to reach %s losing no ships but suffering %d dp of damage. They exceeded the capability of the gates.", fleet.Name, wp0.TargetName, wp1.TargetName, damage)
	} else if totalShipsLost < 5 {
		text = fmt.Sprintf("%s used the stargate at %s to reach %s losing only %d ship%s to the treacherous void. They were fortunate. They exceeded the capability of the gates.", fleet.Name, wp0.TargetName, wp1.TargetName, totalShipsLost, func() string {
			if totalShipsLost == 1 {
				return ""
			} else {
				return "s"
			}
		}())
	} else if totalShipsLost >= 5 && totalShipsLost <= 10 {
		text = fmt.Sprintf("%s used the stargate at %s to reach %s losing %d ships to the unforgiving void. Exceeding the capability of your stargates is not recommended.", fleet.Name, wp0.TargetName, wp1.TargetName, totalShipsLost)
	} else if totalShipsLost >= 10 && totalShipsLost <= 50 {
		text = fmt.Sprintf("%s used the stargate at %s to reach %s unfortunately losing %d ships to the great unknown. Exceeding the capability of your stargates is dangerous.", fleet.Name, wp0.TargetName, wp1.TargetName, totalShipsLost)
	} else if totalShipsLost >= 50 {
		text = fmt.Sprintf("%s used the stargate at %s to reach %s losing an unbelievable %d ships. The jump was far in excess of the capabilities of starbases involved..", fleet.Name, wp0.TargetName, wp1.TargetName, totalShipsLost)
	}
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageFleetStargateDamaged,
		Text:            text,
		TargetType:      TargetFleet,
		TargetPlayerNum: fleet.PlayerNum,
		TargetNum:       fleet.Num,
	})
}

func (m *messageClient) fleetReproduce(player *Player, fleet *Fleet, colonistsGrown int, planet *Planet, over int) {
	var text string
	if planet == nil || over == 0 {
		text = fmt.Sprintf("Your colonists in %s have made good use of their time increasing their on-board number by %d colonists.", fleet.Name, colonistsGrown)
	} else {
		text = fmt.Sprintf("Breeding activities on %s have overflowed living space. %d colonists have been beamed down to %s.", fleet.Name, over, planet.Name)
	}
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetReproduce, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})

}

func (m *messageClient) fleetCompletedAssignedOrders(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageFleetScrapped, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) fleetHitMineField(player *Player, fleet *Fleet, fleetPlayer *Player, mineField *MineField, damage int, shipsDestroyed int) {
	var text string
	if fleet.PlayerNum == player.Num {
		// it's our fleet, it must be someone else's minefield
		if fleet.Spec.TotalShips <= shipsDestroyed {
			text = fmt.Sprintf("%s has been annihilated in a %s mine field at %v",
				fleet.Name, mineField.Type, mineField.Position)
		} else {
			text = fmt.Sprintf("%s has been stopped in a %s mine field at %v.",
				fleet.Name, mineField.Type, mineField.Position)
			if damage > 0 {
				if shipsDestroyed > 0 {
					text += fmt.Sprintf(" Your fleet has taken %d damage points and %d ships were destroyed.",
						damage, shipsDestroyed)
				} else {
					text += fmt.Sprintf(" Your fleet has taken %d damage points but none of your ships were destroyed.",
						damage)
				}
			} else {
				text = fmt.Sprintf("%s has been stopped in a %s mine field at %sv",
					fleet.Name, mineField.Type, mineField.Position)
			}
		}
	} else {
		// it's not our fleet, it must be our minefield
		if fleet.Spec.TotalShips <= shipsDestroyed {
			text = fmt.Sprintf("%s %s has been annihilated in your %s mine field at %v",
				fleetPlayer.Race.PluralName, fleet.Name, mineField.Type, mineField.Position)
		} else {
			text = fmt.Sprintf("%s %s has been stopped in your %s mine field at %v.",
				fleetPlayer.Race.PluralName, fleet.Name, mineField.Type, mineField.Position)
			if damage > 0 {
				if shipsDestroyed > 0 {
					text += fmt.Sprintf(" Your mines have inflicted %d damage points and destroyed %d ships.",
						damage, shipsDestroyed)
				} else {
					text += fmt.Sprintf(" Your mines have inflicted %d damage points but you didn't manage to destroy any ships.",
						damage)
				}
			} else {
				text = fmt.Sprintf("%s has been stopped in your %s mine field at %s.",
					fleet.Name, mineField.Type, mineField.Position)
			}
		}
	}

	player.Messages = append(player.Messages, PlayerMessage{
		Type:       PlayerMessageMineFieldHit,
		Text:       text,
		TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num,
	})

}

func (m *messageClient) fleetMinesLaidFailed(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s has attempted to lay mines. The order has been cancelled because the fleet has no mine layers.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) fleetMinesLaid(player *Player, fleet *Fleet, mineField *MineField, numMinesLaid int) {
	var text string
	if mineField.NumMines == numMinesLaid {
		text = fmt.Sprintf("%s has dispersed %d mines.", fleet.Name, numMinesLaid)
	} else {
		text = fmt.Sprintf("%s has increased a minefield by %d mines.", fleet.Name, numMinesLaid)
	}
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageMinesLaid, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: player.Num})
}

func (m *messageClient) fleetMineFieldSwept(player *Player, fleet *Fleet, mineField *MineField, numMinesSwept int) {
	var text string

	if fleet.PlayerNum == player.Num {
		text = fmt.Sprintf("%s has swept %d from a mineField at %v", fleet.Name, numMinesSwept, mineField.Position)
	} else {
		text = fmt.Sprintf("Someone has swept %d from your mineField at %v", numMinesSwept, mineField.Position)
	}

	targetType := TargetNone
	targetNum := None
	targetPlayerNum := None

	// this will be removed if the mines are gone, so target the fleet
	if mineField.NumMines <= 10 {
		if fleet.PlayerNum == player.Num {
			targetType = TargetFleet
			targetNum = fleet.Num
			targetPlayerNum = fleet.PlayerNum
		}
	} else {
		targetType = TargetMineField
		targetNum = mineField.Num
		targetPlayerNum = mineField.PlayerNum
	}

	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageMinesSwept,
		Text:            text,
		TargetType:      targetType,
		TargetNum:       targetNum,
		TargetPlayerNum: targetPlayerNum,
	})
}

func (m *messageClient) colonizeNonPlanet(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s has attempted to colonize a waypoint with no Planet.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum})
}

func (m *messageClient) colonizeOwnedPlanet(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s has attempted to colonize a planet that is already inhabited.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum})

}

func (m *messageClient) colonizeWithNoModule(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s has attempted to colonize a planet without a colonization module.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum})

}

func (m *messageClient) colonizeWithNoColonists(player *Player, fleet *Fleet) {
	text := fmt.Sprintf("%s has attempted to colonize a planet without bringing any colonists.", fleet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, TargetType: TargetFleet, TargetNum: fleet.Num, TargetPlayerNum: fleet.PlayerNum})
}

func (m *messageClient) planetColonized(player *Player, planet *Planet) {
	text := fmt.Sprintf("Your colonists are now in control of %s", planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePlanetColonized, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) planetInvadeEmpty(player *Player, planet *Planet, fleet *Fleet) {
	text := fmt.Sprintf("%s has attempted to invade %s, but the planet is uninhabited.", fleet.Name, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) planetInvaded(player *Player, planet *Planet, fleet *Fleet, planetOwner string, fleetOwner string, attackersKilled int, defendersKilled int) {
	var text string

	// use this formatter to get commas on the text
	p := message.NewPrinter(language.English)
	if player.Num == fleet.PlayerNum {
		if planet.PlayerNum == fleet.PlayerNum {
			// we invaded and won
			text = p.Sprintf("Your %s has successfully invaded %s planet %s killing off all colonists", fleet.Name, planetOwner, planet.Name)
		} else {
			// we invaded and lost
			text = p.Sprintf("Your %s tried to invade %s, but all of your colonists were killed by %s. You valiant fighters managed to kill %d of their colonists.", fleet.Name, planet.Name, planetOwner, defendersKilled)
		}
		player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageEnemyPlanetInvaded, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
	} else {
		if planet.PlayerNum == fleet.PlayerNum {
			// we were invaded, and lost
			text = p.Sprintf("%s %s has successfully invaded your planet %s, killing off all of your colonists", fleetOwner, fleet.Name, planet.Name)
		} else {
			// we were invaded, and lost
			text = p.Sprintf("%s %s tried to invade %s, but you were able to fend them off. You lost %d colonists in the invasion.", fleetOwner, fleet.Name, planet.Name, defendersKilled)
		}
		player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageMyPlanetInvaded, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
	}
}

func (m *messageClient) planetBombed(player *Player, planet *Planet, fleet *Fleet, planetOwner string, fleetOwner string, colonistsKilled int, minesDestroyed int, factoriesDestroyed int, defensesDestroyed int) {
	var text string

	if player.Num == fleet.PlayerNum {
		if planet.population() == 0 {
			text = fmt.Sprintf("Your %s has bombed %s %s killing off all colonists", fleet.Name, planetOwner, planet.Name)
		} else {
			text = fmt.Sprintf("Your %s has bombed %s %s killing %d colonists, and destroying %d mines, %d factories, and %d defenses.", fleet.Name, planetOwner, planet.Name, colonistsKilled, minesDestroyed, factoriesDestroyed, defensesDestroyed)
		}
		player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageEnemyPlanetBombed, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
	} else {
		if planet.population() == 0 {
			text = fmt.Sprintf("%s %s has bombed your %s killing off all colonists", fleetOwner, fleet.Name, planet.Name)
		} else {
			text = fmt.Sprintf("%s %s has bombed your %s killing %d colonists, and destroying %d mines, %d factories, and %d defenses.", fleetOwner, fleet.Name, planet.Name, colonistsKilled, minesDestroyed, factoriesDestroyed, defensesDestroyed)
		}

		player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageMyPlanetBombed, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
	}
}

func (m *messageClient) planetSmartBombed(player *Player, planet *Planet, fleet *Fleet, planetOwner string, fleetOwner string, colonistsKilled int) {
	var text string

	if player.Num == fleet.PlayerNum {
		if planet.population() == 0 {
			text = fmt.Sprintf("Your fleet %s has bombed %s planet %s with smart bombs killing all colonists", fleet.Name, planetOwner, planet.Name)
		} else {
			text = fmt.Sprintf("Your %s has bombed %s planet %s with smart bombs killing %d colonists.", fleet.Name, planetOwner, planet.Name, colonistsKilled)
		}
		player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageEnemyPlanetBombed, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
	} else {
		if planet.population() == 0 {
			text = fmt.Sprintf("%s %s has bombed your %s with smart bombs killing all colonists", fleetOwner, fleet.Name, planet.Name)
		} else {
			text = fmt.Sprintf("%s %s has bombed your %s with smart bombs killing %d colonists.", fleetOwner, fleet.Name, planet.Name, colonistsKilled)
		}
		player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageMyPlanetBombed, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
	}
}

func (m *messageClient) planetRetroBombed(player *Player, planet *Planet, fleet *Fleet, planetOwner string, fleetOwner string, unterraformAmount Hab) {
	var text string

	if player.Num == fleet.PlayerNum {
		text = fmt.Sprintf("Your fleet %s has retro-bombed %s planet %s, undoing %d%% of its terraforming.", fleet.Name, planetOwner, planet.Name, unterraformAmount.absSum())
		player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageEnemyPlanetRetroBombed, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
	} else {
		text = fmt.Sprintf("%s %s has retro-bombed your %s, undoing %d%% of its terraforming.", fleetOwner, fleet.Name, planet.Name, unterraformAmount.absSum())
		player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageEnemyPlanetRetroBombed, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
	}
}

func (m *messageClient) battle(player *Player, planet *Planet, battle *BattleRecord) {
	var text string

	location := fmt.Sprintf("Space (%0f, %0f)", battle.Position.X, battle.Position.Y)
	if planet != nil {
		location = planet.Name
	}
	text = fmt.Sprintf("A battle took place at %s.", location)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageBattle, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num, BattleNum: battle.Num})
}

func (m *messageClient) techLevel(player *Player, field TechField, level int, nextField TechField) {
	text := fmt.Sprintf("Your scientists have completed research into Tech Level %d for %v.  They will continue their efforts in the %v field.", level, field, nextField)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageGainTechLevel, Text: text})
}

func (m *messageClient) playerDiscovered(player *Player, otherPlayer *Player) {
	text := fmt.Sprintf("You have discovered a new species, the %s. You are not alone in the universe!", otherPlayer.Race.PluralName)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePlayerDiscovery, Text: text})
}

func (m *messageClient) permaform(player *Player, planet *Planet, habType HabType, change int) {
	changeText := "decreased"
	if change > 0 {
		changeText = "increased"
	}
	text := fmt.Sprintf("Your race has permanently %s the %s on %s.", changeText, habType, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePermaform, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) packetTerraform(player *Player, planet *Planet, habType HabType, change int) {
	changeText := ""
	if change > 0 {
		changeText = "increased"
	} else {
		changeText = "decreased"
	}

	newValueText := ""
	newValue := planet.Hab.Get(habType)
	switch habType {
	case Grav:
		newValueText = gravString(newValue)
	case Temp:
		newValueText = tempString(newValue)
	case Rad:
		newValueText = radString(newValue)
	}

	text := fmt.Sprintf("Your mineral packet hitting %s has %s the %s to %s", planet.Name, changeText, habType, newValueText)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePacketTerraform, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) packetPermaform(player *Player, planet *Planet, habType HabType, change int) {
	changeText := ""
	if change > 0 {
		changeText = "increased"
	} else {
		changeText = "decreased"
	}
	text := fmt.Sprintf("Your mineral packet has permanently %s the %s on %s.", changeText, habType, planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessagePacketPermaform, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) instaform(player *Player, planet *Planet, terraformAmount Hab) {
	text := fmt.Sprintf("Your race has instantly terraformed %s up to optimal conditions.", planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInstaform, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) terraform(player *Player, planet *Planet, habType HabType, change int) {
	changeText := "decreased"
	if change > 0 {
		changeText = "increased"
	}

	var newValueText string
	newValue := planet.Hab.Get(habType)
	switch habType {
	case Grav:
		newValueText = gravString(newValue)
	case Temp:
		newValueText = tempString(newValue)
	case Rad:
		newValueText = radString(newValue)
	}

	text := fmt.Sprintf("Your terraforming efforts on %s have %s the %s to %s", planet.Name, changeText, habType, newValueText)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageBuiltTerraform, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) remoteMineNoMiners(player *Player, fleet *Fleet, planet *Planet) {
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageInvalid,
		Text:            fmt.Sprintf("%s had orders to mine %s, but the fleet doesn't have any remote mining modules. The order has been canceled.", fleet.Name, planet.Name),
		TargetType:      TargetFleet,
		TargetNum:       fleet.Num,
		TargetPlayerNum: fleet.PlayerNum,
	})
}

func (m *messageClient) remoteMineInhabited(player *Player, fleet *Fleet, planet *Planet) {
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageInvalid,
		Text:            fmt.Sprintf("Remote mining robots from %s had orders to mine %s, but the planet is inhabited. The order has been canceled.", fleet.Name, planet.Name),
		TargetType:      TargetFleet,
		TargetNum:       fleet.Num,
		TargetPlayerNum: fleet.PlayerNum,
	})
}

func (m *messageClient) remoteMineDeepSpace(player *Player, fleet *Fleet) {
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageInvalid,
		Text:            fmt.Sprintf("Remote mining robots from %s had orders to mine in deep space. The order has been canceled.", fleet.Name),
		TargetType:      TargetFleet,
		TargetNum:       fleet.Num,
		TargetPlayerNum: fleet.PlayerNum,
	})
}

func (m *messageClient) remoteMined(player *Player, fleet *Fleet, planet *Planet, mineral Mineral) {
	player.Messages = append(player.Messages, PlayerMessage{
		Type:            PlayerMessageRemoteMined,
		Text:            fmt.Sprintf("%s has remote mined %s, extracting %dkT of ironium, %dkT of boranium, and %dkT of germanium.", fleet.Name, planet.Name, mineral.Ironium, mineral.Boranium, mineral.Germanium),
		TargetType:      TargetFleet,
		TargetNum:       fleet.Num,
		TargetPlayerNum: fleet.PlayerNum,
	})
}

func (m *messageClient) mineralPacket(player *Player, planet *Planet, packet *MineralPacket, target string) {
	text := fmt.Sprintf("%s has produced a mineral packet which has a destination of %s", planet.Name, target)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageBuiltMineralPacket, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) mineralPacketArrived(player *Player, planet *Planet, packet *MineralPacket) {
	text := fmt.Sprintf("Your mineral packet containing %dkT of minerals has landed at %s.", packet.Cargo.Total(), planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageMineralPacketLanded, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) mineralPacketCaught(player *Player, planet *Planet, packet *MineralPacket) {
	text := fmt.Sprintf("Your mass accelerator at %s has successfully captured a packet containing %dkT of minerals.", planet.Name, packet.Cargo.Total())
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageMineralPacketCaught, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) buildMineralPacketNoMassDriver(player *Player, planet *Planet) {
	text := fmt.Sprintf("You have attempted to build a mineral packet on %s, but you have no Starbase equipped with a mass driver on this planet. Production for this planet has been cancelled.", planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) buildMineralPacketNoTarget(player *Player, planet *Planet) {
	text := fmt.Sprintf("You have attempted to build a mineral packet on %s, but you have not specified a target. The minerals have been returned to the planet and production has been cancelled.", planet.Name)
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageInvalid, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (m *messageClient) mineralPacketDamage(player *Player, planet *Planet, packet *MineralPacket, colonistsKilled, defensesDestroyed int) {
	var text string
	if planet.Spec.HasStarbase && planet.starbase.Spec.HasMassDriver {
		if defensesDestroyed == 0 {
			text = fmt.Sprintf("Your mass accelerator at %s was partially successful in capturing a %dkT mineral packet. Unable to completely slow the packet, %d of your colonists were killed in the collision.", planet.Name, packet.Cargo.Total(), colonistsKilled)
		} else {
			text = fmt.Sprintf("Your mass accelerator at %s was partially successful in capturing a %dkT mineral packet. Unfortunately, %d of your colonists and %d of your defenses were destroyed in the collision.", planet.Name, packet.Cargo.Total(), colonistsKilled, defensesDestroyed)
		}
	} else {
		if planet.population() == 0 {
			text = fmt.Sprintf("%s was annihilated by a mineral packet. All of your colonists were killed.", planet.Name)
		} else if defensesDestroyed == 0 {
			text = fmt.Sprintf("%s was bombarded with a %dkT mineral packet. %d of your colonists were killed by the collision.", planet.Name, packet.Cargo.Total(), colonistsKilled)
		} else {
			text = fmt.Sprintf("%s was bombarded with a %dkT mineral packet. %d of your colonists and %d of your defenses were destroyed by the collision.", planet.Name, packet.Cargo.Total(), colonistsKilled, defensesDestroyed)
		}
	}
	player.Messages = append(player.Messages, PlayerMessage{Type: PlayerMessageMineralPacketDamage, Text: text, TargetType: TargetPlanet, TargetNum: planet.Num})
}

func (mc *messageClient) victory(player *Player, victor *Player) {
	var text string
	if player.Num == victor.Num {
		text = "You have been declared the winner of this game. You may continue to play though, if you wish to really rub everyone's nose in your grand victory."
	} else {
		text = fmt.Sprintf("The forces of %s have been declared the winner of this game. You are advised to accept their supremacy, though you may continue the fight.", player.Race.PluralName)
	}
	// Victory messages are always the first message of the year
	player.Messages = append([]PlayerMessage{PlayerMessage{Type: PlayerMessageVictor, Text: text}}, player.Messages...)
}
