package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
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

func TestRunKill(t *testing.T) {
	var testCases = []struct {
		name   string
		proj   string
		sig    syscall.Signal
		expErr error
	}{
		{"SIGINT", "./testdata/tool", syscall.SIGINT, ErrSignal},
		{"SIGTERM", "./testdata/tool", syscall.SIGTERM, ErrSignal},
		{"SIGQUIT", "./testdata/tool", syscall.SIGQUIT, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			command = mockCmdTimeout
			errCh := make(chan error)
			ignoreChannel := make(chan os.Signal, 3)
			expectChannel := make(chan os.Signal, 1)

			signal.Notify(ignoreChannel, syscall.SIGQUIT)
			defer signal.Stop(ignoreChannel)

			signal.Notify(expectChannel, tc.sig)
			defer signal.Stop(expectChannel)

			go func() {
				errCh <- run(tc.proj, io.Discard)
			}()

			go func() {
				time.Sleep(2 * time.Second)
				syscall.Kill(syscall.Getpid(), tc.sig)
			}()

			select {
			case err := <-errCh:
				if err == nil {
					t.Errorf("Expected error. Got 'nil' instead.")
					return
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expected error: %q. Got %q", tc.expErr, err)
				}
				select {
				case rec := <-expectChannel:
					if rec != tc.sig {
						t.Errorf("Expected signal: %q. Got %q", tc.sig, rec)
					}
				default:
					t.Errorf("signal not received")
				}
			case <-ignoreChannel:
			}
		})
	}
}
