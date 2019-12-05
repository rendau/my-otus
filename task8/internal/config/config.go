package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Debug      bool
	HTTPListen string
	GRPCListen string

	PgDsn            string
	PgMigrationsPath string

	LogFile  string
	LogLevel string // error | warn | info | debug
}

func ParseConfig(path string) (*Config, error) {
	res := &Config{
		Debug:      false,
		HTTPListen: ":80",
		LogFile:    "",
		LogLevel:   "warn",
	}

	confBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(confBytes) > 0 {
		confFileObj := struct {
			Debug            bool   `yaml:"debug"`
			HTTPListen       string `yaml:"http_listen"`
			GRPCListen       string `yaml:"grpc_listen"`
			PgDsn            string `yaml:"pg_dsn"`
			PgMigrationsPath string `yaml:"pg_migrations_path"`
			LogFile          string `yaml:"log_file"`
			LogLevel         string `yaml:"log_level"`
		}{}

		err = yaml.Unmarshal(confBytes, &confFileObj)
		if err != nil {
			return nil, err
		}

		res.Debug = confFileObj.Debug

		if confFileObj.HTTPListen != "" {
			res.HTTPListen = confFileObj.HTTPListen
		}

		if confFileObj.GRPCListen != "" {
			res.GRPCListen = confFileObj.GRPCListen
		}

		if confFileObj.PgDsn != "" {
			res.PgDsn = confFileObj.PgDsn
		}

		if confFileObj.PgMigrationsPath != "" {
			res.PgMigrationsPath = confFileObj.PgMigrationsPath
		}

		if confFileObj.LogFile != "" {
			res.LogFile = confFileObj.LogFile
		}

		if confFileObj.LogLevel != "" {
			res.LogLevel = confFileObj.LogLevel
		}
	}

	return res, nil
}
