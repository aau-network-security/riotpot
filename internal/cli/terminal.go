package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/riotpot/tools/errors"

	tm "github.com/buger/goterm"
	yaml "gopkg.in/yaml.v3"
)

// Line or paragraph from a yml file
type YLine struct {
	// The string content of the paragraph
	Str string `yaml:"str"`
	// how to align the content in the terminal vertically and horizontally
	H_align   string `yaml:"h_align"`
	V_align   string `yaml:"v_align"`
	Animation string `yaml:"animation"`
	Input     bool   `yaml:"input"`
}

// Individual document page from a yml file
type YPage struct {
	Lines []YLine `yaml:"lines,flow"`
}

// Implements a wrapper for yml multi-document splits.
type YBook struct {
	Pages []YPage `yaml:",flow"`
}

func (y *YBook) Parse(source []byte) (err error) {
	dec := yaml.NewDecoder(bytes.NewReader(source))
	for {
		var doc YPage
		if dec.Decode(&doc) != nil {
			break
		}
		y.Pages = append(y.Pages, doc)
	}
	return
}

// Structure that contains terminal printing functions
// to handle the yaml files containing page structures.
type Terminal struct {
	Color string
}

// Parses the desired horizontal alignment of a string
func (t *Terminal) h_align(str string, align string) (s string) {
	// width of the terminal
	wi := tm.Width()

	// mapping function to a string given then width of the terminal
	applyAlign := func(str string, w int, f func(string, int) string) string {
		// split the string by the new line. This is to apply the mapping to
		// any multi-line string
		tmp := strings.Split(str, "\n")
		for i, l := range tmp {
			// apply the function to the current line
			tmp[i] = f(l, w)
		}
		// join the strings with teh separator and return it
		return strings.Join(tmp, "\n")
	}

	// center the string
	center := func(str string, w int) string {
		return fmt.Sprintf("%*s", (w+len(str))/2, str)
	}

	// right align the string
	right := func(str string, w int) string {
		return fmt.Sprintf("%*s", w, str)
	}

	switch {
	case align == "center":
		s = applyAlign(str, wi, center)

	case align == "right":
		s = applyAlign(str, wi, right)
	default:
		s = str
	}

	return s
}

// Parses the desired vertical alignment of a string
// TODO: Not Implemented yet vertical alignment
//func (t *Terminal) v_align(str string, align string) (s string) {}

// Align a string vertically and horizontally**
// **horizontal alignment not implemented yet!
func (t *Terminal) align(line string, horizontal string, vertical string) (l string) {
	// 1. Align the line horizontally
	l = t.h_align(line, horizontal)

	// 2. Align the line vertically
	// l = t.v_align(line, vertical)
	return l
}

// Animate a string
func (t *Terminal) animate(str string, name string) {
	switch {
	// clear the screen
	case name == "clear":
		tm.Clear()
		tm.Flush()
	// sleep for 3s
	case name == "delay":
		time.Sleep(time.Second * 3)
	// simply prints the line otherwise
	default:
		fmt.Printf("%s\n", str)
	}
}

// Read from the console once
// TODO: finish the logic of this feature
func (t *Terminal) Input() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\u001b[38;5;196muser@riotpot:\u001b[0m~# ")
	text, err := reader.ReadString('\n')

	if err != nil {
		errors.Raise(err)
	}

	return text, err
}

// Parse read a YLine struct.
func (t *Terminal) Read(line YLine) {

	// Get the string result of the alignment
	l := t.align(line.Str, line.H_align, line.V_align)

	// Process the animation
	t.animate(l, line.Animation)

	// Note: Not fully implemented yet.
	//if line.Input {t.Input()}
}

func (t *Terminal) Lecture(path string) {

	// Load the data from the file path
	reader, err := os.ReadFile(path)
	errors.Raise(err)

	// Load a new book
	book := YBook{}
	// Load the pages into the book
	book.Parse(reader)

	// Read the pages line by line
	for _, page := range book.Pages {
		for _, line := range page.Lines {
			t.Read(line)
		}
	}
}
