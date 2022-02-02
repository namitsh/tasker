package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"tasker"
)

func main() {
	app := &cli.App{
		Name:     "Tasker",
		Usage:    "A simple CLI program to manage your tasks",
		Action:   tasker.DefaultAction(),
		Commands: tasker.Commands(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
