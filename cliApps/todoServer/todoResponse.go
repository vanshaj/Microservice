package main

import (
	"encoding/json"
	"time"

	todoapp "github.com/vanshaj/Microservice/cliApps/todoapp"
)

type todoResponse struct {
	Results todoapp.List `json:"results"`
}

func (r *todoResponse) MarshalJSON() ([]byte, error) {
	resp := struct {
		Results      todoapp.List `json:"results"`
		Date         int64        `json:"date"`
		TotalResults int          `json:"totalResults"`
	}{
		Results:      r.Results,
		Date:         time.Now().Unix(),
		TotalResults: len(r.Results),
	}
	return json.Marshal(resp)
}
