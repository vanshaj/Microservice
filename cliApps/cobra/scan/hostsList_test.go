package scan

import (
	"bytes"
	"errors"
	"testing"
)

func TestAdd(t *testing.T) {
	testCases := []struct {
		name        string
		host        string
		expectedLen int
		expectedErr error
	}{
		{"Add single item", "192.168.29.5", 1, nil},
		{"Add no item", "", 0, ErrEmptyHost},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := &HostsList{}
			actualErr := h.Add(tc.host)
			if !errors.Is(actualErr, tc.expectedErr) {
				t.Fatalf("Expected error not mataches actual error. Expected '%v' but Actual '%v'", actualErr, tc.expectedErr)
			}
			if len(h.Hosts) != tc.expectedLen {
				t.Fatalf("Expected len doesnot matches actual length")
			}
		})
	}
}

func TestLoad(t *testing.T) {
	testCases := []struct {
		name        string
		ips         string
		expectedLen int
		expectedErr error
	}{
		{"Load 3 ips", "192.168.29.1\n192.168.29.2\n192.168.29.3", 3, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var b bytes.Buffer
			b.Write([]byte(tc.ips))
			h := &HostsList{}
			actualErr := h.Load(&b)
			if !errors.Is(actualErr, tc.expectedErr) {
				t.Fatalf("Expected error not mataches actual error. Expected '%v' but Actual '%v'", actualErr, tc.expectedErr)
			}
			if len(h.Hosts) != tc.expectedLen {
				t.Fatalf("Expected len doesnot matches actual length")
			}
		})
	}
}
