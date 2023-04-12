package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

func addHandler(w http.ResponseWriter, r *http.Request, list *todo.List, todoFile string) {
	item := struct {
		Task string `json:"task"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("Invalid json: %ss", err)
		replyError(w, r, http.StatusBadRequest, message)
		return
	}
	list.Add(item.Task)
	if err := list.Save(todoFile); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	replyTextContent(w, r, http.StatusCreated, "")
}

func getOneHandler(w http.ResponseWriter, r *http.Request, list *todo.List, id int) {
	resp := &todoResponse{
		Results: (*list)[id-1 : id],
	}
	replyJSONContent(w, r, http.StatusOK, resp)
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

func replyTextContent(w http.ResponseWriter, r *http.Request, statusCode int, content string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	w.Write([]byte(content))
}

func todoRouter(todoFile string, l sync.Locker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list := &todo.List{}
		l.Lock()
		defer l.Unlock()
		if err := list.Get(todoFile); err != nil {
			replyError(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		if r.URL.Path == "" {
			switch r.Method {
			case http.MethodGet:
				getAllHandler(w, r, list)
			case http.MethodPost:
				addHandler(w, r, list, todoFile)
			default:
				message := "Method not supported"
				replyError(w, r, http.StatusMethodNotAllowed, message)
			}
			return
		}
		id, err := validateID(r.URL.Path, list)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				replyError(w, r, http.StatusNotFound, err.Error())
				return
			}
			replyError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		switch r.Method {
		case http.MethodGet:
			getOneHandler(w, r, list, id)
		default:
			message := "Method not supported"
			replyError(w, r, http.StatusMethodNotAllowed, message)
		}
	}
}
