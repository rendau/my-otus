/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"github.com/rendau/my-otus/task8/internal/adapters/grpc"
	"github.com/rendau/my-otus/task8/internal/adapters/logger"
	"github.com/rendau/my-otus/task8/internal/adapters/rest"
	"github.com/rendau/my-otus/task8/internal/adapters/storage/memdb"
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
		cfg, err := parseConfig(cfgFile)
		if err != nil {
			log.Fatalln("Fail to parse config")
		}

		//fmt.Println(cfg.LogFile)
		lg, err := logger.NewLogger(cfg.LogFile, cfg.LogLevel, cfg.Debug, false)
		if err != nil {
			log.Fatalln("Fail to create logger")
		}
		defer lg.Sync()

		db, err := memdb.NewMemDb(lg)
		if err != nil {
			lg.Fatalw("Fail to create mem-db", "error", err)
		}

		ucs := usecases.CreateUsecases(lg, db)

		restAPI := rest.CreateAPI(lg, cfg.HTTPListen, ucs)
		restAPI.Start()

		grpcAPI := grpc.CreateAPI(lg, cfg.GRPCListen, ucs)
		grpcAPI.Start()

		time.Sleep(100 * time.Millisecond)

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
