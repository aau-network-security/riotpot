/*
This package implements the `hello` message when riotpot is launch.
*/
package greeting

import (
	"fmt"

	"github.com/riotpot/internal/cli"
)

type Greet struct {
	// Indicates if this is the first time riotpot has been launched
	Tutorial bool
	// Indicates if we want the initial message at all
	Initial bool
}

func (g *Greet) Greeting() {
	// regardless, show the initial welcome
	g.initial(g.Initial)

	switch {
	case g.Tutorial: // Throw a walkthrough and initial greeting
		// just in case, set the welcome now to false, so the user does not
		// get the
		g.Tutorial = false
		g.walkthrough()
	}
}

// Gives a walkthrough RiotPot
func (g *Greet) walkthrough() {}

// Throws a regular salute
func (g *Greet) initial(show bool) {
	if show {
		term := cli.Terminal{}
		term.Lecture("configs/greeting/hello.yml")
	} else {
		fmt.Println("░▒▓███ RIoIPot ███▓▒░")
	}
}
