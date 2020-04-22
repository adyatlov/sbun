package tools

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const tasksDirName = "tasks"

// starting_20200416T110149-running_20200416T112050-killed_20200416T114052__kafka-2-broker__06e119a6-b6bb-4dae-8229-799cdf54c752
var taskIDRegexp = regexp.MustCompile(`__(.+)__(.+)$`)
var taskStatusRegexp = regexp.MustCompile(`(failed|starting|running|killed)_([0-9T]*)`)

type Task struct {
	ID              string
	Name            string
	DirName         string
	DirNameAbsolute string
	Staring         time.Time
	Running         time.Time
	Killed          time.Time
	Failed          time.Time
}

func parseTask(dirName string) (Task, error) {
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
	task.Failed = statuses["failed"]
	return task, nil
}

func parseTasks(bundlePath string) ([]Task, error) {
	tasksDir := filepath.Join(bundlePath, tasksDirName)
	taskFiles, err := ioutil.ReadDir(tasksDir)
	if err != nil {
		return nil, fmt.Errorf("cannot list files in the \"%v\" directory: %v", tasksDirName, err)
	}
	tasks := make([]Task, 0, len(taskFiles))
	for _, f := range taskFiles {
		if !f.IsDir() {
			continue
		}
		task, err := parseTask(f.Name())
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "WARNING: cannot parse the directory name \"%v\". "+
				"If you know that this directory was created by the service diagnostics bundle tool, "+
				"please, create the issue https://github.com/adyatlov/sbun/issues: %v\n",
				f.Name(), err.Error())
			continue
		}
		tasks = append(tasks, task)
	}
	if len(tasks) == 0 {
		return nil, fmt.Errorf("\"%v\" directory doesn't contain task directories", tasksDirName)
	}
	for i, _ := range tasks {
		tasks[i].DirNameAbsolute = filepath.Join(tasksDir, tasks[i].DirName)
	}
	return tasks, nil
}
