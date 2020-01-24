package cmd

import (
	"fmt"
	"github.com/rendau/my-otus/task8/scheduler/internal/adapters/logger"
	"github.com/rendau/my-otus/task8/scheduler/internal/adapters/rmq"
	"github.com/rendau/my-otus/task8/scheduler/internal/adapters/storage/pg"
	"github.com/rendau/my-otus/task8/scheduler/internal/domain/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var rootCmd = &cobra.Command{
	Use: "scheduler",
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

		stg, err := pg.NewPostgresDb(viper.GetString("pg_dsn"))
		if err != nil {
			lg.Errorw("Fail to create pg-storage", "error", err.Error())
			os.Exit(1)
		}

		mq, err := rmq.NewRmq(viper.GetString("rmq_dsn"))
		if err != nil {
			lg.Errorw("Fail to create rabbit-mq client", "error", err.Error())
			os.Exit(1)
		}

		cr := core.CreateCore(
			lg,
			stg,
			mq,
		)

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
}

func initConfig() {
	viper.SetDefault("log_level", "warn")

	viper.AutomaticEnv() // read in environment variables that match

	if viper.GetString("conf_path") != "" {
		viper.SetConfigFile(viper.GetString("conf_path"))
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
