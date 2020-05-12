package main

import (
	"github.com/t-xinlin/doc/pkg/cli/cmd"
)

func main() {
	cmd.RootCmd.SetArgs([]string{"version", "1000"})
	cmd.RootCmd.Execute()
	cmd.RootCmd.SetArgs([]string{"add", "1000", "51"})
	cmd.RootCmd.Execute()
	cmd.RootCmd.SetArgs([]string{"--help", "1000"})
	cmd.RootCmd.Execute()
}
