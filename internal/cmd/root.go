package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd := &cobra.Command{
		Use:     "rosseti-curier",
		Short:   "rosseti-curier is a service for rosseti-app",
		Version: AppInfo,
	}

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Command for starting HTTP server and connection with third-party services",
		Run:   serveHandler,
	}

	rootCmd.AddCommand(serveCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
