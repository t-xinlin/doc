package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd =&cobra.Command{
	Use:"cmd",
	Short:"Cmd",
	Long:"It is a good",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute()  {
	if err := RootCmd.Execute(); err !=nil{
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
