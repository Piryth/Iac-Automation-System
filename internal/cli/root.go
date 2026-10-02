package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "iascli",
	Short: "CLI tool for Ias automation plateform",
	Long:  `IasCli manages an IaS system, from workflow deployment to management`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to iascli v1.0.0. Run iascli --help for usage")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
