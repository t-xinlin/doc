package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

func init() {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "a + b",
		Long:  "add long",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 2 {
				a, err := strconv.Atoi(args[0])
				if err != nil {
					fmt.Println("para error")
					return

				}
				b, err := strconv.Atoi(args[1])
				if err != nil {
					fmt.Println("para error")
					return

				}
				fmt.Println("Result: ", a+b)
			} else {
				fmt.Println("para error")
			}

		},
	}
	cmd.PersistentFlags().Int("a", 0, "a")
	cmd.PersistentFlags().Int("b", 0, "b")
	RootCmd.AddCommand(cmd)
}
