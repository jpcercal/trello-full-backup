package main

import (
	"github.com/spf13/cobra"

	app "github.com/jpcercal/go-trello-backup/internal"
)

func main() {
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(app.NewFullBackupCommand())

	if err := rootCmd.Execute(); err != nil {
		app.GetLogger().WithError(err).Fatalf("failed to execute command")
	}
}
