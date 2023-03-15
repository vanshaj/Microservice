package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/vanshaj/Microservice/cliApps/cobra/scan"
)

func setup(t *testing.T, hosts []string, initList bool) (string, func()) {
	f, err := os.CreateTemp("/tmp", "pscan")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if initList {
		h1 := &scan.HostsList{}
		for _, h := range hosts {
			h1.Add(h)
		}
		if err := h1.Save(f); err != nil {
			t.Fatal(err)
		}
	}
	return f.Name(), func() {
		os.Remove(f.Name())
	}
}

func TestHostActions(t *testing.T) {
	hosts := []string{"host1", "host2", "host3"}
	testCases := []struct {
		name        string
		args        []string
		expectedOut string
		initList    bool
		actionFunc  func(io.Writer, string, []string) error
	}{
		{
			name:        "addAction",
			args:        hosts,
			expectedOut: "added host:  host1\nadded host:  host2\nadded host:  host3\n",
			initList:    false,
			actionFunc:  addAction,
		},
		{
			name:        "listAction",
			args:        hosts,
			expectedOut: "host1\nhost2\nhost3\n",
			initList:    true,
			actionFunc:  listAction,
		},
		{
			name:        "deleteAction",
			args:        []string{"host1", "host2"},
			expectedOut: "deleted host: host1\ndeleted host: host2\n",
			initList:    false,
			actionFunc:  deleteAction,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, teardown := setup(t, hosts, tc.initList)
			defer teardown()
			var out bytes.Buffer
			if err := tc.actionFunc(&out, f, tc.args); err != nil {
				t.Fatalf("Expected no error, got %q\n", err)
			}
			if out.String() != tc.expectedOut {
				t.Errorf("expected output %q, got %q\n", tc.expectedOut, out.String())
			}
		})
	}
}

func TestIntegration(t *testing.T) {
	hosts := []string{"host1", "host2", "host3"}
	f, teardown := setup(t, hosts, false)
	defer teardown()

	var out bytes.Buffer
	expectedOut := ""
	for _, v := range hosts {
		expectedOut += fmt.Sprintf("added host:  %s\n", v)
	}
	expectedOut += strings.Join(hosts, "\n")
	expectedOut += fmt.Sprintln()

	if err := addAction(&out, f, hosts); err != nil {
		t.Fatalf("Expected no error, got %q\n", err)
	}
	if err := listAction(&out, f, nil); err != nil {
		t.Fatalf("Expected no error, got %q\n", err)
	}

	if out.String() != expectedOut {
		t.Errorf("expected output %q but got %q\n", expectedOut, out.String())
	}

}
