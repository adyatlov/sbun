package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

// starting_20200416T110149-running_20200416T112050-killed_20200416T114052__kafka-2-broker__06e119a6-b6bb-4dae-8229-799cdf54c752
var taskIDRegexp = regexp.MustCompile(`__(.+)__(.+)$`)
var taskStatusRegexp = regexp.MustCompile(`(failed|starting|running|killed)_([0-9T]*)`)

type Task struct {
	ID      string
	Name    string
	DirName string
	Staring time.Time
	Running time.Time
	Killed  time.Time
}

func (t Task) IsKilled() bool {
	return t.Killed.IsZero()
}

func ParseTask(dirName string) (Task, error) {
	task := Task{}
	idTokens := taskIDRegexp.FindStringSubmatch(dirName)
	if len(idTokens) != 3 {
		return task, fmt.Errorf("cannot parse ID and name for task: %v", dirName)
	}
	task.ID = idTokens[2]
	task.Name = idTokens[1]
	statusTokens := taskStatusRegexp.FindAllStringSubmatch(dirName, -1)
	if len(statusTokens) == 0 {
		return task, fmt.Errorf("cannot parse statuses for task: %v", dirName)
	}
	statuses := make(map[string]time.Time)
	for _, token := range statusTokens {
		// Mon Jan 2 15:04:05 -0700 MST 2006
		t, err := time.Parse("20060102T150405", token[2])
		if err != nil {
			return task, err
		}
		statuses[token[1]] = t
	}
	task.DirName = dirName
	task.Staring = statuses["starting"]
	task.Running = statuses["running"]
	task.Killed = statuses["killed"]
	return task, nil
}

func main() {
	taskFiles, err := ioutil.ReadDir("./tasks")
	if err != nil {
		panic(err)
	}
	tasks := make([]Task, 0, len(taskFiles))
	for _, f := range taskFiles {
		if !f.IsDir() {
			continue
		}
		task, err := ParseTask(f.Name())
		if err != nil {
			panic(err.Error())
		}
		tasks = append(tasks, task)
	}
	csvWriter := csv.NewWriter(os.Stdout)
	for _, t := range tasks {
		csvWriter.Write([]string{t.Name,
			printTime(t.Staring),
			printTime(t.Running),
			printTime(t.Killed),
			t.ID,
			t.DirName})
	}
	csvWriter.Flush()
}

func printTime(t time.Time) string {
	if t.IsZero() {
		return "N/A"
	}
	return t.String()
}
