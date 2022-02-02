package tasker

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/gookit/color.v1"
	"time"
)

func Commands() []*cli.Command {
	var commands []*cli.Command

	commands = []*cli.Command{
		addTask(),
		listTasks(),
		markCompleteTasks(),
		finishedTasks(),
		deleteTask(),
	}
	return commands
}

func DefaultAction() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		tasks, err := PendingTasks()
		if err != nil {
			if err == mongo.ErrNoDocuments {
				fmt.Println("Nothing to look here")
				return nil
			}
			return err
		}
		printTasks(tasks)
		return nil
	}
}

func addTask() *cli.Command {
	add := cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "add a task to the list",
		Action: func(context *cli.Context) error {
			str := context.Args().First()
			if str == "" {
				return errors.New("cannot add empty task")
			}
			task := &Task{
				ID:        primitive.NewObjectID(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Text:      str,
				Completed: false,
			}
			return CreateTask(task)
		},
	}
	return &add
}

func listTasks() *cli.Command {
	return &cli.Command{
		Name:    "all",
		Aliases: []string{"l"},
		Usage:   "list all tasks",
		Action: func(context *cli.Context) error {
			tasks, err := GetAllTasks()
			if err != nil {
				if err == mongo.ErrNoDocuments {
					fmt.Println("No tasks to list\nStart creating tasks with command : add 'task'")
					return nil
				}
				return err
			}
			printTasks(tasks)
			return nil
		},
	}
}

func markCompleteTasks() *cli.Command {
	return &cli.Command{
		Name:    "done",
		Aliases: []string{"d"},
		Usage:   "complete a task on the list",
		Action: func(context *cli.Context) error {
			text := context.Args().First()
			return MarkCompleted(text)
		},
	}
}

func finishedTasks() *cli.Command {
	return &cli.Command{
		Name:    "finished",
		Aliases: []string{"f"},
		Usage:   "show all finished tasks on the list",
		Action: func(context *cli.Context) error {
			tasks, err := FinishedTasks()
			if err != nil {
				if err == mongo.ErrNoDocuments {
					fmt.Println("No tasks to list\nStart creating tasks with command : add 'task'")
					return nil
				}
				return err
			}
			printTasks(tasks)
			return nil
		},
	}
}

func deleteTask() *cli.Command {
	return &cli.Command{
		Name:    "rm",
		Aliases: []string{"rm"},
		Usage:   "deletes a task on the list",
		Action: func(context *cli.Context) error {
			text := context.Args().First()
			return DeleteTask(text)
		},
	}
}

func printTasks(tasks []*Task) {
	for i, task := range tasks {
		if task.Completed {
			color.Green.Printf("%d: %s\n", i+1, task.Text)
		} else {
			color.Red.Printf("%d: %s\n", i+1, task.Text)
		}
	}
}
