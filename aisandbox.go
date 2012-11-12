// This file is part of The AI Sandbox Go Bindings by errnoh.
// Copyright (c) 2012, errnoh@github
// License: See LICENSE file.
package aisandbox

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
)

var (
	conn net.Conn
	// Buffer the connection so we can read it line by line
	bufConn *bufio.Reader
)

const (
	STATE_UNKNOWN = iota
	STATE_IDLE
	STATE_DEFENDING
	STATE_MOVING
	STATE_ATTACKING
	STATE_CHARGING
	STATE_SHOOTING
)

// Runs in the background and listens for messages from conn
// then parsing the JSON data into structs and forwarding those
// to the commander through a channel.
//
// NOTE: Possibly add:
// "Control" struct to inform the commander for shutdown etc.
// "Initialize" struct that holds both LevelInfo and GameInfo struct and is only sent on <Initialize>
func run(name string, c chan interface{}) {
	var (
		err         error
		buffer      []byte
		message     string
		gameinfo    *json_GameInfo
		levelinfo   *json_LevelInfo
		initialized bool
	)
	// Register with the server
	conn.Write([]byte(name))
loop:
	for {
		if buffer, err = bufConn.ReadBytes('\n'); err != nil {
			log.Println(err)
			break
		}
		message = strings.TrimSpace(string(buffer))
		switch message {
		case "":
			log.Println("Empty line from conn")
			continue loop
		case "<initialize>":
			if initialized {
				log.Printf("Unexpected initialize message '%s'", message)
			}
			// Read level info
			levelinfo = new(json_LevelInfo)
			if err = parsejson(&levelinfo); err != nil {
				log.Println(err)
				continue
			}
			// Read game info
			gameinfo = new(json_GameInfo)
			if err = parsejson(&gameinfo); err != nil {
				log.Println(err)
				continue
			}
			// Only send the actual data of LevelInfo
			c <- levelinfo.Value
			// Make GameInfo more intuitive to use
			c <- gameinfo.simplify()
			initialized = true
		case "<tick>":
			if !initialized {
				log.Printf("Unexpected message '%s' while waiting for initialize", message)
			}
			gameinfo = new(json_GameInfo)
			if err = parsejson(&gameinfo); err != nil {
				log.Println(err)
				continue
			}
			c <- gameinfo.simplify()
		case "<shutdown>":
			if !initialized {
				log.Printf("Unexpected message '%s' while waiting for initialize", message)
			}
			break loop
		default:
			log.Printf("unknown message received: '%s'", message)
		}
	}
	// Tell the commander that we're done here.
	close(c)
}

// Runs in the background and listens to the channel for commands sent by the commander.
func listen(c chan Command) {
	for v := range c {
		// NOTE: Possibly add a check that the message is one of the accepted commands.
		conn.Write([]byte("<command>\n"))
		conn.Write(v.JSON())
	}

	conn.Close()
}

// NOTE: AiSandbox spec says that messages can't contain newlines.
func parsejson(target interface{}) (err error) {
	var buffer []byte

	if buffer, err = bufConn.ReadBytes('\n'); err != nil {
		log.Println(err)
		return
	}
	json.Unmarshal(buffer, target)
	return
}

// Trims newlines and adds one newline to the end.
func trim(b []byte) []byte {
	var count int
	for i := 0; i < len(b); i++ {
		switch b[i] {
		case '\n', '\r':
		default:
			b[count] = b[i]
			count++
		}
	}
	b = b[:count]
	b = append(b, '\n')
	return b
}

func marshal(data interface{}) []byte {
	var (
		buffer []byte
		err    error
	)

	if buffer, err = json.Marshal(data); err != nil {
		log.Println(err)
	}
	return trim(buffer)
}

// Opens a connection to the server
// NOTE: In case of shutdown the "in" -channel will be closed to inform commander that it should shut down.
// NOTE: Close "out" channel when finished to close the connection to the server.
// Params:
// host, port - address and port of the server
// name - name of the commander
// Returns:
// in - incoming updates, being either LevelInfo or GameInfo structs (possibly add control struct to inform about Shutdown etc)
// out - outgoing channel where commander can send his commands, preferably Defend, Attack, Move or Charge structs.
func Connect(host string, port int, name string) (in <-chan interface{}, out chan<- Command, err error) {
	conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Printf("Failed to connect to the server: %s", err.Error())
		return
	}
	bufConn = bufio.NewReader(conn)

	i, o := make(chan interface{}), make(chan Command)

	in = i
	out = o
	go run(name, i)
	go listen(o)
	return
}
