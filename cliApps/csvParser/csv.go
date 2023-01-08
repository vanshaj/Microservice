package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

type statFunc func([]float64) float64

func (s statFunc) operation(data []float64) float64 {
	return s(data)
}

type stat interface {
	operation(data []float64) float64
}

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
	csvR.ReuseRecord = true
	returnData := make([]float64, 0)
	for i := 0; ; i++ {
		row, err := csvR.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("can't read data from file")
		}
		if i == 0 {
			continue
		}
		if len(row) <= column {
			return nil, ErrInvalidColumn
		}
		value, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, err
		}
		returnData = append(returnData, value)
	}
	return returnData, nil
}
