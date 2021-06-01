package cli

import (
	"os"
	"testing"

	"github.com/riotpot/internal/cli"
)

var (
	// TODO add this file
	path = "./greeting_test.yml"
)

func TestYBookParse(t *testing.T) {
	book := cli.YBook{}
	term := cli.Terminal{}

	reader, err := os.ReadFile(path)

	if err != nil {
		t.Error(err)
	}

	book.Parse(reader)

	for _, page := range book.Pages {
		for _, line := range page.Lines {
			term.Read(line)
		}
	}
}

func TestLecture(t *testing.T) {
	term := cli.Terminal{}
	term.Lecture(path)
}
