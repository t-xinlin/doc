package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	VersionCmd.Flags().StringP("user", "u", "username", "usage")// add falg
	RootCmd.AddCommand(VersionCmd)
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("User name: ",cmd.Flag("user").Name) // get from flag
		if len(args) > 0 {
			fmt.Println("Version Static Site Generator v0.9 -- HEAD ", args[0])
		}
	},
}
