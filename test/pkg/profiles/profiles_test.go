package profiles

import (
	"testing"

	"github.com/riotpot/pkg/profiles"
)

func TestRandomPaths(t *testing.T) {
	topics := profiles.RandomNumericTopics("/ps", 10)

	for _, top := range topics {
		t.Log(top)
	}
}

func TestLoadProfile(t *testing.T) {
	p := profiles.Profile{}
	p.Load("profile_test.yml")
	t.Log(p)
}
