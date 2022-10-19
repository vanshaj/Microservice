package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	todo "github.com/vanshaj/Microservice/cliApps/todoapp"
)

func main() {
	fileName := "todo.json"
	list := flag.Bool("list", false, "List all tasks")
	task := flag.String("task", "", "Task to be included in the todo list")
	complete := flag.String("complete", "", "Item to be completed")
	add := flag.Bool("add", false, "Add Task to the Todo List")
	flag.Parse()

	l := &todo.List{}

	switch {
	case *add:
		l.Get(fileName)
		taskData, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		l.Add(taskData)
		if err := l.Save(fileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
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

func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	if len(scanner.Text()) == 0 {
		return "", fmt.Errorf("task cannot be blank")
	}
	return scanner.Text(), nil
}
