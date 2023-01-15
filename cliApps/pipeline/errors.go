package main

import "fmt"

type stepErr struct {
	step  string
	msg   string
	cause error
}

func (s *stepErr) Error() string {
	return fmt.Sprintf("step %q: %s: Cause: %v", s.step, s.msg, s.cause)
}

func (s *stepErr) Is(t error) bool {
	target, ok := t.(*stepErr)
	if !ok {
		return false
	}
	return s.step == target.step
}
