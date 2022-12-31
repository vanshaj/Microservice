package main

import (
	"encoding/csv"
	"io"
	"strconv"
)

type statFunc func(data []float64) float64

func sum(data []float64) float64 {
	var sum float64
	for _, val := range data {
		sum = sum + val
	}
	return sum
}

func avg(data []float64) float64 {
	total := sum(data)
	size := float64(len(data))
	return total / size
}

func csvToFloat(r io.Reader, column int) ([]float64, error) {
	csvR := csv.NewReader(r)
	data, err := csvR.ReadAll()
	if err != nil {
		return nil, err
	}
	returnData := make([]float64, 0)
	if len(data[0]) <= column {
		return nil, ErrInvalidColumn
	}
	for i, row := range data {
		if i == 0 {
			continue
		}
		value, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, err
		}
		returnData = append(returnData, value)
	}
	return returnData, nil
}
