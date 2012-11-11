package main

import (
	"fmt"
	"github.com/errnoh/aisandbox"
	"log"
	"math/rand"
)

func main() {
	var (
		width, height float64
	)

	in, out, err := aisandbox.Connect("GoRandom", "localhost", 5557)
	if err != nil {
		log.Fatalln(err)
	}

	for msg := range in {
		switch m := msg.(type) {
		// When LevelInfo is received..
		case *aisandbox.LevelInfo:
			log.Println("Level loaded")
			width, height = m.Value.Width, m.Value.Height
			log.Println(width, height)
		// And the main game status updates
		case *aisandbox.GameInfo:
			var target []float64
			var text string

			// Select one bot from own team (doesn't check if the bot is dead :)
			bot := rand.Intn(len(m.Value.Teams[m.Value.Team].Value.Members))

			// Throw in couple random numbers
			r, r2 := rand.Intn(3), rand.Intn(2)
			// Select either own or enemy team as target
			if r2 == 0 {
				text = m.Value.Team
			} else {
				text = m.Value.EnemyTeam
			}
			switch r {
			case 0:
				// Attack current flag position of target team
				target = m.Value.Flags[text+"Flag"].Value.Position
				text = fmt.Sprintf("%s flag.", text)

			case 1:
				// Attack spawn location of target teams flag
				target = m.Value.Teams[text].Value.FlagScoreLocation
				text = fmt.Sprintf("%s score location.", text)
			case 2:
				// Attack random point on the map
				// XXX: Doesn't check if target is possible.
				target = []float64{rand.Float64() * width, rand.Float64() * height}
				text = fmt.Sprintf("[%.2f, %.2f].", target[0], target[1])
			}
			out <- attack(m.Value.Teams[m.Value.Team].Value.Members[bot], target, text)
		}
	}

}

func attack(name string, coords []float64, description string) *aisandbox.Attack {
	a := new(aisandbox.Attack)
	a.Class = "Attack"
	a.Value.Bot = name
	a.Value.Target = coords
	a.Value.LookAt = coords
	a.Value.Description = fmt.Sprintf("Attacking %s", description)
	return a
}
