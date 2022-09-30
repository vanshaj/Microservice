package todo

import (
	"fmt"
	"testing"
	"time"
)

func execution(t *testing.T, tests []struct{}, callback func(*testing.T)) {
	for _, tt := range tests {
		fmt.Println(tt)
		t.Run("", callback)
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name        string
		taskName    string
		done        bool
		createdAt   time.Time
		completedAt time.Time
	}{
		{"SimpleTask", "FirstTask", true, time.Now(), time.Now().Add(1 * time.Hour)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := List{}
			l.Add(tt.name)

			if l[0].Task != tt.name {
				t.Errorf("Expected %s, got %s instead", tt.name, l[0].Task)
			}
		})
	}
}

func TestComplete(t *testing.T) {
	tests := []struct {
		name     string
		taskName string
	}{
		{"simpleComplete", "First Completed Task"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := List{}
			l.Add(tt.taskName)
			l.Complete(tt.taskName)

			if l[0].Done != true {
				t.Errorf("Task didnot get completed")
			}
		})
	}
}

func TestSaveGet(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{"TestTheSaveAndGet", "output.json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//defer os.Remove(tt.filename)
			l := List{}
			l.Add("FirstTask")
			l.Add("SecondTask")
			l.Add("ThirdTask")

			if err := l.Save(tt.filename); err != nil {
				t.Fatalf("Error saving list to file: %s", err)
			}

			l2 := List{}
			if err := l2.Get(tt.filename); err != nil {
				t.Fatalf("Error getting file to list: %s", err)
			}

			if l[0].Task != l2[0].Task {
				t.Errorf("task %s from file must match %s from list", l[0].Task, l2[0].Task)
			}

		})
	}
}
