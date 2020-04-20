package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/adyatlov/sbun/taskcsv"
)

func printTasks(cmd *cobra.Command, args []string) {
	writer := os.Stdout
	o := cmd.Flag("output")
	if o.Changed {
		var err error
		writer, err = os.Create(o.Value.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create file: %v", err.Error())
			os.Exit(1)
		}
	}
	err := taskcsv.WriteCsv(bundlePath, writer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
	}
}

func init() {
	taskcsvCmd := &cobra.Command{
		Use:   "taskcsv",
		Short: "Print service task list",
		Long: "Print service task list in the CSV format to the standard output or file. The order of columns is: " +
			"<task name>, <starting timestamp>, <running timestamp>, <killed timestamp>, <failed timestamp>, <task ID>, " +
			"<path to the task directory>",
		Run: printTasks,
	}
	taskcsvCmd.Flags().StringP("output", "o", "",
		"path to the output CSV file")
	rootCmd.AddCommand(taskcsvCmd)
}
