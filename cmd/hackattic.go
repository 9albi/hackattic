package cmd

import (
	"errors"
	"flag"
	"fmt"

	"github.com/9albi/hackattic/pkg/challenge"
)

var Challenges = map[string]challenge.Challenge{
	"mini_miner": MiniMiner{},
}

func Run() error {
	var challengeName string
	flag.StringVar(&challengeName, "challenge", "undefined", "challenge name")
	flag.Parse()

	if challengeName == "" {
		return errors.New("challenge name cannot be empty")
	}

	challenge, found := Challenges[challengeName]
	if !found {
		return fmt.Errorf("couldn't find challenge '%s'", challengeName)
	}

	return challenge.Solve()
}
