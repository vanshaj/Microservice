package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

func mockCmdContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess"}
	cs = append(cs, name)
	cs = append(cs, arg...)
	cmd := exec.CommandContext(ctx, os.Args[0], cs...)
	cmd.Env = []string{"HELPER_PROCESS=1"}
	return cmd
}

func mockCmdTimeout(ctx context.Context, exe string, args ...string) *exec.Cmd {
	cmd := mockCmdContext(ctx, exe, args...)
	cmd.Env = append(cmd.Env, "GO_HELPER_TIMEOUT=1")
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("HELPER_PROCESS") != "1" {
		return
	}
	if os.Getenv("GO_HELPER_TIMEOUT") == "1" {
		time.Sleep(15 * time.Second)
	}
	if os.Args[2] == "git" {
		fmt.Fprintln(os.Stdout, "Everything up-to-date")
		os.Exit(0)
	}
	os.Exit(1)
}

func TestRun(t *testing.T) {
	testCases := []struct {
		name    string
		dir     string
		expMsg  string
		expErr  error
		mockCmd func(ctx context.Context, name string, arg ...string) *exec.Cmd
	}{
		{name: "No error",
			dir:     "./testdata/tool",
			expMsg:  "Go Build: Success\n Go test: Success\nGofmt: Success\n Git Push: Success\n",
			expErr:  nil,
			mockCmd: mockCmdContext},

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
		{name: "Error timeout",
			dir:     "./testdata/tool",
			expMsg:  "",
			expErr:  &stepErr{step: "git push"},
			mockCmd: mockCmdTimeout},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockCmd != nil {
				command = tc.mockCmd
			}

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
