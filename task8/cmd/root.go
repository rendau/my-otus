package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "clncnd",
	Short: "CleanCalendar is a calendar micorservice demo",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello world!")
	},
}
