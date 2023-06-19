package cs

import (
	"fmt"
	"time"
)

type Tags map[string]string

type NewGamePlayerType string

const (
	NewGamePlayerTypeHost   NewGamePlayerType = "Host"
	NewGamePlayerTypeInvite NewGamePlayerType = "Invite"
	NewGamePlayerTypeOpen   NewGamePlayerType = "Open"
	NewGamePlayerTypeAI     NewGamePlayerType = "AI"
)

type AIDifficulty string

const (
	AIDifficultyNone   AIDifficulty = ""
	AIDifficultyEasy   AIDifficulty = "Easy"
	AIDifficultyNormal AIDifficulty = "Normal"
	AIDifficultyHard   AIDifficulty = "Hard"
)

type NewGamePlayer struct {
	Type         NewGamePlayerType `json:"type,omitempty"`
	RaceID       int64             `json:"raceId,omitempty"`
	AIDifficulty AIDifficulty      `json:"aiDifficulty,omitempty"`
	Color        string            `json:"color,omitempty"`
}

type GameSettings struct {
	Name                         string            `json:"name"`
	QuickStartTurns              int               `json:"quickStartTurns"`
	Size                         Size              `json:"size"`
	Density                      Density           `json:"density"`
	PlayerPositions              PlayerPositions   `json:"playerPositions"`
	RandomEvents                 bool              `json:"randomEvents"`
	ComputerPlayersFormAlliances bool              `json:"computerPlayersFormAlliances"`
	PublicPlayerScores           bool              `json:"publicPlayerScores"`
	StartMode                    GameStartMode     `json:"startMode"`
	VictoryConditions            VictoryConditions `json:"victoryConditions"`
	Players                      []NewGamePlayer   `json:"players,omitempty"`
}

type Game struct {
	ID                           int64             `json:"id" header:"ID"`
	CreatedAt                    time.Time         `json:"createdAt"`
	UpdatedAt                    time.Time         `json:"updatedAt"`
	Name                         string            `json:"name" header:"Name"`
	HostID                       int64             `json:"hostId"`
	QuickStartTurns              int               `json:"quickStartTurns"`
	Size                         Size              `json:"size"`
	Density                      Density           `json:"density"`
	PlayerPositions              PlayerPositions   `json:"playerPositions"`
	RandomEvents                 bool              `json:"randomEvents"`
	ComputerPlayersFormAlliances bool              `json:"computerPlayersFormAlliances"`
	PublicPlayerScores           bool              `json:"publicPlayerScores"`
	StartMode                    GameStartMode     `json:"startMode"`
	Year                         int               `json:"year"`
	State                        GameState         `json:"state"`
	OpenPlayerSlots              uint              `json:"openPlayerSlots"`
	NumPlayers                   int               `json:"numPlayers"`
	VictoryConditions            VictoryConditions `json:"victoryConditions"`
	VictorDeclared               bool              `json:"victorDeclared"`
	Seed                         int64             `json:"seed"`
	Rules                        Rules             `json:"rules"`
	Area                         Vector            `json:"area,omitempty"`
}

// A game with players and a universe, used in universe and turn generation
type FullGame struct {
	*Game
	*Universe
	Players []*Player `json:"players,omitempty"`
}

type VictoryConditions struct {
	Conditions               []VictoryCondition `json:"conditions"`
	NumCriteriaRequired      int                `json:"numCriteriaRequired"`
	YearsPassed              int                `json:"yearsPassed"`
	OwnPlanets               int                `json:"ownPlanets"`
	AttainTechLevel          int                `json:"attainTechLevel"`
	AttainTechLevelNumFields int                `json:"attainTechLevelNumFields"`
	ExceedsScore             int                `json:"exceedsScore"`
	ExceedsSecondPlaceScore  int                `json:"exceedsSecondPlaceScore"`
	ProductionCapacity       int                `json:"productionCapacity"`
	OwnCapitalShips          int                `json:"ownCapitalShips"`
	HighestScoreAfterYears   int                `json:"highestScoreAfterYears"`
}

type Size string

const (
	SizeTiny       Size = "Tiny"
	SizeTinyWide   Size = "TinyWide"
	SizeSmall      Size = "Small"
	SizeSmallWide  Size = "SmallWide"
	SizeMedium     Size = "Medium"
	SizeMediumWide Size = "MediumWide"
	SizeLarge      Size = "Large"
	SizeLargeWide  Size = "LargeWide"
	SizeHuge       Size = "Huge"
	SizeHugeWide   Size = "HugeWide"
)

type Density string

const (
	DensitySparse Density = "Sparse"
	DensityNormal Density = "Normal"
	DensityDense  Density = "Dense"
	DensityPacked Density = "Packed"
)

type PlayerPositions string

const (
	PlayerPositionsClose    PlayerPositions = "Close"
	PlayerPositionsModerate PlayerPositions = "Moderate"
	PlayerPositionsFarther  PlayerPositions = "Farther"
	PlayerPositionsDistant  PlayerPositions = "Distant"
)

type GameStartMode string

