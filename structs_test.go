// This file is part of The AI Sandbox Go Bindings by errnoh.
// Copyright (c) 2012, errnoh@github
// License: See LICENSE file.
package aisandbox

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"
)

func BenchmarkSimplify(b *testing.B) {
	b.StopTimer()
	gi := new(json_GameInfo)
	err := json.Unmarshal([]byte(json_gameinfo), &gi)
	if err != nil {
		b.Fatalf(err.Error())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		gi.simplify()
	}
}

func TestJSONFlagInfo(t *testing.T) {
	s := new(json_FlagInfo)
	err := json.Unmarshal([]byte(json_flaginfo), &s)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestJSONBotInfo(t *testing.T) {
	s := new(json_BotInfo)
	err := json.Unmarshal([]byte(json_botinfo), &s)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestJSONMatchInfo(t *testing.T) {
	s := new(json_MatchInfo)
	err := json.Unmarshal([]byte(json_matchinfo), &s)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestJSONGameInfo(t *testing.T) {
	s := new(json_GameInfo)
	err := json.Unmarshal([]byte(json_gameinfo), &s)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestServer(t *testing.T) {
	var (
		buf     []byte
		bufConn *bufio.Reader
		err     error
	)

	go func() {
		ln, err := net.Listen("tcp", ":41187")
		if err != nil {
			t.Fatal(err.Error())
		}
		defer ln.Close()

		for {
			conn, err := ln.Accept()
			defer conn.Close()
			bufConn = bufio.NewReader(conn)
			if err != nil {
				t.Fatal(err.Error())
				continue
			}
			buf, err = bufConn.ReadBytes('\n')
			fmt.Printf("Joined: %s", buf)
			buf, err = bufConn.ReadBytes('\n')
			fmt.Printf("Command: %s", buf)
			conn.Write([]byte(json_init))
			conn.Write([]byte(json_tick))
			conn.Write([]byte(json_shutdown))
			return
		}
	}()

	go func() {
		<-time.After(time.Second * 3)
		t.Fatalf("Timeout after 3 seconds")
	}()

	in, out, err := Connect("Bacon", "localhost", 41187)
	if err != nil {
		// close(out)
		t.Fatal(err)
	}
	defer close(out)

	attack := new(Attack)
	attack.Class = "Attack"
	attack.Value.Bot = "red1"
	attack.Value.Description = "Go shoot stuff"
	attack.Value.LookAt = []float64{3, 5}
	attack.Value.Target = []float64{23, 93}
	out <- attack

	for v := range in {
		switch vtype := v.(type) {
		case *GameInfo:
			fmt.Println("GameInfo")
		case *LevelInfo:
			fmt.Println("LevelInfo")
		default:
			t.Fatalf("Unexpected type from 'in' channel: %s", vtype)
		}
	}
}

func TestJSON(t *testing.T) {
	// A bit ugly to test because of the anonymous structs. It's not a problem when actually using it though.
	expected_li := new(json_LevelInfo)
	expected_li.Class = "LevelInfo"
	expected_li.Value = new(LevelInfo)
	expected_li.Value.Width = 88
	expected_li.Value.Height = 50
	expected_li.Value.BlockHeights = [][]float64{{1, 2, 3}, {4, 5}, {6, 7, 8, 9}}
	expected_li.Value.TeamNames = []string{"Blue", "Red"}
	expected_li.Value.FlagSpawnLocations = map[string][]float64{"Blue": {82.0, 20.0}, "Red": {6.0, 30.0}}
	expected_li.Value.FlagScoreLocations = map[string][]float64{"Blue": {82.0, 20.0}, "Red": {6.0, 30.0}}
	expected_li.Value.BotSpawnAreas = map[string][][]float64{"Blue": {{79.0, 2.0}, {85.0, 9.0}}, "Red": {{3.0, 41.0}, {9.0, 48.0}}}
	expected_li.Value.FOVangle = 1.5707963267948966
	expected_li.Value.CharacterRadius = 0.25
	expected_li.Value.WalkingSpeed = 3.0
	expected_li.Value.RunningSpeed = 6.0
	expected_li.Value.FiringDistance = 15.0

	li := new(json_LevelInfo)
	err := json.Unmarshal([]byte(json_levelinfo), &li)
	if err != nil {
		t.Fatalf(err.Error())
	}

	test := [][]float64{
		{li.Value.BlockHeights[0][1], expected_li.Value.BlockHeights[0][1]},
		{li.Value.Width, expected_li.Value.Width},
		{li.Value.FOVangle, expected_li.Value.FOVangle},
		{li.Value.CharacterRadius, expected_li.Value.CharacterRadius},
		{li.Value.BotSpawnAreas["Blue"][0][1], expected_li.Value.BotSpawnAreas["Blue"][0][1]},
	}

	for _, pair := range test {
		if pair[0] != pair[1] {
			t.Errorf("JSON Unmarshal: Expected %d, got %d", pair[1], pair[0])
		}
	}
}

var (
	json_levelinfo = `{
  "__class__": "LevelInfo",
  "__value__": {
    "width": 88,
    "height": 50,
    "blockHeights": [[1, 2, 3], [4, 5], [6, 7, 8, 9]],           
    "teamNames": ["Blue", "Red"],                        
    "flagSpawnLocations": {                                          
      "Blue": [82.0, 20.0],
      "Red": [6.0, 30.0]
    },
    "flagScoreLocations": {                                               
      "Blue": [82.0, 20.0],
      "Red": [6.0, 30.0]
    },
    "botSpawnAreas": {                                              
      "Blue": [[79.0, 2.0], [85.0, 9.0]],
      "Red": [[3.0, 41.0], [9.0, 48.0]]
    },
    "FOVangle": 1.5707963267948966,
    "characterRadius": 0.25,
    "walkingSpeed": 3.0,
    "runningSpeed": 6.0,
    "firingDistance": 15.0
  }
}`

	json_init = `<initialize>
{"__class__": "LevelInfo", "__value__": {"runningSpeed": 6.0, "flagSpawnLocations": {"Blue": [82.0, 20.0], "Red": [6.0, 30.0]}, "teamNames": ["Blue", "Red"], "blockHeights": [[0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 2, 2, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 4, 4, 4, 4, 2, 2, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 2, 2, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 4, 4, 4, 4, 2, 2, 0, 0, 0, 0, 0, 0, 1, 2, 2, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 2, 2, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 4, 4, 4, 4, 1, 1, 0, 0, 0, 0, 0, 0, 1, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2], [0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2], [0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 1, 4, 4, 4, 4, 0, 0, 0, 0, 1, 1, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 2, 2, 2, 2, 2, 2, 0, 0, 0, 0, 0, 0], [0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 1, 4, 4, 4, 4, 0, 0, 0, 0, 0, 1, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 2, 2, 2, 2, 2, 2, 0, 0, 0, 0, 0, 0], [0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 1, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 0, 0, 0, 1, 1, 2, 2, 1, 1, 0, 0, 0, 0, 0, 0], [0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0], [0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0], [0, 0, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0], [0, 0, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 1, 4, 4, 4, 4, 1, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 1, 4, 4, 4, 4, 1, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0], [0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 1, 1, 2, 2, 1, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0], [0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0, 2, 2, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0], [0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 4, 4, 4, 4, 1, 0, 0, 0, 1, 2, 2, 1, 1, 1, 0], [0, 0, 0, 0, 2, 2, 4, 4, 4, 4, 0, 0, 0, 0, 0, 4, 4, 4, 4, 2, 2, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 4, 4, 4, 4, 1, 0, 0, 0, 1, 2, 2, 1, 2, 2, 0], [0, 0, 0, 0, 2, 2, 2, 2, 1, 1, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 4, 4, 4, 4, 1, 0, 0, 0, 0, 1, 1, 1, 2, 2, 0], [0, 0, 0, 0, 1, 1, 2, 2, 1, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 1, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 1, 2, 2, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 4, 4, 4, 4, 2, 2, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0, 1, 2, 2, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 4, 4, 4, 4, 2, 2, 1, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0], [1, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 4, 4, 4, 4, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 1, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 1, 2, 2, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0], [0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 4, 4, 4, 4, 1, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 4, 4, 4, 4, 1, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 4, 4, 4, 4, 1, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 4, 4, 4, 4, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 4, 4, 4, 4, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 4, 4, 4, 4, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0], [0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 1, 2, 2, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 1, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 4, 4, 4, 4, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 4, 4, 4, 4, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 1], [0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 1, 4, 4, 4, 4, 2, 2, 1, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 4, 4, 4, 4, 2, 2, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 2, 2, 0, 0, 0, 0], [0, 2, 2, 2, 2, 1, 0, 0, 0, 0, 1, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 4, 4, 4, 4, 2, 2, 0, 0, 0, 0], [0, 2, 2, 2, 2, 2, 2, 0, 0, 0, 1, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 1, 1, 1, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 1, 0, 0, 0, 0], [0, 1, 1, 1, 1, 2, 2, 0, 0, 0, 1, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 0, 0, 1, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0], [0, 0, 0, 1, 2, 2, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 1, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0], [0, 0, 0, 1, 2, 2, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 0, 0, 0, 2, 2, 1, 2, 2, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 1, 1, 0, 0, 0, 0, 0, 0], [0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 0, 0, 4, 4, 4, 4, 2, 2, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 4, 4, 4, 4, 1, 1, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 0, 0], [0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 0, 0], [0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 4, 4, 4, 4, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0], [0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 4, 4, 4, 4, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0], [0, 0, 0, 0, 0, 0, 1, 2, 2, 2, 2, 1, 0, 0, 0, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 4, 4, 4, 4, 1, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0], [0, 0, 0, 0, 0, 0, 1, 2, 2, 2, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0], [0, 0, 0, 0, 0, 0, 1, 2, 2, 1, 2, 2, 1, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 1, 0, 0, 0, 0, 1, 1, 2, 2, 1, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0], [2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 1, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0], [2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 0, 0, 0, 2, 2, 4, 4, 4, 4, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 2, 2, 0, 0, 0, 0, 0, 0, 2, 2, 4, 4, 4, 4, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 1, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 1, 0, 0, 0, 0, 0, 0, 1, 1, 4, 4, 4, 4, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 1, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]], "height": 50, "characterRadius": 0.25, "walkingSpeed": 3.0, "FOVangle": 1.5707963267948966, "botSpawnAreas": {"Blue": [[79.0, 2.0], [85.0, 9.0]], "Red": [[3.0, 41.0], [9.0, 48.0]]}, "firingDistance": 15.0, "width": 88, "flagScoreLocations": {"Blue": [82.0, 20.0], "Red": [6.0, 30.0]}}}
{"__class__": "GameInfo", "__value__": {"teams": {"Blue": {"__class__": "TeamInfo", "__value__": {"flagScoreLocation": [82.0, 20.0], "name": "Blue", "flagSpawnLocation": [82.0, 20.0], "flag": "BlueFlag", "members": ["Blue0", "Blue1", "Blue2", "Blue3", "Blue4"], "botSpawnArea": [[79.0, 2.0], [85.0, 9.0]]}}, "Red": {"__class__": "TeamInfo", "__value__": {"flagScoreLocation": [6.0, 30.0], "name": "Red", "flagSpawnLocation": [6.0, 30.0], "flag": "RedFlag", "members": ["Red0", "Red1", "Red2", "Red3", "Red4"], "botSpawnArea": [[3.0, 41.0], [9.0, 48.0]]}}}, "flags": {"BlueFlag": {"__class__": "FlagInfo", "__value__": {"position": [82.0, 20.0], "carrier": null, "name": "BlueFlag", "respawnTimer": 0.10000000149011612, "team": "Blue"}}, "RedFlag": {"__class__": "FlagInfo", "__value__": {"position": [6.0, 30.0], "carrier": null, "name": "RedFlag", "respawnTimer": 0.10000000149011612, "team": "Red"}}}, "enemyTeam": "Red", "team": "Blue", "bots": {"Red3": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Red3", "facingDirection": null, "state": 0, "health": 0.0, "seenlast": null, "team": "Red", "currentAction": null, "position": null, "visibleEnemies": []}}, "Red2": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Red2", "facingDirection": null, "state": 0, "health": 0.0, "seenlast": null, "team": "Red", "currentAction": null, "position": null, "visibleEnemies": []}}, "Red1": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Red1", "facingDirection": null, "state": 0, "health": 0.0, "seenlast": null, "team": "Red", "currentAction": null, "position": null, "visibleEnemies": []}}, "Red0": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Red0", "facingDirection": null, "state": 0, "health": 0.0, "seenlast": null, "team": "Red", "currentAction": null, "position": null, "visibleEnemies": []}}, "Red4": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Red4", "facingDirection": null, "state": 0, "health": 0.0, "seenlast": null, "team": "Red", "currentAction": null, "position": null, "visibleEnemies": []}}, "Blue1": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Blue1", "facingDirection": [0.06574580073356628, 0.9978364109992981], "state": 1, "health": 100.0, "seenlast": null, "team": "Blue", "currentAction": null, "position": [81.11000061035156, 6.492311954498291], "visibleEnemies": []}}, "Blue0": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Blue0", "facingDirection": [0.10403892397880554, 0.9945732355117798], "state": 1, "health": 100.0, "seenlast": null, "team": "Blue", "currentAction": null, "position": [80.45407104492188, 5.22149658203125], "visibleEnemies": []}}, "Blue3": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Blue3", "facingDirection": [0.22079943120479584, 0.9753192663192749], "state": 1, "health": 100.0, "seenlast": null, "team": "Blue", "currentAction": null, "position": [79.2674331665039, 7.929657459259033], "visibleEnemies": []}}, "Blue2": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Blue2", "facingDirection": [0.0015204440569505095, 0.9999988675117493], "state": 1, "health": 100.0, "seenlast": null, "team": "Blue", "currentAction": null, "position": [81.97264862060547, 2.010946273803711], "visibleEnemies": []}}, "Blue4": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Blue4", "facingDirection": [0.2348455935716629, 0.9720326662063599], "state": 1, "health": 100.0, "seenlast": null, "team": "Blue", "currentAction": null, "position": [79.1587905883789, 8.240152359008789], "visibleEnemies": []}}}, "match": {"__class__": "MatchInfo", "__value__": {"timeRemaining": 180.0, "timeToNextRespawn": 45.0, "combatEvents": [], "timePassed": 0.0, "scores": {"Blue": 0, "Red": 0}}}}}
`

	json_tick = `<tick>
{"__class__": "GameInfo", "__value__": {"teams": {"Blue": {"__class__": "TeamInfo", "__value__": {"flagScoreLocation": [82.0, 20.0], "name": "Blue", "flagSpawnLocation": [82.0, 20.0], "flag": "BlueFlag", "members": ["Blue0", "Blue1", "Blue2", "Blue3", "Blue4"], "botSpawnArea": [[79.0, 2.0], [85.0, 9.0]]}}, "Red": {"__class__": "TeamInfo", "__value__": {"flagScoreLocation": [6.0, 30.0], "name": "Red", "flagSpawnLocation": [6.0, 30.0], "flag": "RedFlag", "members": ["Red0", "Red1", "Red2", "Red3", "Red4"], "botSpawnArea": [[3.0, 41.0], [9.0, 48.0]]}}}, "flags": {"BlueFlag": {"__class__": "FlagInfo", "__value__": {"position": [82.0, 20.0], "carrier": null, "name": "BlueFlag", "respawnTimer": -7.450580596923828e-09, "team": "Blue"}}, "RedFlag": {"__class__": "FlagInfo", "__value__": {"position": [9.723822593688965, 28.638526916503906], "carrier": "Blue1", "name": "RedFlag", "respawnTimer": -7.450580596923828e-09, "team": "Red"}}}, "enemyTeam": "Red", "team": "Blue", "bots": {"Red3": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Red3", "facingDirection": [0.9375345706939697, -0.3478919267654419], "state": 6, "health": 0, "seenlast": 13.370665550231934, "team": "Red", "currentAction": "ShootAtCommand", "position": [35.6309928894043, 26.81215476989746], "visibleEnemies": []}}, "Red2": {"__class__": "BotInfo", "__value__": {"seenBy": ["Blue0"], "flag": null, "name": "Red2", "facingDirection": [0.9123391509056091, -0.4094350337982178], "state": 6, "health": 0, "seenlast": 0.0, "team": "Red", "currentAction": "ShootAtCommand", "position": [68.28890991210938, 25.360763549804688], "visibleEnemies": []}}, "Red1": {"__class__": "BotInfo", "__value__": {"seenBy": ["Blue0"], "flag": null, "name": "Red1", "facingDirection": [-0.9972056150436401, 0.07470673322677612], "state": 4, "health": 0, "seenlast": 0.0, "team": "Red", "currentAction": "AttackCommand", "position": [68.53483581542969, 25.27260398864746], "visibleEnemies": []}}, "Red0": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Red0", "facingDirection": [0.9994280338287354, -0.033820152282714844], "state": 6, "health": 0, "seenlast": 13.370665550231934, "team": "Red", "currentAction": "ShootAtCommand", "position": [34.46906280517578, 24.155515670776367], "visibleEnemies": []}}, "Red4": {"__class__": "BotInfo", "__value__": {"seenBy": ["Blue0"], "flag": null, "name": "Red4", "facingDirection": [0.912505030632019, -0.4090656042098999], "state": 6, "health": 0, "seenlast": 0.0, "team": "Red", "currentAction": "ShootAtCommand", "position": [68.30572509765625, 25.36515998840332], "visibleEnemies": []}}, "Blue1": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": "RedFlag", "name": "Blue1", "facingDirection": [0.9242773652076721, -0.3817223310470581], "state": 3, "health": 100.0, "seenlast": null, "team": "Blue", "currentAction": "MoveCommand", "position": [9.723822593688965, 28.638526916503906], "visibleEnemies": []}}, "Blue0": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Blue0", "facingDirection": [-0.9890086054801941, 0.14785832166671753], "state": 1, "health": 100.0, "seenlast": null, "team": "Blue", "currentAction": null, "position": [81.625, 19.375], "visibleEnemies": ["Red2", "Red1", "Red4"]}}, "Blue3": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Blue3", "facingDirection": [-0.9994280338287354, 0.03381979465484619], "state": 1, "health": 0, "seenlast": null, "team": "Blue", "currentAction": null, "position": [48.790069580078125, 23.665205001831055], "visibleEnemies": []}}, "Blue2": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Blue2", "facingDirection": [-0.9112738966941833, 0.411800742149353], "state": 6, "health": 0, "seenlast": null, "team": "Blue", "currentAction": "ShootAtCommand", "position": [57.94633102416992, 32.63374710083008], "visibleEnemies": []}}, "Blue4": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Blue4", "facingDirection": [-0.9575538635253906, 0.2882544994354248], "state": 6, "health": 0, "seenlast": null, "team": "Blue", "currentAction": "ShootAtCommand", "position": [47.545501708984375, 19.977867126464844], "visibleEnemies": []}}}, "match": {"__class__": "MatchInfo", "__value__": {"timeRemaining": 148.42462158203125, "timeToNextRespawn": 13.427755355834961, "combatEvents": [{"__class__": "MatchCombatEvent", "__value__": {"instigator": "Blue3", "time": 14.939663887023926, "type": 1, "subject": "Red3"}}, {"__class__": "MatchCombatEvent", "__value__": {"instigator": "Red2", "time": 16.550338745117188, "type": 1, "subject": "Blue2"}}, {"__class__": "MatchCombatEvent", "__value__": {"instigator": "Red4", "time": 16.550338745117188, "type": 1, "subject": "Blue2"}}, {"__class__": "MatchCombatEvent", "__value__": {"instigator": "Red0", "time": 17.310344696044922, "type": 1, "subject": "Blue4"}}, {"__class__": "MatchCombatEvent", "__value__": {"instigator": "Blue3", "time": 18.036685943603516, "type": 1, "subject": "Red0"}}, {"__class__": "MatchCombatEvent", "__value__": {"instigator": "Red1", "time": 18.201021194458008, "type": 1, "subject": "Blue3"}}, {"__class__": "MatchCombatEvent", "__value__": {"instigator": "Blue0", "time": 28.15752601623535, "type": 1, "subject": "Red4"}}, {"__class__": "MatchCombatEvent", "__value__": {"instigator": "Blue1", "time": 28.15752601623535, "type": 2, "subject": "RedFlag"}}, {"__class__": "MatchCombatEvent", "__value__": {"instigator": "Blue0", "time": 28.616199493408203, "type": 1, "subject": "Red2"}}, {"__class__": "MatchCombatEvent", "__value__": {"instigator": "Blue0", "time": 29.308876037597656, "type": 1, "subject": "Red1"}}], "timePassed": 31.5719051361084, "scores": {"Blue": 0, "Red": 0}}}}}
`
	json_shutdown = `<shutdown>
 `

	json_flaginfo = `{
                "__class__": "FlagInfo",
                "__value__": {
                    "position": [
                        82,
                        20
                    ],
                    "carrier": null,
                    "name": "BlueFlag",
                    "respawnTimer": -7.450580596923828e-9,
                    "team": "Blue"
                }
            }`

	json_botinfo = `{
                "__class__": "BotInfo",
                "__value__": {
                    "seenBy": [],
                    "flag": null,
                    "name": "Red3",
                    "facingDirection": [
                        0.9375345706939697,
                        -0.3478919267654419
                    ],
                    "state": 6,
                    "health": 0,
                    "seenlast": 13.370665550231934,
                    "team": "Red",
                    "currentAction": "ShootAtCommand",
                    "position": [
                        35.6309928894043,
                        26.81215476989746
                    ],
                    "visibleEnemies": []
                }
            }`

	json_matchinfo = `{
            "__class__": "MatchInfo",
            "__value__": {
                "timeRemaining": 148.42462158203125,
                "timeToNextRespawn": 13.427755355834961,
                "combatEvents": [
                    {
                        "__class__": "MatchCombatEvent",
                        "__value__": {
                            "instigator": "Blue3",
                            "time": 14.939663887023926,
                            "type": 1,
                            "subject": "Red3"
                        }
                    },
                    {
                        "__class__": "MatchCombatEvent",
                        "__value__": {
                            "instigator": "Red2",
                            "time": 16.550338745117188,
                            "type": 1,
                            "subject": "Blue2"
                        }
                    },
                    {
                        "__class__": "MatchCombatEvent",
                        "__value__": {
                            "instigator": "Blue0",
                            "time": 29.308876037597656,
                            "type": 1,
                            "subject": "Red1"
                        }
                    }
                ],
                "timePassed": 31.5719051361084,
                "scores": {
                    "Blue": 0,
                    "Red": 0
                }
            }
        }`

	json_gameinfo = `{"__class__": "GameInfo", "__value__": {"teams": {"Blue": {"__class__": "TeamInfo", "__value__": {"flagScoreLocation": [82.0, 20.0], "name": "Blue", "flagSpawnLocation": [82.0, 20.0], "flag": "BlueFlag", "members": ["Blue0", "Blue1", "Blue2", "Blue3", "Blue4"], "botSpawnArea": [[79.0, 2.0], [85.0, 9.0]]}}, "Red": {"__class__": "TeamInfo", "__value__": {"flagScoreLocation": [6.0, 30.0], "name": "Red", "flagSpawnLocation": [6.0, 30.0], "flag": "RedFlag", "members": ["Red0", "Red1", "Red2", "Red3", "Red4"], "botSpawnArea": [[3.0, 41.0], [9.0, 48.0]]}}}, "flags": {"BlueFlag": {"__class__": "FlagInfo", "__value__": {"position": [82.0, 20.0], "carrier": null, "name": "BlueFlag", "respawnTimer": -7.450580596923828e-09, "team": "Blue"}}, "RedFlag": {"__class__": "FlagInfo", "__value__": {"position": [9.723822593688965, 28.638526916503906], "carrier": "Blue1", "name": "RedFlag", "respawnTimer": -7.450580596923828e-09, "team": "Red"}}}, "enemyTeam": "Red", "team": "Blue", "bots": {"Red3": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Red3", "facingDirection": [0.9375345706939697, -0.3478919267654419], "state": 6, "health": 0, "seenlast": 13.370665550231934, "team": "Red", "currentAction": "ShootAtCommand", "position": [35.6309928894043, 26.81215476989746], "visibleEnemies": []}}, "Red2": {"__class__": "BotInfo", "__value__": {"seenBy": ["Blue0"], "flag": null, "name": "Red2", "facingDirection": [0.9123391509056091, -0.4094350337982178], "state": 6, "health": 0, "seenlast": 0.0, "team": "Red", "currentAction": "ShootAtCommand", "position": [68.28890991210938, 25.360763549804688], "visibleEnemies": []}}, "Red1": {"__class__": "BotInfo", "__value__": {"seenBy": ["Blue0"], "flag": null, "name": "Red1", "facingDirection": [-0.9972056150436401, 0.07470673322677612], "state": 4, "health": 0, "seenlast": 0.0, "team": "Red", "currentAction": "AttackCommand", "position": [68.53483581542969, 25.27260398864746], "visibleEnemies": []}}, "Red0": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Red0", "facingDirection": [0.9994280338287354, -0.033820152282714844], "state": 6, "health": 0, "seenlast": 13.370665550231934, "team": "Red", "currentAction": "ShootAtCommand", "position": [34.46906280517578, 24.155515670776367], "visibleEnemies": []}}, "Red4": {"__class__": "BotInfo", "__value__": {"seenBy": ["Blue0"], "flag": null, "name": "Red4", "facingDirection": [0.912505030632019, -0.4090656042098999], "state": 6, "health": 0, "seenlast": 0.0, "team": "Red", "currentAction": "ShootAtCommand", "position": [68.30572509765625, 25.36515998840332], "visibleEnemies": []}}, "Blue1": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": "RedFlag", "name": "Blue1", "facingDirection": [0.9242773652076721, -0.3817223310470581], "state": 3, "health": 100.0, "seenlast": null, "team": "Blue", "currentAction": "MoveCommand", "position": [9.723822593688965, 28.638526916503906], "visibleEnemies": []}}, "Blue0": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Blue0", "facingDirection": [-0.9890086054801941, 0.14785832166671753], "state": 1, "health": 100.0, "seenlast": null, "team": "Blue", "currentAction": null, "position": [81.625, 19.375], "visibleEnemies": ["Red2", "Red1", "Red4"]}}, "Blue3": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Blue3", "facingDirection": [-0.9994280338287354, 0.03381979465484619], "state": 1, "health": 0, "seenlast": null, "team": "Blue", "currentAction": null, "position": [48.790069580078125, 23.665205001831055], "visibleEnemies": []}}, "Blue2": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Blue2", "facingDirection": [-0.9112738966941833, 0.411800742149353], "state": 6, "health": 0, "seenlast": null, "team": "Blue", "currentAction": "ShootAtCommand", "position": [57.94633102416992, 32.63374710083008], "visibleEnemies": []}}, "Blue4": {"__class__": "BotInfo", "__value__": {"seenBy": [], "flag": null, "name": "Blue4", "facingDirection": [-0.9575538635253906, 0.2882544994354248], "state": 6, "health": 0, "seenlast": null, "team": "Blue", "currentAction": "ShootAtCommand", "position": [47.545501708984375, 19.977867126464844], "visibleEnemies": []}}}, "match": {"__class__": "MatchInfo", "__value__": {"timeRemaining": 148.42462158203125, "timeToNextRespawn": 13.427755355834961, "combatEvents": [{"__class__": "MatchCombatEvent", "__value__": {"instigator": "Blue3", "time": 14.939663887023926, "type": 1, "subject": "Red3"}}, {"__class__": "MatchCombatEvent", "__value__": {"instigator": "Red2", "time": 16.550338745117188, "type": 1, "subject": "Blue2"}}, {"__class__": "MatchCombatEvent", "__value__": {"instigator": "Red4", "time": 16.550338745117188, "type": 1, "subject": "Blue2"}}, {"__class__": "MatchCombatEvent", "__value__": {"instigator": "Red0", "time": 17.310344696044922, "type": 1, "subject": "Blue4"}}, {"__class__": "MatchCombatEvent", "__value__": {"instigator": "Blue3", "time": 18.036685943603516, "type": 1, "subject": "Red0"}}, {"__class__": "MatchCombatEvent", "__value__": {"instigator": "Red1", "time": 18.201021194458008, "type": 1, "subject": "Blue3"}}, {"__class__": "MatchCombatEvent", "__value__": {"instigator": "Blue0", "time": 28.15752601623535, "type": 1, "subject": "Red4"}}, {"__class__": "MatchCombatEvent", "__value__": {"instigator": "Blue1", "time": 28.15752601623535, "type": 2, "subject": "RedFlag"}}, {"__class__": "MatchCombatEvent", "__value__": {"instigator": "Blue0", "time": 28.616199493408203, "type": 1, "subject": "Red2"}}, {"__class__": "MatchCombatEvent", "__value__": {"instigator": "Blue0", "time": 29.308876037597656, "type": 1, "subject": "Red1"}}], "timePassed": 31.5719051361084, "scores": {"Blue": 0, "Red": 0}}}}}`
)
