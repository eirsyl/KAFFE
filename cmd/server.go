package cmd

import "github.com/spf13/cobra"

func init() {
	RootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "webserver with UI and API",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
