package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd =&cobra.Command{
	Use:"hugo",
	Short:"Hugo",
	Long:"It is a good",
	Run: func(cmd *cobra.Command, args []string) {
		
	},
}

func Execute()  {
	if err := rootCmd.Execute(); err !=nil{
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
