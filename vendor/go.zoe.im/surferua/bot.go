package surferua

import (
	"math/rand"
	"time"
)

var botDBSize = 0
var botDB = []string{}

// Bot Crawler  UA strings
type Bot struct {
	// The name of bot
	Name string

	// The version of bot
	Version string

	// The url of bot
	URL string
}

// String generate strings of UserAgent
func (b *Bot) String() (s string) {
	return b.Name + "/" + b.Version + " (+" + b.URL + ")"
}

// NewBot returns a bot ua randomly
func NewBot() string {
	return botDB[rand.Intn(botDBSize)]
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
