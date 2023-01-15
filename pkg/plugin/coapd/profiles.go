// Package that implements the creation of profiles that resemble real devices,
// including topics that clients can subscribe to, and generate random values.
//
// TODO include a method to load profile topics and determine whether the returning
// value should be a number in a range, a string or something else.
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/riotpot/internal/logger"
	"gopkg.in/yaml.v3"
)

const (
	WORD   = "word"
	NUMBER = "number"
)

type Profile struct {
	Name    string // Name of the device that will be faked
	Version string // Version of the device
	Banner  string // Banner of the device

	Topics []Topic // list of topics registered in the profile
}

func (p *Profile) Load(path string) {
	data, err := os.ReadFile(path)
	logger.Log.Error().Err(err)
	err = yaml.Unmarshal(data, &p)
	logger.Log.Error().Err(err)
}

// Method that provides a getter for topics and anso creates the topic
// and addds it to the profile if it didn't exist previously
func (p *Profile) GetOrCreateTopic(path string, msgType string) Topic {
	t, i := p.Topic(path)

	// if the index of the topic is -1, the topic could not be found
	// in the list of topics in the profile.
	if i == -1 {
		// update the empty topic returned from the `Topic` function.
		// NOTE: there is no implementation to differentiate the topic
		// content, therefore, we just include a random starting point
		// and a default range from 0 to 100.
		t.Path = path

		switch msgType {
		case NUMBER:
			t.MessageType = NUMBER
			t.Interval = [2]float32{0, 100}
			t.pNum = (rand.Float32() * (100 - 0)) + 0
		case WORD:
			t.MessageType = WORD
		}
	}

	return t
}

// Method to get a topic based on the path.
// Note: The path is currently very strict, it would be nice to ease it.
func (p *Profile) Topic(path string) (topic Topic, index int) {
	for i, t := range p.Topics {
		if t.Path == path {
			topic = t
			index = i
			return
		}
	}
	// update the index to -1
	index = -1

	return
}

type Topic struct {
	Path string // path string e.g. `/path/to/topic/*`

	MessageType string     // type of the message
	Words       []string   // list of words that can be used on the topic
	Interval    [2]float32 // floor and ceil values of the topic

	pNum  float32 // Previous number
	pWord string  // Previous word
}

func (t *Topic) PathAndMessage() (string, string) {
	return t.Path, t.Message()
}

func (t *Topic) SetMessage(msg string) {
	if t.MessageType == WORD {
		t.pWord = msg
	} else if t.MessageType == NUMBER {
		value, err := strconv.ParseFloat(msg, 32)
		if err != nil {
			panic("Value could not be converted to float32")
		}
		t.pNum = float32(value)
	}
}

// Implements a generic method that generates a string message
// based on the characteristics of the topic.
func (t *Topic) Message() string {
	var msg string

	// we add 1 to the result just in case is 0.
	index := rand.Intn(100) + 1

	switch t.MessageType {
	case WORD:
		msg = t.word(index)
		t.pWord = msg
	case NUMBER:
		n := t.genFloatNum(0.1)
		t.pNum = n
		msg = fmt.Sprintf("%v", n)
	}

	return msg
}

// Implements a getter for the attribute `words` based on an integer
func (t *Topic) word(i int) string {
	i = i % len(t.Words)
	return t.Words[i]
}

// Generates a number based on a previous seed and the topic range.
// It takes a weight value that considers how far at maximum we want the
// result to be from the given number in between the range.
// The weight must be a value between 0-1 as a percercentage of the value,
// however, we expect to see very low values such as 0.01!
func (t *Topic) genFloatNum(weight float32) float32 {
	// check if the weight is higher or lower than the extremes, and if so,
	// recalculate it using a pseudo-random float value.
	if weight > 1 || weight < 0 {
		weight = float32(rand.Float64())
	}

	min, max := t.Interval[0], t.Interval[1]

	// calculate a window of values that we want to use around
	// the previous number.
	variance := t.pNum * weight

	// check if the variance minimum is higher than the local minimum
	v_min := t.pNum - variance
	if t.pNum-variance > min {
		min = v_min
	}

	// check if the variance maximum is lower than the local maximum
	v_max := t.pNum + variance
	if v_max < max {
		max = v_max
	}

	// calculate a random float32 in between the given range
	val := (rand.Float32() * (max - min)) + min

	return val
}

// Function that creates a random amount of topics using uuid's for
// the topics and pathing.
func RandomNumericTopics(path string, n int) (topics []Topic) {
	rand.Seed(time.Now().UnixNano())
	// random amount of topics from with n as max
	n = rand.Intn(n)
	// random amount of sub-topics and levels
	nSubs := rand.Intn(n + 1)
	maxLevels := rand.Intn(nSubs + 1)

	paths := []string{}

	for i := 0; i <= n; i++ {
		var p string

		// update the number of subs to at least 1
		if nSubs < 1 {
			nSubs = 1
		}

		// increase the probability of hitting an existing path
		s := rand.Intn((nSubs + 1) / 2)

		if s >= len(paths)-1 || len(paths) <= 1 {
			l := rand.Intn(maxLevels + 1)
			p = pathUuid(path, l)
			paths = append(paths, p)
		} else {
			p = paths[s]
		}

		t := randTopicNumber(p)
		topics = append(topics, t)

	}

	return
}

// Creates a full path of n amount of `levels` with a prefix
func pathUuid(prefix string, levels int) string {
	p := []string{prefix}

	for i := 0; i < levels; i++ {
		u := uuid.NewString()
		p = append(p, u)
	}
	return strings.Join(p, "/")
}

func randTopicNumber(path string) Topic {
	u := uuid.NewString()
	path = fmt.Sprintf("%s/%s", path, u)
	return Topic{
		Path:        path,
		MessageType: "number",
		Interval:    [2]float32{0, float32(rand.Intn(100))},
	}
}
