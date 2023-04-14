package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	errNotFound       = errors.New("Not found")
	errInternalServer = errors.New("Internal Server Error")
	errBadRequest     = errors.New("Bad Request")
)

type Item struct {
	Task        string    `json:"Task"`
	Done        bool      `json:"Done"`
	CreateAt    time.Time `json:"CreatedAt"`
	CompletedAt time.Time `json:"CompletedAt"`
}

type Response struct {
	Results []Item `json:"results"`
}

func (resp Response) ToJson(w io.Writer) error {
	return json.NewEncoder(w).Encode(resp)
}

func newClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}

func GetAll(apiRoot string) (Response, error) {
	c := newClient()
	resp, err := c.Get(fmt.Sprintf("%s/todo", apiRoot))
	if err != nil {
		return Response{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return Response{}, errors.New(fmt.Sprintf("Status code returned is %d", resp.StatusCode))
	}
	content, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	items := Response{}
	err = json.NewDecoder(strings.NewReader(string(content))).Decode(&items)
	if err != nil {
		return Response{}, errBadRequest
	}
	return items, nil
}

func AddOne(apiRoot string, taskName string) (Response, error) {
	c := newClient()
	resp, err := c.Post(fmt.Sprintf("%s/todo/%s", apiRoot, taskName), "application/json", strings.NewReader(""))
	if err != nil {
		return Response{}, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return Response{}, errors.New(fmt.Sprintf("Status code returned is %d", resp.StatusCode))
	}
	content, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	items := Response{}
	err = json.NewDecoder(strings.NewReader(string(content))).Decode(&items)
	if err != nil {
		return Response{}, errBadRequest
	}
	return items, nil
}
