package cmd

import (
	"context"
	"github.com/rendau/my-otus/task8/internal/adapters/grpc"
	"github.com/rendau/my-otus/task8/internal/adapters/logger"
	"github.com/rendau/my-otus/task8/internal/adapters/rest"
	"github.com/rendau/my-otus/task8/internal/adapters/storage/pg"
	"github.com/rendau/my-otus/task8/internal/config"
	"github.com/rendau/my-otus/task8/internal/domain/usecases"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serves calendar service",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.ParseConfig(cfgFile)
		if err != nil {
			log.Fatalln("Fail to parse config")
		}

		lg, err := logger.NewLogger(cfg.LogFile, cfg.LogLevel, cfg.Debug, false)
		if err != nil {
			log.Fatalln("Fail to create logger")
		}
		defer lg.Sync()

		err = pg.MigrationDo(cfg.PgDsn, cfg.PgMigrationsPath, "up")
		if err != nil {
			lg.Fatalw("Fail to apply migrations", "error", err)
		}

		db, err := pg.NewPostgresDb(cfg.PgDsn)
		if err != nil {
			lg.Fatalw("Fail to create postgres-db", "error", err)
		}

		ucs := usecases.CreateUsecases(lg, db)

		restAPI := rest.CreateAPI(lg, cfg.HTTPListen, ucs)
		restAPI.Start()

		grpcAPI := grpc.CreateAPI(lg, cfg.GRPCListen, ucs)
		grpcAPI.Start()

		lg.Infof("Started listen, address: %s", cfg.HTTPListen)

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

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVar(&cfgFile, "config", "", "config file eg. conf.yml")
	_ = serveCmd.MarkFlagRequired("config")
}
