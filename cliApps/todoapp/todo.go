package todo

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type List []item

func (items *List) findIndex(task string) (int, error) {
	index := -1
	for i, val := range *items {
		if val.Task == task {
			index = i
		}
	}
	if index == -1 {
		return -1, errors.New("no such task is present")
	}
	return index, nil
}

func (items *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	*items = append(*items, t)
}

func (items *List) Delete(task string) error {
	ls := *items
	index, err := items.findIndex(task)
	if err != nil {
		return err
	}
	*items = append(ls[:index], ls[index+1:]...)
	return nil
}

func (items *List) Complete(task string) error {
	index, err := items.findIndex(task)
	if err != nil {
		return err
	}
	ls := *items
	ls[index].Done = true
	ls[index].CompletedAt = time.Now()
	return nil
}

func (items *List) Save(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	fileEncoder := json.NewEncoder(file)
	err = fileEncoder.Encode(items)
	return err
}

func (items *List) Get(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	fileDecoder := json.NewDecoder(file)
	err = fileDecoder.Decode(items)
	return err
}
