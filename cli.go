package main

import (
	"errors"
	"fmt"
)

var _ = errors.New("")

type SubCommand struct {
	Name        string
	Description string
	Options     []Option
	Run         func(_ SubCommand, _ []string)
}

type Option struct {
	Description string
	Flag        string
}

type Application struct {
	Name        string
	Author      string
	Description string
	Usage       string
	SubCommands []SubCommand
}

func (app *Application) Run(args []string) error {
	if len(args) == 1 {
		app.Help()
		return nil
	}

	var command SubCommand
	for _, c := range app.SubCommands {
		if c.Name == args[1] {
			command = c
			break
		}
	}

	if command.Run == nil {
		app.Help()
		return nil
	}

	command.Run(command, args)
	return nil
}

func (app *Application) Help() error {
	fmt.Print(app.Name, "\n")
	fmt.Print(app.Description, "\n")
	fmt.Print("\nUSAGE:\n\t", app.Usage, "\n")
	/*
		fmt.Print("\nOPTIONS:\n")
		for _, option := range app.Options {
			fmt.Print("\t")
			for i := 0; i < 10; i++ {
				if i < len(option.Flag) {
					fmt.Print(string(option.Flag[i]))
				} else {
					fmt.Print(" ")
				}
			}
			fmt.Print("\t", option.Description, "\n")
		}
	*/
	fmt.Print("\nSUBCOMMANDS:\n")
	for _, SubCommand := range app.SubCommands {
		fmt.Print("\t")
		for i := 0; i < 10; i++ {
			if i < len(SubCommand.Name) {
				fmt.Print(string(SubCommand.Name[i]))
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print("\t", SubCommand.Description, "\n")
	}
	fmt.Println("\nBy", app.Author)
	return nil
}
