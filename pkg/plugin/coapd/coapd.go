// This package implements a CoAP server that listen for connections
// and logs them into a database
// CoAP specs:
// - https://coap.technology/spec.html
// https://tools.ietf.org/html/rfc8974
package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	coap "github.com/plgd-dev/go-coap/v2"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/mux"
	"github.com/riotpot/internal/globals"
	"github.com/riotpot/internal/logger"
	"github.com/riotpot/internal/services"
)

var Plugin string

const (
	name    = "CoAP"
	port    = 5683
	network = globals.UDP
)

func init() {
	Plugin = "Coapd"
}

func Coapd() services.Service {
	mx := services.NewPluginService(name, port, network)

	profile := Profile{
		Topics: RandomNumericTopics("/ps", 10),
	}

	return &Coap{
		mx,
		profile,
	}
}

type Coap struct {
	services.Service
	Profile Profile
}

func (c *Coap) Run() (err error) {

	r := mux.NewRouter()

	// Adds a logger function to the router that we will use
	// to log all the details of the request.
	r.Use(c.loggingMiddleware)

	// Set a default handler for any given path.
	// This will cause all the requests to go through this function.
	r.DefaultHandleFunc(c.observeHandler)

	// Run the server listening on the given port and using the defined
	// lvl4 layer protocol.
	err = coap.ListenAndServe(c.GetNetwork().String(), c.GetAddress(), r)

	return
}

// Method used by the coap mux to log requests.
func (c *Coap) loggingMiddleware(next mux.Handler) mux.Handler {

	// the function takes two arguments: `w` a writer, to respond to
	// the client; and `r` the message to send to the client.
	fn := func(w mux.ResponseWriter, r *mux.Message) {
		// save the client request information
		c.save(w, r)

		// serves the response to the client.
		next.ServeCOAP(w, r)
	}
	return mux.HandlerFunc(fn)
}

func (c *Coap) observeHandler(w mux.ResponseWriter, req *mux.Message) {
	path, err := req.Options.Path()
	if err != nil {
		logger.Log.Error().Err(err)
	}

	// divide the path query into /<topic/sub>/?<flag>=<query>
	r := regexp.MustCompile(`^(?P<topic>(?:\/)?.*?)(?:\?(?P<flag>rt|ct)=(?P<query>(?:.*)))?$`)
	matches := r.FindStringSubmatch(path)

	switch req.Code {
	// if the observe option is set to 0, it means that the client wants to
	// get notified of the current state of a given topic (READ)
	// 		0 (register) adds the entry to the list, if not present.
	// 		1 (deregister) removes the entry from the list, if present.
	case codes.GET:

		// check if the path contains a query, parse it, and understand
		// the request then as a discovery req.
		qIndex := r.SubexpIndex("query")
		q := ""
		if qIndex > -1 {
			q = matches[qIndex]
		}

		if q != "" || path == ".well-known/core" {
			t := matches[r.SubexpIndex("topic")]
			f := matches[r.SubexpIndex("flag")]
			err = c.discovery(w.Client(), req.Token, t, f, q)
			break
		}

		// check if the options came with the observe flag, then it must be a
		// subscription.
		obs, err := req.Options.Observe()
		if err != nil {
			logger.Log.Error().Err(err)
		}

		topic := c.Profile.GetOrCreateTopic(path, "number")

		if obs == 0 {
			// Register the topic and send data
			go c.periodicTransmitter(w.Client(), req.Token, topic)
			break
		}

		if obs == 1 {
			// Deregister the observer and stop sending messages to it for the path.
			break
		}

		// otherwise we consider the connection as a simple request on the state
		msg := c.msg(topic)
		err = c.get(w.Client(), req.Token, msg, -1) // Error is handled below

	// Create a topic. It must indicate the path and the Content Format (ct)
	// It should return the path and a response 2.01 created, since any
	// push to the honeypot will also `create` the topic.
	// TODO: the creation might include an expiration date, we do not include it
	// in the response, it might be nice to implement it.
	case codes.POST, codes.PUT:
		msg := "word"
		if _, err := strconv.Atoi(req.Message.String()); err == nil {
			msg = "number"
		}

		topic := c.Profile.GetOrCreateTopic(path, msg)
		topic.SetMessage(req.Message.String())

		err = c.post(w.Client(), req.Token, path)
	case codes.DELETE:
		err = c.delete(w.Client(), req.Token)
	}

	if err != nil {
		logger.Log.Error().Err(err).Msg("Error on transmitter")
	}
}

// CoAP provides a specific path that returns the list of resources
// normally available at `.well-known/core`, however, it can be at other paths
// depending on the API implementation.
//	 - https://tools.ietf.org/html/draft-ietf-core-coap-pubsub-09#section-4.1
//
// 		Example req: GET /ps/?rt="temperature"
//
// The discovery is meant to localize available topics, however, this is not
// a requirement of the protocol, the server MAY register the topics there. This method
// returns the list of topics under the the query from the `Profile`.
func (c *Coap) discovery(cc mux.Client, token []byte, path string, flag string, query string) error {
	var body []string
	topics := c.Profile.Topics

	// if the path is the `.well-known/core`, use all the topics, otherwise filter them.
	if path != ".well-known/core" {
		topics = filter(topics, path, flag, query)
	}

	// create simple strings following the CoRE Link format
	// - https://tools.ietf.org/id/draft-ietf-core-link-format-11.html#rfc.section.1.2.1
	for _, t := range topics {
		l := fmt.Sprintf(
			`<%s>;%s=%s`,
			t.Path,
			flag,
			query,
		)
		body = append(body, l)
	}

	m := message.Message{
		Code:    codes.Content,
		Context: cc.Context(),
		Body: bytes.NewReader(
			[]byte(strings.Join(body, ",")),
		),
	}

	// declare a packet options object and a buffer in where we can store
	// both the headers and the message type.
	var opts message.Options
	var buf []byte

	// put a header indicating that we will use plain text.
	opts, n, err := opts.SetContentFormat(buf, message.AppLinkFormat)

	// check if the buffer is too small. This might be caused because the bufer was never allocated.
	if err == message.ErrTooSmall {
		buf = append(buf, make([]byte, n)...)
		opts, _, err = opts.SetContentFormat(buf, message.AppLinkFormat)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error on transmitter")
		}
	}

	m.Options = opts

	return cc.WriteMessage(&m)
}

