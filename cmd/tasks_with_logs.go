package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adyatlov/sbun/tools"
	"github.com/spf13/cobra"
)

func tasksWithLogs(cmd *cobra.Command, _ []string) {
	f := cmd.Flag("save-to")
	oldDir := filepath.Join("..", tools.DirNameTasks)
	newDir := filepath.Join(bundlePath, "tasks_with_logs")
	if f.Changed {
		oldDir = filepath.Join(bundlePath, tools.DirNameTasks)
		newDir = filepath.Join(f.Value.String(), "tasks_with_logs")
	}
	tasks, err := tools.FindTasks(bundlePath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot find tasks: %v", err)
		return
	}
	if len(tasks) == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "No tasks with logs found.")
		return
	}
	fmt.Printf("DEBUG: %v\n", newDir)
	if err := os.RemoveAll(newDir); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot remove directory: %v\n", err)
	}
	if err := os.MkdirAll(newDir, 0777); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot create directory: %v\n", err)
	}
	for _, task := range tasks {
		if !task.HasLogs {
			continue
		}
		oldDirName := filepath.Join(oldDir, task.DirName)
		newDirName := filepath.Join(newDir, task.DirName)
		if err := os.Symlink(oldDirName, newDirName); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Cannot create link: %v", err)
		}
	}
}

func init() {
	taskCsvCmd := &cobra.Command{
		Use:   "tasks-with-logs",
		Short: "Find tasks which have logs",
		Long: "The command creates a tasks_with_logs directory and a sym-link in this directory to each task which has logs. " +
			"By default, it create links with relative paths in <bundle path>/tasks_with_logs.",
		Run: tasksWithLogs,
	}
	taskCsvCmd.Flags().StringP("save-to", "s", "",
		"path to the directory in which the command creates a tasks_with_logs direcoty")
	rootCmd.AddCommand(taskCsvCmd)
}
