package internal

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type fullBackupCommand struct {
	Config struct {
		Debug       bool
		SaveTo      string
		TrelloKey   string
		TrelloToken string
	}
	cmd *cobra.Command
}

func NewFullBackupCommand() *cobra.Command {
	command := &fullBackupCommand{}

	command.cmd = &cobra.Command{
		Use:   "full-backup",
		Short: "Fetches all the information of a given board and saves it locally",
		Long:  `Use it to perform full-backups, it will bring all information of a given Trello board and make it available locally. After the first execution, consider using incremental backups instead.`,
		Run: func(cmd *cobra.Command, args []string) {
			if command.Config.Debug {
				logger.SetLevel(logrus.DebugLevel)
			}

			BackupAllBoards(
				NewTrello(command.Config.TrelloKey, command.Config.TrelloToken),
				command.Config.SaveTo,
			)

			fmt.Println(command.Config.TrelloKey)
			fmt.Println(command.Config.TrelloToken)
		},
	}

	saveTo, err := os.Getwd()

	if err != nil {
		logger.WithError(err).WithField("directory", saveTo).Fatalf("failed to get current directory")
	}

	saveTo = fmt.Sprint(saveTo, fmt.Sprintf("%c", os.PathSeparator), "go-trello-backup")

	command.cmd.Flags().BoolVarP(&command.Config.Debug, "debug", "d", false, "Shows additional information that might be useful for debugging the application")
	command.cmd.Flags().StringVarP(&command.Config.SaveTo, "save-to", "b", saveTo, "Directory in which the backup files will be saved to")
	command.cmd.Flags().StringVarP(&command.Config.TrelloKey, "trello-api-key", "k", "", "API key that will be used in order to communicate with Trello REST API")
	command.cmd.Flags().StringVarP(&command.Config.TrelloToken, "trello-api-token", "t", "", "API token that will be used in order to communicate with Trello REST API")

	return command.cmd
}
