package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	bundlePath string
)

var rootCmd = &cobra.Command{
	Use:   "sbun",
	Short: "Service diagnostics bundle analysis tool",
	Long: "SBun is a CLI tool which helps to analyze DC/OS service diagnostics bundle: " +
		"https://support.d2iq.com/s/article/create-service-diag-bundle",
}

func init() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error while detecting a working directory: %v\n", err.Error())
		os.Exit(1)
	}
	rootCmd.PersistentFlags().StringVarP(&bundlePath, "path", "p", wd,
		"path to the bundle directory")
}

// Execute starts Bun.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