// Implements a response to a `DELETE` request. It responds by including the `DELETED` code.
func (c *Coap) delete(cc mux.Client, token []byte) error {
	m := message.Message{
		Code:    codes.Deleted,
		Token:   token,
		Context: cc.Context(),
	}

	return cc.WriteMessage(&m)
}

// Implements a response to a `POST` request. It responds by including a `CREATED` message
// to the topic for which the message is being pushed, or rather the topic that wants to
// be created.
func (c *Coap) post(cc mux.Client, token []byte, path string) error {
	// create a message from the client connection that we can use to respond.
	m := message.Message{
		Code:    codes.Created,
		Context: cc.Context(),
	}

	// declare a packet options object and a buffer in where we can store
	// both the headers and the message type.
	var opts message.Options
	var buf []byte

	// put a header indicating that we will use plain text.
	opts, n, err := opts.SetContentFormat(buf, message.AppLinkFormat)

	// check if the buffer is too small. This might be caused because the bufer was never allocated.
	if err == message.ErrTooSmall {
		buf = append(buf, make([]byte, n)...)
		opts, _, err = opts.SetContentFormat(buf, message.AppLinkFormat)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error on transmitter")
		}
	}

	// The server MUST add the path to the created topic
	// force behaviour of CREATE on PUBLISH
	opts = opts.Set(message.Option{
		ID:    message.LocationPath,
		Value: []byte(path),
	})

	// send the message with the options
	m.Options = opts
	return cc.WriteMessage(&m)
}

func (c *Coap) get(cc mux.Client, token []byte, msg []byte, obs int64) error {
	// create a message from the client connection that we can use to respond.
	m := message.Message{
		Code:    codes.Content,
		Token:   token,
		Context: cc.Context(),
		Body:    bytes.NewReader(msg),
	}

	// declare a packet options object and a buffer in where we can store
	// both the headers and the message type.
	var opts message.Options
	var buf []byte

	// put a header indicating that we will use plain text.
	opts, n, err := opts.SetContentFormat(buf, message.TextPlain)

	// check if the buffer is too small. This might be caused because the bufer was never allocated.
	if err == message.ErrTooSmall {
		buf = append(buf, make([]byte, n)...)
		opts, _, err = opts.SetContentFormat(buf, message.TextPlain)
	}

	if err != nil {
		return fmt.Errorf("cannot set content format to response: %w", err)
	}

	if obs >= 0 {
		opts, n, err = opts.SetObserve(buf, uint32(obs))
		if err == message.ErrTooSmall {
			buf = append(buf, make([]byte, n)...)

			// set the observer in the options, making the message a notification,
			// this value is simply now a counter for reordering.
			opts, _, err = opts.SetObserve(buf, uint32(obs))
		}
		if err != nil {
			return fmt.Errorf("cannot set options to response: %w", err)
		}
	}

	// send the message with the options
	m.Options = opts
	return cc.WriteMessage(&m)
}

func (c *Coap) msg(topic Topic) []byte {
	msg := topic.Message()
	return []byte(msg)
}

func (c *Coap) periodicTransmitter(cc mux.Client, token []byte, topic Topic) {

	obs := int64(2)

	// create a channel from where we can read that will sleep and then put
	// a value into the channel.
	// when the channel is filled, the loop breaks
	stop := make(chan int, 1)
	go timeout(60*time.Second, stop) // sleeps for a minute before timeout

	for {
		select {
		case <-stop:
			return
		default:
			msg := c.msg(topic)
			err := c.get(cc, token, msg, obs)
			if err != nil {
				logger.Log.Error().Err(err).Msg("Error on transmitter. Shutting down.")
				return
			}

			// sleep for few seconds before trying to send a message again
			secs := time.Duration(rand.Intn(60) + 1)
			time.Sleep(secs * time.Second)
		}
	}
}

func (c *Coap) save(w mux.ResponseWriter, r *mux.Message) {
	// rmtAddr := w.Client().RemoteAddr()
	// path :=
}

// Filter a list of topics based on the query string included and the flag
func filter(topics []Topic, path string, flag string, query string) (t []Topic) {
	// Replace all the quotation marks on the query.
	// Queryes might contain quotation marks e.g. `/ps/a?ct="50"`
	query = strings.Replace(query, `"`, "", -1)

	switch flag {
	case "ct":
		if query == "0" {
			t = append(t, topics...)
		}
	case "rt":
		for _, topic := range topics {
			if strings.HasPrefix(topic.Path, path) && strings.Contains(topic.Path, query) {
				t = append(t, topic)
			}
		}
	// TODO: implement this filters. They might have an impact on the `Topic` struct though!
	case "anchor", "rel", "if":
	}

	return
}

func timeout(t time.Duration, stop chan int) {
	time.Sleep(t)
	stop <- 1
}
