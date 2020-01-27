package cmd

import (
	"context"
	"fmt"
	"github.com/rendau/my-otus/task8/api/internal/adapters/grpc"
	"github.com/rendau/my-otus/task8/api/internal/adapters/logger"
	"github.com/rendau/my-otus/task8/api/internal/adapters/rest"
	"github.com/rendau/my-otus/task8/api/internal/adapters/storage/pg"
	"github.com/rendau/my-otus/task8/api/internal/domain/usecases"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var rootCmd = &cobra.Command{
	Use: "calendar_api",
	Run: func(cmd *cobra.Command, args []string) {
		lg, err := logger.NewLogger(
			viper.GetString("log_file"),
			viper.GetString("log_level"),
			viper.GetBool("debug"),
			false,
		)
		if err != nil {
			log.Fatalln("Fail to create logger")
		}
		defer lg.Sync()

		db, err := pg.NewPostgresDb(viper.GetString("pg_dsn"))
		if err != nil {
			lg.Fatalw("Fail to create postgres-db", "error", err)
		}

		ucs := usecases.CreateUsecases(lg, db)

		restAPI := rest.CreateAPI(
			lg,
			viper.GetString("http_listen"),
			viper.GetString("m_http_listen"),
			ucs,
		)
		restAPI.Start()

		grpcAPI := grpc.CreateAPI(lg, viper.GetString("grpc_listen"), ucs)
		grpcAPI.Start()

		lg.Infof("Started listen, address: %s", viper.GetString("http_listen"))

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop

		lg.Info("Shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() {
			cancel()
		}()
		err = restAPI.Shutdown(ctx)
		if err != nil {
			lg.Fatalw("Fail to shutdown rest-api", "error", err)
		}

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
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetDefault("debug", "false")
	viper.SetDefault("http_listen", ":80")
	viper.SetDefault("grpc_listen", ":8080")
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
