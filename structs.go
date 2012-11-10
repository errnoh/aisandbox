// This file is part of The AI Sandbox Go Bindings by errnoh.
// Copyright (c) 2012, errnoh@github
// License: See LICENSE file.
package aisandbox

type LevelInfo struct {
	Class string `json:"__class__"`
	Value struct {
		Width              float64                `json:"width"`
		Height             float64                `json:"height"`
		BlockHeights       [][]float64            `json:"blockHeights"`       // a 'width' list of 'height' lengthed list of integers
		TeamNames          []string               `json:"teamNames"`          // list of team names
		FlagSpawnLocations map[string][]float64   `json:"flagSpawnLocations"` // map of team name to position
		FlagScoreLocations map[string][]float64   `json:"flagScoreLocations"` // map of team name to position
		BotSpawnAreas      map[string][][]float64 `json:"botSpawnAreas"`      // map of team name to min and max positions
		FOVangle           float64                `json:"FOVangle"`
		CharacterRadius    float64                `json:"characterRadius"`
		WalkingSpeed       float64                `json:"walkingSpeed"`
		RunningSpeed       float64                `json:"runningSpeed"`
		FiringDistance     float64                `json:"firingDistance"`
	} `json:"__value__"`
}

type GameInfo struct {
	Class string `json:"__class__"`
	Value struct {
		Teams     map[string]*TeamInfo `json:"teams"` // map of team names to TeamInfo objects
		Team      string               `json:"team"`
		EnemyTeam string               `json:"enemyTeam"`
		Flags     map[string]*FlagInfo `json:"flags"` // map of team names to FlagInfo objects
		Bots      map[string]*BotInfo  `json:"bots"`  // map of bot names to BotInfo objects
		Match     *MatchInfo           `json:"match"` // MatchInfo object
	} `json:"__value__"`
}

type TeamInfo struct {
	Class string `json:"__class__"`
	Value struct {
		Name              string      `json:"name"`
		Flag              string      `json:"flag"`
		Members           []string    `json:"members"`           // list of bot names
		FlagSpawnLocation []float64   `json:"flagSpawnLocation"` // (may be removed as this is available in LevelInfo)
		FlagScoreLocation []float64   `json:"flagScoreLocation"` // (may be removed as this is available in LevelInfo)
		botSpawnArea      [][]float64 `json:"flagSpawnArea"`     // min and max positions (may be removed as this is available in LevelInfo)
	} `json:"__value__"`
}

type FlagInfo struct {
	Class string `json:"__class__"`
	Value struct {
		Name         string    `json:"name"`
		Team         string    `json:"team"`
		Position     []float64 `json:"position"`
		Carrier      string    `json:"carrier, omitempty"` // optional bot name, null if the flag is not being carried
		RespawnTimer float64   `json:"respawnTimer"`
	} `json:"__value__"`
}

type BotInfo struct {
	Class string `json:"__class__"`
	Value struct {
		Name            string    `json:"name"`
		Team            string    `json:"team"`
		Position        []float64 `json:"position, omitempty"`        // optional, null if the bot is not visible
		FacingDirection []float64 `json:"facingDirection, omitempty"` // optional, null if the bot is not visible
		Flag            string    `json:"flag, omitempty"`            // optional flag name, null if the bot is not carrying a flag
		CurrentAction   string    `json:"currentAction, omitempty"`   // optional current action name, null if the bot is not visible (will be removed)
		// values are 0 = unknown, 1 = idle, 2 = defending, 3 = moving, 4 = attacking, 5 = charging, 6 = shooting
		State          float64  `json:"state, omitempty"`    // optional current action name, null if the bot is not visible
		Health         float64  `json:"health, omitempty"`   // optional, null if the bot is not visible
		SeenLast       float64  `json:"seenlast, omitempty"` // time since the object was last seen, null if the object was never seen
		VisibleEnemies []string `json:"visibleEnemies"`      // list of bot names for bots which this bot can see
		SeenBy         []string `json:"seenBy"`              // list of bot names for bots which can see this bot
	} `json:"__value__"`
}

type MatchInfo struct {
	Class string `json:"__class__"`
	Value struct {
		TimeRemaining     float64             `json:"timeRemaining"`
		TimeToNextRespawn float64             `json:"timeToNextRespawn"`
		CombatEvents      []*MatchCombatEvent `json:"combatEvents"` // list of MatchCombatEvent objects
	} `json:"__value__"`
}

type MatchCombatEvent struct {
	Class string `json:"__class__"`
	Value struct {
		Type       float64 `json:"type"`                  // values are 0 = none, 1 = bot killed, 2 = flag picked up, 3 = flag dropped (more to be added soon)
		Instigator string  `json:"instigator, omitempty"` // optional bot name that caused the event, null if the event was automatic (eg flag reset, bot respawn)
		// can either be a FlagInfo or a BotInfo name
		Subject string `json:"subject"` // bot or flag name that was the subject of the event
		Time    string `json:"time"`
	} `json:"__value__"`
}

type Defend struct {
	Class string `json:"__class__"`
	Value struct {
		Bot             string    `json:"bot"`
		FacingDirection []float64 `json:"facingDirection"`
		Description     string    `json:"description"`
	} `json:"__value__"`
}

type Move struct {
	Class string `json:"__class__"`
	Value struct {
		Bot         string    `json:"bot"`
		Target      []float64 `json:"target"`
		Description string    `json:"description"`
	} `json:"__value__"`
}

type Attack struct {
	Class string `json:"__class__"`
	Value struct {
		Bot         string    `json:"bot"`
		Target      []float64 `json:"target"`
		LookAt      []float64 `json:"lookAt"`
		Description string    `json:"description"`
	} `json:"__value__"`
}

type Charge struct {
	Class string `json:"__class__"`
	Value struct {
		Bot         string    `json:"bot"`
		Target      []float64 `json:"target"`
		Description string    `json:"description"`
	} `json:"__value__"`
}
