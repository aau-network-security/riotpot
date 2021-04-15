package greeting

import (
	"testing"

	"github.com/riotpot/internal/greeting"
)

func TestExampleGreet(t *testing.T) {
	g := greeting.Greet{
		Tutorial: false,
	}

	g.Greeting()
}
