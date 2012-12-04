// This file is part of The AI Sandbox Go Bindings by errnoh.
// Copyright (c) 2012, errnoh@github
// License: See LICENSE file.

package aisandbox

import (
	"errors"
	"fmt"
)

// NOTE: This file contains exported structs

// Exported structs that contain the server messages

type LevelInfo struct {
	Width              float64                `json:"width"`
	Height             float64                `json:"height"`
	BlockHeights       [][]float64            `json:"blockHeights"`       // a 'width' list of 'height' lengthed list of integers (BlockHeights[x][y])
	TeamNames          []string               `json:"teamNames"`          // list of team names
	FlagSpawnLocations map[string][]float64   `json:"flagSpawnLocations"` // map of team name to position
	FlagScoreLocations map[string][]float64   `json:"flagScoreLocations"` // map of team name to position
	BotSpawnAreas      map[string][][]float64 `json:"botSpawnAreas"`      // map of team name to min and max positions
	FOVangle           float64                `json:"FOVangle"`
	CharacterRadius    float64                `json:"characterRadius"`
	WalkingSpeed       float64                `json:"walkingSpeed"`
	RunningSpeed       float64                `json:"runningSpeed"`
	FiringDistance     float64                `json:"firingDistance"`
	GameLength         float64                `json:"gameLength"`         // the time (seconds) that a game will last
	InitializationTime float64                `json:"initializationTime"` // the time (seconds) allowed to the commanders for initialization
	RespawnTime        float64                `json:"respawnTime"`
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
	Score             float64
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
	State           float64 // values are 0 = unknown, 1 = idle, 2 = defending, 3 = moving, 4 = attacking, 5 = charging, 6 = shooting
	Health          float64
	SeenLast        float64
	VisibleEnemies  []*BotInfo
	SeenBy          []*BotInfo
}

type MatchInfo struct {
	TimeRemaining     float64
	TimeToNextRespawn float64
	TimePassed        float64
	CombatEvents      []*CombatEvent
}

type CombatEvent struct {
	Type       float64 // values are 0 = none, 1 = bot killed, 2 = flag picked up, 3 = flag dropped 4 = flag captured, 5 = flag restored, 6 = bot respawned
	Instigator string
	Subject    string // can either be a FlagInfo or a BotInfo name
	Time       float64
}

// Command structs for the commander

type Command interface {
	JSON() []byte
}

// Since update 1.4 Defend can be passed as many [direction], duration pairs as one wants.
// That's not really a slice of any actual type, so Defend became a bit more complicated.
// Because of this, constructors are added to the API.
//
// If you still want to do manual constructor, there's a sample below.

/*
	Example defend:
	&Defend{
		Bot:		"Bacon",
		Description:	"Mmmm",
		FacingDirections: []FacingDirection{
			FacingDirection{[]float64{1.2, 2}, 3.7},
			FacingDirection{[]float64{4, 5.1}, 6.3},
		},
	}
*/
type Defend struct {
	Bot              string             `json:"bot"`
	FacingDirections []*FacingDirection `json:"facingDirections"`
	Description      string             `json:"description"`
}

type FacingDirection struct {
	Direction []float64
	Duration  float64
}

func (fd FacingDirection) MarshalJSON() (b []byte, err error) {
	if len(fd.Direction) != 2 {
		return nil, errors.New("Invalid coordinates in FacingDirection")
	}
	return []byte(fmt.Sprintf("[[%f, %f], %f]", fd.Direction[0], fd.Direction[1], fd.Duration)), nil

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
	Bot         string      `json:"bot"`
	Target      [][]float64 `json:"target"`
	Description string      `json:"description"`
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
	Bot         string      `json:"bot"`
	Target      [][]float64 `json:"target"`
	LookAt      []float64   `json:"lookAt, omitempty"` // Optional
	Description string      `json:"description"`
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
	Bot         string      `json:"bot"`
	Target      [][]float64 `json:"target"`
	Description string      `json:"description"`
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
