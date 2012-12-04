package main

import (
	"fmt"
	"github.com/errnoh/aisandbox"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var (
		width, height float64
		initialized   bool
		err           error

		host = "localhost"
		port = 41041
	)

	// Competition server needs to be able to use custom host and port.
	// These should be given as first and second command-line arguments when launching the bot.
	if len(os.Args) >= 3 {
		host = os.Args[1]
		port, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid port number: %s", os.Args[2])
		}
	}

	in, out, err := aisandbox.Connect(host, port, "GoRandom")
	if err != nil {
		log.Fatalln(err)
	}

	for msg := range in {
		switch m := msg.(type) {
		// When LevelInfo is received..
		case *aisandbox.LevelInfo:
			// Process LevelInfo data
			log.Println("Level loaded")
			width, height = m.Width, m.Height
			log.Println(width, height)
		// And the main game status updates
		case *aisandbox.GameInfo:
			var target []float64
			var text string
			var team *aisandbox.TeamInfo

			// After getting the first GameInfo packet, process it if you need and send server information that you have processed the data
			if !initialized {
				// Process initial GameInfo data here
				// doStuff()
				initialized = true

				// Inform server that you're ready.
				aisandbox.Ready()
				// Don't send commands yet, wait for first actual game tick.
				continue
			}

			for _, bot := range m.Team.Members {
				// Skip dead bots and bots who already have something to do.
				if bot.Health == 0 || bot.State > 1 {
					continue
				}

				// Throw in couple random numbers
				r, r2 := rand.Intn(3), rand.Intn(2)
				// Select either own or enemy team as target
				if r2 == 0 {
					team = m.Team
				} else {
					team = m.EnemyTeam
				}

				switch r {
				case 0:
					// Attack current flag position of target team
					target = team.Flag.Position
					text = fmt.Sprintf("Attacking %s flag.", team.Name)

				case 1:
					// Attack spawn location of target teams flag
					target = team.FlagScoreLocation
					text = fmt.Sprintf("Attacking %s score location.", team.Name)
				case 2:
					// Attack random point on the map
					// XXX: Doesn't check if target is possible.
					target = []float64{rand.Float64() * width, rand.Float64() * height}
					text = fmt.Sprintf("Attacking [%.2f, %.2f].", target[0], target[1])
				}

				r = rand.Intn(4)
				switch r {
				case 0:
					out <- aisandbox.NewMove(bot.Name, text, target)
				case 1:
					out <- aisandbox.NewAttack(bot.Name, text, nil, target)
				case 2:
					out <- aisandbox.NewDefend(bot.Name, "Defending", []float64{1, 1}, []float64{1, -1, 1}, []float64{-1, -1, 2}, []float64{-1, 1, 0})
				case 3:
					out <- aisandbox.NewCharge(bot.Name, text, target)
				}
			}
		}
	}
	fmt.Println("Received <shutdown> from server")
	close(out)
}
