package main

import (
	"errors"
)

type command struct {
	name string
	args []string
}

type commands struct {
	cmdMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	val, ok := c.cmdMap[cmd.name]
	if !ok {
		return errors.New("command doesn't exist")
	}

	err := val(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	_, ok := c.cmdMap[name]
	if !ok {
		c.cmdMap[name] = f
	}
}
