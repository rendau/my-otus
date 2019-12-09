package cmd

import (
	"fmt"
	"github.com/rendau/my-otus/task8/sender/internal/adapters/logger"
	"github.com/rendau/my-otus/task8/sender/internal/adapters/rmq"
	"github.com/rendau/my-otus/task8/sender/internal/domain/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use: "sender",
	Run: func(cmd *cobra.Command, args []string) {
		lg, err := logger.NewLogger(
			"",
			viper.GetString("log_level"),
			viper.GetBool("debug"),
			false,
		)
		if err != nil {
			log.Fatalln("Fail to create logger:", err.Error())
		}

		mq, err := rmq.NewRmq(viper.GetString("rmq_dsn"))
		if err != nil {
			lg.Errorw("Fail to create rabbit-mq client", "error", err.Error())
			os.Exit(1)
		}

		cr := core.CreateCore(lg, mq)

		cr.Run()

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		var exitCode int

		select {
		case <-cr.Done():
			exitCode = 1
			break
		case <-stop:
			break
		}

		cr.Stop()

		err = mq.Stop()
		if err != nil {
			lg.Error("Fail to stop mq", "error", err.Error())
			exitCode = 1
		}

		os.Exit(exitCode)
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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

func initConfig() {
	viper.SetDefault("log_level", "warn")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
