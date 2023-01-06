package main

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

var positiveTestData = `name,data
ak,20
bc,30`

func TestCsvToFloat(t *testing.T) {
	testCases := []struct {
		name          string
		data          io.Reader
		column        int
		expected      []float64
		expectedError error
	}{
		{
			name:          "correct column",
			data:          strings.NewReader(positiveTestData),
			column:        1,
			expected:      []float64{20, 30},
			expectedError: nil,
		},
		{
			name:          "incorrect column",
			data:          strings.NewReader(positiveTestData),
			column:        3,
			expected:      nil,
			expectedError: ErrInvalidColumn,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := csvToFloat(tc.data, tc.column)
			if err != nil {
				if err != tc.expectedError {
					t.Fatalf("got %s wanted %s", err, tc.expectedError)
				}
			} else {
				if !reflect.DeepEqual(output, tc.expected) {
					t.Fatalf("out is not expected")
				}
			}
		})
	}
}
