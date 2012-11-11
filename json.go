// This file is part of The AI Sandbox Go Bindings by errnoh.
// Copyright (c) 2012, errnoh@github
// License: See LICENSE file.
package aisandbox

type JSON_GameInfo struct {
	Class string `json:"__class__"`
	Value struct {
		Teams     map[string]*JSON_TeamInfo `json:"teams"` // map of team names to TeamInfo objects
		Team      string                    `json:"team"`
		EnemyTeam string                    `json:"enemyTeam"`
		Flags     map[string]*JSON_FlagInfo `json:"flags"` // map of team names to FlagInfo objects
		Bots      map[string]*JSON_BotInfo  `json:"bots"`  // map of bot names to BotInfo objects
		Match     *JSON_MatchInfo           `json:"match"` // MatchInfo object
	} `json:"__value__"`
}

// Parse JSON_GameInfo struct into more intuitive GameInfo struct
// before sending it to the commander
func (data *JSON_GameInfo) simplify() *GameInfo {
	v := data.Value
	own := v.Teams[v.Team].Value
	enemy := v.Teams[v.EnemyTeam].Value
	ownflag := v.Flags[own.Flag].Value
	enemyflag := v.Flags[enemy.Flag].Value
	match := v.Match.Value

	// BotInfo

	ownbots := make(map[string]*BotInfo)
	enemybots := make(map[string]*BotInfo)

	for _, name := range own.Members {
		ownbots[v.Bots[name].Value.Name] = &BotInfo{
			Name:            v.Bots[name].Value.Name,
			Team:            v.Bots[name].Value.Team,
			Position:        v.Bots[name].Value.Position,
			FacingDirection: v.Bots[name].Value.FacingDirection,
			Flag:            v.Bots[name].Value.Flag,
			CurrentAction:   v.Bots[name].Value.CurrentAction,
			State:           v.Bots[name].Value.State,
			Health:          v.Bots[name].Value.Health,
			SeenLast:        v.Bots[name].Value.SeenLast,
		}
	}

	for _, name := range enemy.Members {
		enemybots[v.Bots[name].Value.Name] = &BotInfo{
			Name:            v.Bots[name].Value.Name,
			Team:            v.Bots[name].Value.Team,
			Position:        v.Bots[name].Value.Position,
			FacingDirection: v.Bots[name].Value.FacingDirection,
			Flag:            v.Bots[name].Value.Flag,
			CurrentAction:   v.Bots[name].Value.CurrentAction,
			State:           v.Bots[name].Value.State,
			Health:          v.Bots[name].Value.Health,
			SeenLast:        v.Bots[name].Value.SeenLast,
		}

		for _, seenby := range v.Bots[name].Value.SeenBy {
			enemybots[name].SeenBy = append(enemybots[name].SeenBy, ownbots[seenby])
		}

		for _, visible := range v.Bots[name].Value.VisibleEnemies {
			enemybots[name].VisibleEnemies = append(enemybots[name].VisibleEnemies, ownbots[visible])
		}
	}

	for _, name := range own.Members {
		for _, seenby := range v.Bots[name].Value.SeenBy {
			ownbots[name].SeenBy = append(ownbots[name].SeenBy, enemybots[seenby])
		}

		for _, visible := range v.Bots[name].Value.VisibleEnemies {
			ownbots[name].VisibleEnemies = append(ownbots[name].VisibleEnemies, enemybots[visible])
		}
	}

	// FlagInfo

	ownflaginfo := &FlagInfo{
		Position:     ownflag.Position,
		Carrier:      enemybots[ownflag.Carrier],
		RespawnTimer: ownflag.RespawnTimer,
	}

	enemyflaginfo := &FlagInfo{
		Position:     enemyflag.Position,
		Carrier:      ownbots[enemyflag.Carrier],
		RespawnTimer: enemyflag.RespawnTimer,
	}

	// TeamInfo

	ownteaminfo := &TeamInfo{
		Name:              own.Name,
		Flag:              ownflaginfo,
		Members:           ownbots,
		FlagSpawnLocation: own.FlagSpawnLocation,
		FlagScoreLocation: own.FlagScoreLocation,
		BotSpawnArea:      own.BotSpawnArea,
	}

	enemyteaminfo := &TeamInfo{
		Name:              enemy.Name,
		Flag:              enemyflaginfo,
		Members:           enemybots,
		FlagSpawnLocation: enemy.FlagSpawnLocation,
		FlagScoreLocation: enemy.FlagScoreLocation,
		BotSpawnArea:      enemy.BotSpawnArea,
	}

	// MatchInfo

	matchinfo := &MatchInfo{
		TimeRemaining:     match.TimeRemaining,
		TimeToNextRespawn: match.TimeToNextRespawn,
	}

	for _, event := range match.CombatEvents {
		matchinfo.CombatEvents = append(matchinfo.CombatEvents, event.Value)
	}

	// GameInfo

	return &GameInfo{
		Team:      ownteaminfo,
		EnemyTeam: enemyteaminfo,
		Match:     matchinfo,
	}
}

type JSON_LevelInfo struct {
	Class string     `json:"__class__"`
	Value *LevelInfo `json:"__value__"`
}

type JSON_TeamInfo struct {
	Class string `json:"__class__"`
	Value struct {
		Name              string      `json:"name"`
		Flag              string      `json:"flag"`
		Members           []string    `json:"members"`           // list of bot names
		FlagSpawnLocation []float64   `json:"flagSpawnLocation"` // (may be removed as this is available in LevelInfo)
		FlagScoreLocation []float64   `json:"flagScoreLocation"` // (may be removed as this is available in LevelInfo)
		BotSpawnArea      [][]float64 `json:"flagSpawnArea"`     // min and max positions (may be removed as this is available in LevelInfo)
	} `json:"__value__"`
}

type JSON_FlagInfo struct {
	Class string `json:"__class__"`
	Value struct {
		Name         string    `json:"name"`
		Team         string    `json:"team"`
		Position     []float64 `json:"position"`
		Carrier      string    `json:"carrier, omitempty"` // optional bot name, null if the flag is not being carried
		RespawnTimer float64   `json:"respawnTimer"`
	} `json:"__value__"`
}

type JSON_BotInfo struct {
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

type JSON_MatchInfo struct {
	Class string `json:"__class__"`
	Value struct {
		TimeRemaining     float64                  `json:"timeRemaining"`
		TimeToNextRespawn float64                  `json:"timeToNextRespawn"`
		CombatEvents      []*JSON_MatchCombatEvent `json:"combatEvents"` // list of MatchCombatEvent objects
	} `json:"__value__"`
}

type JSON_MatchCombatEvent struct {
	Class string       `json:"__class__"`
	Value *CombatEvent `json:"__value__"`
}
