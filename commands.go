package main

import "fmt"

type commands struct {
	cmdNamesMap map[string]func(*state, command) error
}

func (cmds *commands) register(name string, f func(*state, command) error) {
	cmds.cmdNamesMap[name] = f
}

func (cmds *commands) run(s *state, cmd command) error {
	val, ok := cmds.cmdNamesMap[cmd.name]
	if !ok {
		return fmt.Errorf("command %s does not exist", cmd.name)
	}

	err := val(s, cmd)
	if err != nil {
		return fmt.Errorf("command failed to execute: %v", err)
	}

	return nil
}
