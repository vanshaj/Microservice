package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"

	todo "github.com/vanshaj/Microservice/cliApps/todoapp"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidData = errors.New("invalid data")
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	content := "API is here\n"
	replyTextContent(w, r, http.StatusOK, content)
}

func getAllHandler(w http.ResponseWriter, r *http.Request, list *todo.List) {
	resp := &todoResponse{
		Results: *list,
	}
	replyJSONContent(w, r, http.StatusOK, resp)
}

func getOneHandler(w http.ResponseWriter, r *http.Request, list *todo.List, taskName string) {
	index := -1
	for i, val := range *list {
		if val.Task == taskName {
			index = i
		}
	}
	if index == -1 {
		replyError(w, r, http.StatusNotFound, "Invalid task")
		return
	}
	resp := &todoResponse{}
	if index-1 == len(*list) {
		resp.Results = (*list)[index:]
	} else {
		resp.Results = (*list)[index : index+1]
	}
	replyJSONContent(w, r, http.StatusOK, resp)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request, list *todo.List, taskName string, todoFile string) {
	list.Add(taskName)
	if err := list.Save(todoFile); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	replyEmptyContent(w, r, http.StatusCreated)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request, list *todo.List, taskName string, todoFile string) {
	err := list.Delete(taskName)
	if err != nil {
		replyError(w, r, http.StatusBadRequest, err.Error())
	}
	list.Save(todoFile)
	replyEmptyContent(w, r, http.StatusOK)
}

func validateID(path string, list *todo.List) (int, error) {
	id, err := strconv.Atoi(path)
	if err != nil {
		return 0, fmt.Errorf("%w: Invalid ID: %s", ErrInvalidData, err)
	}

	if id < 1 {
		return 0, fmt.Errorf("%w, Invalid ID: Less than one", ErrInvalidData)
	}

	if id > len(*list) {
		return id, fmt.Errorf("%w: ID %d not found", ErrNotFound, id)
	}

	return id, nil
}

func todoRouter(todoFile string, l sync.Locker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list := &todo.List{}
		l.Lock()
		defer l.Unlock()
		file_size, err := os.Stat(todoFile)
		if err != nil {
			replyError(w, r, http.StatusInternalServerError, err.Error())
		}
		if file_size.Size() != 0 {
			list.Get(todoFile)
		}
		if r.URL.Path == "" {
			switch r.Method {
			case http.MethodGet:
				getAllHandler(w, r, list)
			default:
				message := "Method not supported"
				replyError(w, r, http.StatusMethodNotAllowed, message)
			}
			return
		}
		taskName, err := url.QueryUnescape(r.URL.Path)
		if err != nil {
			replyError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		switch r.Method {
		case http.MethodGet:
			getOneHandler(w, r, list, taskName)
		case http.MethodPost:
			addTaskHandler(w, r, list, taskName, todoFile)
		case http.MethodDelete:
			deleteTaskHandler(w, r, list, taskName, todoFile)
		default:
			message := "Method not supported"
			replyError(w, r, http.StatusMethodNotAllowed, message)
		}
	}
}
