package main

import (
	"context"
	"os/exec"
	"time"
)

type timeoutStep struct {
	step
	timeout time.Duration
}

var command = exec.CommandContext

func newTimeoutStep(name, exe, message, proj string, args []string, timeout time.Duration) *timeoutStep {
	s := &timeoutStep{}
	s.step = *newStep(name, exe, message, proj, args)
	if timeout == 0 {
		s.timeout = 30 * time.Second
	} else {
		s.timeout = timeout
	}
	return s
}

func (s *timeoutStep) execute() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	cmdCtx := command(ctx, s.exe, s.args...)
	cmdCtx.Dir = s.proj
	if err := cmdCtx.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", &stepErr{
				step:  s.name,
				msg:   "failed timeout",
				cause: context.DeadlineExceeded,
			}
		}
		return "", &stepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err,
		}
	}
	return s.message, nil
}
