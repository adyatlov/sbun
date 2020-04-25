package tools

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const DirNameTasks = "tasks"

// starting_20200416T110149-running_20200416T112050-killed_20200416T114052__kafka-2-broker__06e119a6-b6bb-4dae-8229-799cdf54c752
var taskIDRegexp = regexp.MustCompile(`__(.+)__(.+)$`)
var taskStatusRegexp = regexp.MustCompile(`(failed|starting|running|killed)_([0-9T]*)`)

// stdout.1.gz, stdout.gz, stdout, stdout.1
var stdoutRegexp = regexp.MustCompile(`^stdout(\.[0-9]+)?(\.gz)?$`)
var stderrRegexp = regexp.MustCompile(`^stderr(\.[0-9]*)?(\.gz)?$`)

// stdout_all, stderr_all, stdout_all.gz, stderr_all.gz
var stdAllRegexp = regexp.MustCompile(`^(stderr|stdout)_all(\.gz)?$`)

type Task struct {
	ID              string
	Name            string
	DirName         string
	DirNameAbsolute string
	Staring         time.Time
	Running         time.Time
	Killed          time.Time
	Failed          time.Time
	HasLogs         bool
}

func parseTaskDirName(dirName string) (Task, error) {
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

func FindTasks(bundlePath string) ([]Task, error) {
	tasksDir := filepath.Join(bundlePath, DirNameTasks)
	taskFiles, err := ioutil.ReadDir(tasksDir)
	if err != nil {
		return nil, fmt.Errorf("cannot list files in the \"%v\" directory: %v", DirNameTasks, err)
	}
	tasks := make([]Task, 0, len(taskFiles))
	for _, f := range taskFiles {
		if !f.IsDir() {
			continue
		}
		task, err := parseTaskDirName(f.Name())
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "WARNING: cannot parse the directory name \"%v\". "+
				"If you know that this directory was created by the service diagnostics bundle tool, "+
				"please, create the issue https://github.com/adyatlov/sbun/issues: %v\n",
				f.Name(), err.Error())
			continue
		}
		task.DirNameAbsolute = filepath.Join(tasksDir, task.DirName)
		task.HasLogs = hasLogs(task.DirNameAbsolute)
		tasks = append(tasks, task)
	}
	if len(tasks) == 0 {
		return nil, fmt.Errorf("\"%v\" directory doesn't contain task directories", DirNameTasks)
	}
	return tasks, nil
}

func hasLogs(taskDir string) bool {
	var has bool
	err := filepath.Walk(taskDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Cannot walk into path %v: %v", path, err)
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if stdoutRegexp.MatchString(info.Name()) ||
			stderrRegexp.MatchString(info.Name()) ||
			stdAllRegexp.MatchString(info.Name()) {
			has = true
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot walk: %v", err)
	}
	return has
}
