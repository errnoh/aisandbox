// This file is part of The AI Sandbox Go Bindings by errnoh.
// Copyright (c) 2012, errnoh@github
// License: See LICENSE file.

package aisandbox

type LevelInfo struct {
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
}

type GameInfo struct {
	Team      *TeamInfo
	EnemyTeam *TeamInfo
	Match     *MatchInfo
}

type TeamInfo struct {
	Name              string
	Flag              *FlagInfo
	Members           map[string]*BotInfo
	FlagSpawnLocation []float64
	FlagScoreLocation []float64
	BotSpawnArea      [][]float64
}

type FlagInfo struct {
	Position     []float64
	Carrier      *BotInfo
	RespawnTimer float64
}

type BotInfo struct {
	Name            string
	Team            string
	Position        []float64
	FacingDirection []float64
	Flag            string
	CurrentAction   string
	State           float64 // values are 0 = unknown, 1 = idle, 2 = defending, 3 = moving, 4 = attacking, 5 = charging, 6 = shooting
	Health          float64
	SeenLast        float64
	VisibleEnemies  []*BotInfo
	SeenBy          []*BotInfo
}

type MatchInfo struct {
	TimeRemaining     float64
	TimeToNextRespawn float64
	CombatEvents      []*CombatEvent
}

type CombatEvent struct {
	Type       float64
	Instigator string
	Subject    string // can either be a FlagInfo or a BotInfo name
	Time       float64
}

// Commands

type Command interface {
	JSON() []byte
}

type Defend struct {
	Bot             string    `json:"bot"`
	FacingDirection []float64 `json:"facingDirection"`
	Description     string    `json:"description"`
}

func (c *Defend) JSON() []byte {
	cmd := struct {
		Class string  `json:"__class__"`
		Value *Defend `json:"__value__"`
	}{
		Class: "Defend",
		Value: c,
	}

	return marshal(cmd)
}

type Move struct {
	Bot         string    `json:"bot"`
	Target      []float64 `json:"target"`
	Description string    `json:"description"`
}

func (c *Move) JSON() []byte {
	cmd := struct {
		Class string `json:"__class__"`
		Value *Move  `json:"__value__"`
	}{
		Class: "Move",
		Value: c,
	}

	return marshal(cmd)
}

type Attack struct {
	Bot         string    `json:"bot"`
	Target      []float64 `json:"target"`
	LookAt      []float64 `json:"lookAt, omitempty"` // Optional
	Description string    `json:"description"`
}

func (c *Attack) JSON() []byte {
	cmd := struct {
		Class string  `json:"__class__"`
		Value *Attack `json:"__value__"`
	}{
		Class: "Attack",
		Value: c,
	}

	return marshal(cmd)
}

type Charge struct {
	Bot         string    `json:"bot"`
	Target      []float64 `json:"target"`
	Description string    `json:"description"`
}

func (c *Charge) JSON() []byte {
	cmd := struct {
		Class string  `json:"__class__"`
		Value *Charge `json:"__value__"`
	}{
		Class: "Charge",
		Value: c,
	}

	return marshal(cmd)
}
