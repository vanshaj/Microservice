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
	del := flag.Bool("del", false, "Delete the task from the list")
	listAll := flag.Bool("listall", false, "List only incomplete items")
	flag.Parse()

	l := &todo.List{}
	l.Get(fileName)

	switch {
	case *add:
		taskData, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, eachTask := range taskData {
			l.Add(eachTask)
		}
		if err := l.Save(fileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *del:
		taskData, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, eachTask := range taskData {
			l.Delete(eachTask)
		}
		if err := l.Save(fileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *list:
		for _, item := range *l {
			if item.Done != true {
				fmt.Println(item.Task)
			}
		}
	case *complete != "":
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
		l.Add(*task)
		if err := l.Save(fileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *listAll:
		for _, val := range *l {
			fmt.Println(val.Task)
		}
	}

}

func getTask(r io.Reader, args ...string) ([]string, error) {
	if len(args) > 0 {
		return []string{strings.Join(args, " ")}, nil
	}
	scanner := bufio.NewScanner(r)
	output := make([]string, 0, 20)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		} else {
			output = append(output, scanner.Text())
		}
	}
	return output, nil
}
