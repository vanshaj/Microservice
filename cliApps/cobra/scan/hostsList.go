package scan

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
)

var (
	ErrEmptyHost  = errors.New("Empty Host")
	ErrNoSuchHost = errors.New("No such Host")
	ErrNoSort     = errors.New("unable to sort")
)

type HostsList struct {
	Hosts []string
}

func (h *HostsList) Add(host string) error {
	if host == "" {
		return ErrEmptyHost
	} else {
		h.Hosts = append(h.Hosts, host)
	}
	return nil
}

func (h *HostsList) Remove(host string) error {
	index, found := h.Search(host)
	if !found {
		return ErrNoSuchHost
	} else {
		h.Hosts = append(h.Hosts[:index], h.Hosts[index+1:]...)
		return nil
	}
}

func (h *HostsList) Search(host string) (int, bool) {
	if !sort.StringsAreSorted(h.Hosts) {
		sort.Strings(h.Hosts)
	}
	index := sort.SearchStrings(h.Hosts, host)
	if index < len(h.Hosts) && h.Hosts[index] == host {
		return index, true
	}
	return -1, false
}

func (h *HostsList) Load(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if err := h.Add(scanner.Text()); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (h *HostsList) Save(w io.Writer) error {
	for _, val := range h.Hosts {
		_, err := w.Write([]byte(val + "\n"))
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *HostsList) List() error {
	for _, val := range h.Hosts {
		fmt.Fprint(os.Stdout, val)
	}
	return nil
}
