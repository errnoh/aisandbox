// This file is part of The AI Sandbox Go Bindings by errnoh.
// Copyright (c) 2012, errnoh@github
// License: See LICENSE file.

package aisandbox

// NOTE: This file contains unexported things that are used when parsing JSON data from the game server.

import (
	"encoding/json"
)

// Workaround for null JSON strings
type nstring string

func (n *nstring) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*string)(n))
}

// Workaround for null JSON floats
type nfloat64 float64

func (n *nfloat64) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	return json.Unmarshal(b, (*float64)(n))
}

// Unexported structs where the raw JSON is parsed.

type json_GameInfo struct {
	Class string `json:"__class__"`
	Value struct {
		Teams     map[string]*json_TeamInfo `json:"teams"` // map of team names to TeamInfo objects
		Team      string                    `json:"team"`
		EnemyTeam string                    `json:"enemyTeam"`
		Flags     map[string]*json_FlagInfo `json:"flags"` // map of team names to FlagInfo objects
		Bots      map[string]*json_BotInfo  `json:"bots"`  // map of bot names to BotInfo objects
		Match     *json_MatchInfo           `json:"match"` // MatchInfo object
	} `json:"__value__"`
}

// Parse json_GameInfo struct into more intuitive GameInfo struct
// before sending it to the commander
func (data *json_GameInfo) simplify() *GameInfo {
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
			Flag:            string(v.Bots[name].Value.Flag),
			State:           float64(v.Bots[name].Value.State),
			Health:          float64(v.Bots[name].Value.Health),
			SeenLast:        float64(v.Bots[name].Value.SeenLast),
		}
	}

	for _, name := range enemy.Members {
		enemybots[v.Bots[name].Value.Name] = &BotInfo{
			Name:            v.Bots[name].Value.Name,
			Team:            v.Bots[name].Value.Team,
			Position:        v.Bots[name].Value.Position,
			FacingDirection: v.Bots[name].Value.FacingDirection,
			Flag:            string(v.Bots[name].Value.Flag),
			State:           float64(v.Bots[name].Value.State),
			Health:          float64(v.Bots[name].Value.Health),
			SeenLast:        float64(v.Bots[name].Value.SeenLast),
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
		Carrier:      enemybots[string(ownflag.Carrier)],
		RespawnTimer: ownflag.RespawnTimer,
	}

	enemyflaginfo := &FlagInfo{
		Position:     enemyflag.Position,
		Carrier:      ownbots[string(enemyflag.Carrier)],
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
		Score:             match.Scores[own.Name],
	}

	enemyteaminfo := &TeamInfo{
		Name:              enemy.Name,
		Flag:              enemyflaginfo,
		Members:           enemybots,
		FlagSpawnLocation: enemy.FlagSpawnLocation,
		FlagScoreLocation: enemy.FlagScoreLocation,
		BotSpawnArea:      enemy.BotSpawnArea,
		Score:             match.Scores[enemy.Name],
	}

	// MatchInfo

	matchinfo := &MatchInfo{
		TimeRemaining:     match.TimeRemaining,
		TimeToNextRespawn: match.TimeToNextRespawn,
		TimePassed:        match.TimePassed,
	}

	// TODO: map instigator field to target?
	for _, event := range match.CombatEvents {
		matchinfo.CombatEvents = append(
			matchinfo.CombatEvents,
			&CombatEvent{
				Type:       event.Value.Type,
				Instigator: string(event.Value.Instigator),
				Subject:    event.Value.Subject,
				Time:       event.Value.Time,
			},
		)
	}

	// GameInfo

	return &GameInfo{
		Team:      ownteaminfo,
		EnemyTeam: enemyteaminfo,
		Match:     matchinfo,
	}
}

type json_LevelInfo struct {
	Class string     `json:"__class__"`
	Value *LevelInfo `json:"__value__"`
}

type json_TeamInfo struct {
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

type json_FlagInfo struct {
	Class string `json:"__class__"`
	Value struct {
		Name         string    `json:"name"`
		Team         string    `json:"team"`
		Position     []float64 `json:"position"`
		Carrier      nstring   `json:"carrier, omitempty"` // optional bot name, null if the flag is not being carried
		RespawnTimer float64   `json:"respawnTimer"`
	} `json:"__value__"`
}

type json_BotInfo struct {
	Class string `json:"__class__"`
	Value struct {
		Name            string    `json:"name"`
		Team            string    `json:"team"`
		Position        []float64 `json:"position, omitempty"`        // optional, null if the bot is not visible
		FacingDirection []float64 `json:"facingDirection, omitempty"` // optional, null if the bot is not visible
		Flag            nstring   `json:"flag, omitempty"`            // optional flag name, null if the bot is not carrying a flag
		// values are 0 = unknown, 1 = idle, 2 = defending, 3 = moving, 4 = attacking, 5 = charging, 6 = shooting
		State          nfloat64 `json:"state, omitempty"`    // optional current action name, null if the bot is not visible
		Health         nfloat64 `json:"health, omitempty"`   // optional, null if the bot is not visible
		SeenLast       nfloat64 `json:"seenlast, omitempty"` // time since the object was last seen, null if the object was never seen
		VisibleEnemies []string `json:"visibleEnemies"`      // list of bot names for bots which this bot can see
		SeenBy         []string `json:"seenBy"`              // list of bot names for bots which can see this bot
	} `json:"__value__"`
}

type json_MatchInfo struct {
	Class string `json:"__class__"`
	Value struct {
		TimeRemaining     float64                  `json:"timeRemaining"`
		TimeToNextRespawn float64                  `json:"timeToNextRespawn"`
		CombatEvents      []*json_MatchCombatEvent `json:"combatEvents"` // list of MatchCombatEvent objects
		TimePassed        float64                  `json:"timePassed"`
		Scores            map[string]float64       `json:"scores"`
	} `json:"__value__"`
}

type json_MatchCombatEvent struct {
	Class string            `json:"__class__"`
	Value *json_CombatEvent `json:"__value__"`
}

type json_CombatEvent struct {
	Type       float64 `json:"type"`                  // values are 0 = none, 1 = bot killed, 2 = flag picked up, 3 = flag dropped (more to be added soon)
	Instigator nstring `json:"instigator, omitempty"` // optional bot name that caused the event, null if the event was automatic (eg flag reset, bot respawn)
	// can either be a FlagInfo or a BotInfo name
	Subject string  `json:"subject"` // bot or flag name that was the subject of the event
	Time    float64 `json:"time"`
}

type json_ConnectServer struct {
	Class string `json:"__class__"`
	Value struct {
		ProtocolVersion string `json:"protocolVersion"`
	} `json:"__value__"`
}

type json_ClientConnect struct {
	Class string `json:"__class__"`
	Value struct {
		CommanderName string `json:"commanderName"`
		Language      string `json:"language"`
	} `json:"__value__"`
}
