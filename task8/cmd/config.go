package cmd

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type config struct {
	Debug      bool
	HTTPListen string

	LogFile  string
	LogLevel string // error | warn | info | debug
}

func parseConfig(configPath string) (*config, error) {
	res := &config{
		Debug:      false,
		HTTPListen: ":80",
		LogFile:    "",
		LogLevel:   "warn",
	}

	confBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if len(confBytes) > 0 {
		confFileObj := struct {
			Debug      bool   `yaml:"debug"`
			HTTPListen string `yaml:"http_listen"`
			LogFile    string `yaml:"log_file"`
			LogLevel   string `yaml:"log_level"`
		}{}

		err = yaml.Unmarshal(confBytes, &confFileObj)
		if err != nil {
			return nil, err
		}

		res.Debug = confFileObj.Debug

		if confFileObj.HTTPListen != "" {
			res.HTTPListen = confFileObj.HTTPListen
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