const (
	GameStartModeNormal   GameStartMode = "Normal"
	GameStartModeMidGame  GameStartMode = "MidGame"
	GameStartModeLateGame GameStartMode = "LateGame"
	GameStartModeEndGame  GameStartMode = "EndGame"
)

type GameState string

const (
	GameStateSetup             GameState = "Setup"
	GameStateWaitingForPlayers GameState = "WaitingForPlayers"
	GameStateGeneratingTurn    GameState = "GeneratingTurn"
)

type VictoryCondition string

const (
	VictoryConditionOwnPlanets              VictoryCondition = "OwnPlanets"
	VictoryConditionAttainTechLevels        VictoryCondition = "AttainTechLevels"
	VictoryConditionExceedsScore            VictoryCondition = "ExceedsScore"
	VictoryConditionExceedsSecondPlaceScore VictoryCondition = "ExceedsSecondPlaceScore"
	VictoryConditionProductionCapacity      VictoryCondition = "ProductionCapacity"
	VictoryConditionOwnCapitalShips         VictoryCondition = "OwnCapitalShips"
	VictoryConditionHighestScoreAfterYears  VictoryCondition = "HighestScoreAfterYears"
)

func NewGame() *Game {
	seed := time.Now().UnixNano()
	rules := NewRulesWithSeed(seed)
	return &Game{
		Name:            "A Barefoot Jaywalk",
		Size:            SizeSmall,
		Density:         DensityNormal,
		PlayerPositions: PlayerPositionsModerate,
		RandomEvents:    true,
		StartMode:       GameStartModeNormal,
		Year:            2400,
		State:           GameStateSetup,
		VictoryConditions: VictoryConditions{
			Conditions: []VictoryCondition{
				VictoryConditionOwnPlanets,
				VictoryConditionAttainTechLevels,
				VictoryConditionExceedsSecondPlaceScore,
			},
			NumCriteriaRequired:      1,
			YearsPassed:              50,
			OwnPlanets:               60,
			AttainTechLevel:          22,
			AttainTechLevelNumFields: 4,
			ExceedsScore:             11000,
			ExceedsSecondPlaceScore:  100,
			ProductionCapacity:       100000,
			OwnCapitalShips:          100,
			HighestScoreAfterYears:   100,
		},
		Seed:  seed,
		Rules: rules,
	}
}

// create a new GameSettings object for the default game
func NewGameSettings() *GameSettings {
	return &GameSettings{
		Name:            "A Barefoot Jaywalk",
		Size:            SizeSmall,
		Density:         DensityNormal,
		PlayerPositions: PlayerPositionsModerate,
		RandomEvents:    true,
		StartMode:       GameStartModeNormal,
		VictoryConditions: VictoryConditions{
			Conditions: []VictoryCondition{
				VictoryConditionOwnPlanets,
				VictoryConditionAttainTechLevels,
				VictoryConditionExceedsSecondPlaceScore,
			},
			NumCriteriaRequired:      1,
			YearsPassed:              50,
			OwnPlanets:               60,
			AttainTechLevel:          22,
			AttainTechLevelNumFields: 4,
			ExceedsScore:             11000,
			ExceedsSecondPlaceScore:  100,
			ProductionCapacity:       100000,
			OwnCapitalShips:          100,
			HighestScoreAfterYears:   100,
		},
	}
}

// set the name to this game
func (settings *GameSettings) WithName(name string) *GameSettings {
	settings.Name = name
	return settings
}

func (settings *GameSettings) WithSize(size Size) *GameSettings {
	settings.Size = size
	return settings
}

func (settings *GameSettings) WithDensity(density Density) *GameSettings {
	settings.Density = density
	return settings
}

// add a host to this game
func (settings *GameSettings) WithHost(raceID int64) *GameSettings {
	settings.Players = append(settings.Players, NewGamePlayer{Type: NewGamePlayerTypeHost, RaceID: raceID, Color: "#0000FF"})
	return settings
}

// Add a player slot open to any players
func (settings *GameSettings) WithOpenPlayerSlot() *GameSettings {
	settings.Players = append(settings.Players, NewGamePlayer{Type: NewGamePlayerTypeOpen})
	return settings
}

// Add an AI player
func (settings *GameSettings) WithAIPlayer(aiDifficulty AIDifficulty) *GameSettings {
	settings.Players = append(settings.Players, NewGamePlayer{Type: NewGamePlayerTypeAI, AIDifficulty: aiDifficulty})
	return settings
}

func (g *Game) String() string {
	return fmt.Sprintf("%s (%d)", g.Name, g.ID)
}

// update the game with new settings
func (g *Game) WithSettings(settings GameSettings) *Game {
	g.Name = settings.Name
	g.QuickStartTurns = settings.QuickStartTurns
	g.Size = settings.Size
	g.Density = settings.Density
	g.PlayerPositions = settings.PlayerPositions
	g.RandomEvents = settings.RandomEvents
	g.ComputerPlayersFormAlliances = settings.ComputerPlayersFormAlliances
	g.PublicPlayerScores = settings.PublicPlayerScores
	g.StartMode = settings.StartMode
	g.VictoryConditions = settings.VictoryConditions

	return g
}
