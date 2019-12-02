package cmd

import (
	"context"
	"fmt"
	"github.com/rendau/my-otus/task10/internal"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var timeout uint32

var rootCmd = &cobra.Command{
	Use: "go-telnet [flags] host port. Example: go-telnet --timeout=3 localhost 5050",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			err := cmd.Usage()
			if err != nil {
				log.Fatalln(err.Error())
			}
			os.Exit(1)
		}

		host := args[0]
		port, _ := strconv.ParseInt(args[1], 10, 64)
		if port == 0 {
			log.Fatalln("Bad port, port must be integer number")
		}

		client := internal.NewClient()

		startCtx, _ := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

		doneCh, err := client.Start(startCtx, host, port, os.Stdin)
		if err != nil {
			log.Fatalln(err)
		}

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

		select {
		case <-stop:
		case <-doneCh:
		}

		fmt.Println("Stopping")

		client.Stop()

		os.Exit(0)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().Uint32VarP(&timeout, "timeout", "t", 10, "Timeout in seconds, default 10")
}
