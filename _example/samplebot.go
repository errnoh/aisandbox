package main

import (
	"fmt"
	"github.com/errnoh/aisandbox"
	"log"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var (
		width, height float64
	)

	in, out, err := aisandbox.Connect("localhost", 5557, "GoRandom")
	if err != nil {
		log.Fatalln(err)
	}

	for msg := range in {
		switch m := msg.(type) {
		// When LevelInfo is received..
		case *aisandbox.LevelInfo:
			log.Println("Level loaded")
			width, height = m.Width, m.Height
			log.Println(width, height)
		// And the main game status updates
		case *aisandbox.GameInfo:
			var target []float64
			var text string
			var team *aisandbox.TeamInfo

			for _, bot := range m.Team.Members {
				if bot.Health == 0 {
					continue
				}

				// Only update command 1/10 of the time.
				if r := rand.Float64(); r > 0.1 {
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
					text = fmt.Sprintf("%s flag.", team.Name)

				case 1:
					// Attack spawn location of target teams flag
					target = team.FlagScoreLocation
					text = fmt.Sprintf("%s score location.", team.Name)
				case 2:
					// Attack random point on the map
					// XXX: Doesn't check if target is possible.
					target = []float64{rand.Float64() * width, rand.Float64() * height}
					text = fmt.Sprintf("[%.2f, %.2f].", target[0], target[1])
				}
				out <- attack(bot.Name, target, text)
			}
		}
	}

}

func attack(name string, coords []float64, description string) *aisandbox.Attack {
	return &aisandbox.Attack{
		Bot:         name,
		Target:      coords,
		LookAt:      coords,
		Description: fmt.Sprintf("Attacking %s", description),
	}
}
