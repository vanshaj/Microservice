package main

import (
	"flag"
	"fmt"
	"os"

	todo "github.com/vanshaj/Microservice/cliApps/todoapp"
)

func main() {
	fileName := "todo.json"
	list := flag.Bool("list", false, "List all tasks")
	task := flag.String("task", "", "Task to be included in the todo list")
	complete := flag.String("complete", "", "Item to be completed")
	flag.Parse()

	l := &todo.List{}

	switch {
	case *list:
		l.Get(fileName)
		for _, item := range *l {
			if item.Done != true {
				fmt.Println(item.Task)
			}
		}
	case *complete != "":
		l.Get(fileName)
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := l.Save(fileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		} else {
			fmt.Printf("%s task has been completed", *complete)
		}
	case *task != "":
		l.Get(fileName)
		l.Add(*task)
		if err := l.Save(fileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

}
