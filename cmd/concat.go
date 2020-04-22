package cmd

import (
	"fmt"
	"os"

	"github.com/adyatlov/sbun/tools"
	"github.com/spf13/cobra"
)

func concatLogs(cmd *cobra.Command, _ []string) {
	compress := true
	if cmd.Flag("dont-compress").Changed {
		compress = false
	}
	if err := tools.Concat(bundlePath, compress); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: error when concatenating logs: %v", err.Error())
	}
}

func init() {
	concatLogsCmd := &cobra.Command{
		Use:   "concat-logs",
		Short: "Concatenate task logs to a single file",
		Long:  "Concatenate all task stdout and stderr logs to a single file: stdout_all, stderr_all",
		Run:   concatLogs,
	}
	concatLogsCmd.Flags().BoolP("dont-compress", "d", false,
		"do not compress stdout_all and stderr_all files")
	rootCmd.AddCommand(concatLogsCmd)
}
