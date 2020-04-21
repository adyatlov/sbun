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
	O := cmd.Flag("default-name")
	if o.Changed && O.Changed {
		fmt.Fprintln(os.Stderr, "ERROR: Flags -o (--output) and -O (--default-name) are mutually exclusive. "+
			"Please use only one of them.")
		os.Exit(1)
	}
	var err error
	if o.Changed {
		writer, err = os.Create(o.Value.String())
	}
	if O.Changed {
		writer, err = os.Create("tasks.csv")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Cannot create file: %v", err.Error())
		os.Exit(1)
	}
	err = taskcsv.WriteCsv(bundlePath, writer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err.Error())
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
	taskcsvCmd.Flags().BoolP("default-name", "O", false,
		"write output to the tasks.csv file")
	rootCmd.AddCommand(taskcsvCmd)
}
