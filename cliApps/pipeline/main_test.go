package main

import (
	"bytes"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name   string
		dir    string
		expMsg string
		expErr error
	}{
		{name: "No error",
			dir:    "./testdata/tool",
			expMsg: "Go Build: Success\n Go test: Success\n",
			expErr: nil},

		{name: "Error Build",
			dir:    "./testdata/toolerr",
			expMsg: "",
			expErr: &stepErr{step: "Build"}},

		{name: "Error Test",
			dir:    "./testdata/toolerr_test",
			expMsg: "",
			expErr: &stepErr{step: "Test"}},
		{name: "Error fmt",
			dir:    "./testdata/toolerr_fmt",
			expMsg: "",
			expErr: &stepErr{step: "Format"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			err := run(tc.dir, b)
			if tc.expErr == nil {
				if err != nil {
					t.Fatal("expected no error but get error")
				}
			} else {
				if !tc.expErr.(*stepErr).Is(err) {
					t.Fatalf("expected error in %s but got in %s", tc.expErr.(*stepErr).step, err.(*stepErr).step)
				}

			}
		})
	}
}
