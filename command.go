// Copyright 2013 Google Inc. All Rights Reserved.
// Modifications Copyright 2016 Drachenfels GmbH. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package command allows you to define subcommands
// for your command line interfaces. It extends the flag package
// to provide flag support for subcommands.
package command

import (
	"errors"
	"flag"
	"fmt"
)

// A map of all of the registered sub-commands.
type Path struct {
	entries map[string]*CmdCont
}

func NewPath() *Path {
	return &Path{
		entries: make(map[string]*CmdCont),
	}
}

var (
	ErrCmdUsage  = errors.New("Invalid command usage.")
	ErrNoSuchCmd = errors.New("No such command.")
)

// Cmd represents a sub command, allowing to define subcommand
// flags and runnable to run once arguments match the subcommand
// requirements.
type Cmd interface {
	// Callback used to register flags for the subcommand
	Flags(*flag.FlagSet)
	Run(args ...string) error
}

// A func that implements the Cmd interface.
// For registering simple commands without flags.
type CmdFunc func(args []string) error

func (s CmdFunc) Flags(fs *flag.FlagSet) {
}

func (s CmdFunc) Run(args ...string) error {
	return s(args)
}

type CmdCont struct {
	Cmd
	Name          string
	Desc          string
	RequiredFlags []string
	Flags         *flag.FlagSet
}

// Registers a Cmd for the provided sub-command Name.
// E.g. Name is the `status` in `git status`.
func (p *Path) Add(name, description string, command Cmd, requiredFlags ...string) *CmdCont {
	c := &CmdCont{
		Cmd:           command,
		Name:          name,
		Desc:          description,
		RequiredFlags: requiredFlags,
		Flags:         flag.NewFlagSet(name, flag.ContinueOnError),
	}
	// register subcommand flags
	c.Cmd.Flags(c.Flags)
	// TODO warn before overwriting an existing command ?
	p.entries[name] = c
	return c
}

// Parses the flags and leftover arguments to match them with a
// sub-command. Evaluate all of the global flags and register
// sub-command handlers before calling it. Sub-command handler's
// `Run` will be called if there is a match.
// A usage with flag defaults will be printed if provided arguments
// don't match the configuration.
// Global flags are accessible once Parse executes.
func (p *Path) Run(args ...string) (*CmdCont, error) {
	// if there are no subcommands registered,
	// return immediately
	if len(p.entries) < 1 || len(args) < 1 {
		return nil, ErrCmdUsage
	}
	// first argument is the subcommand
	if cont, ok := p.entries[args[0]]; ok {
		if len(args) > 1 {
			err := cont.Flags.Parse(args[1:])
			if err != nil {
				return cont, err
			}
		}

		// check for required / mandatory flags.
		missingFlags := make(map[string]bool)
		for _, flagName := range cont.RequiredFlags {
			missingFlags[flagName] = true
		}
		cont.Flags.Visit(func(f *flag.Flag) {
			delete(missingFlags, f.Name)
		})

		if len(missingFlags) > 0 {
			keys := make([]string, 0, len(missingFlags))
			for k := range missingFlags {
				keys = append(keys, k)
			}
			return cont, fmt.Errorf("Required flags not set: %q\n", keys)
		}
		return cont, cont.Run(cont.Flags.Args()...)
	}
	return nil, ErrNoSuchCmd
}

func (p *Path) PrintAvailableCommands() {
	fmt.Println("Available commands:")
	for _, c := range p.entries {
		fmt.Printf("\t%s\t%s\n", c.Name, c.Desc)
	}
}

var globalPath = NewPath()

func Add(name, description string, command Cmd, requiredFlags ...string) *CmdCont {
	return globalPath.Add(name, description, command, requiredFlags...)
}

func PrintAvailableCommands() {
	globalPath.PrintAvailableCommands()
}

func Run(args ...string) (*CmdCont, error) {
	return globalPath.Run(args...)
}
