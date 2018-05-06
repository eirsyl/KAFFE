package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/webkom/KAFFE/pkg"
)

var RootCmd = &cobra.Command{
	Use:     "kaffe",
	Short:   "Moccamaster tracker",
	Version: fmt.Sprintf("%s %s", pkg.Version, pkg.BuildDate),
	Long:    `KAFFE tracks the usage of the modified Moccamaster owned by Webkom.`,
}
