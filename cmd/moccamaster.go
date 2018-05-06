package cmd

import "github.com/spf13/cobra"

func init() {
	RootCmd.AddCommand(moccamasterCmd)
}

var moccamasterCmd = &cobra.Command{
	Use:   "moccamaster",
	Short: "moccamaster agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
