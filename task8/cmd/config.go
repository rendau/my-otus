package cmd

import (
	"github.com/rendau/my-otus/task8/internal/config"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func parseConfig(configPath string) (*config.Config, error) {
	res := &config.Config{
		HttpListen: ":80",
		LogFile:    "./log.log",
		LogLevel:   "warn",
	}

	confBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if len(confBytes) > 0 {
		confFileObj := struct {
			HttpListen string `yaml:"http_listen"`
			LogFile    string `yaml:"log_file"`
			LogLevel   string `yaml:"log_level"`
		}{}

		err = yaml.Unmarshal(confBytes, &confFileObj)
		if err != nil {
			return nil, err
		}

		if confFileObj.HttpListen != "" {
			res.HttpListen = confFileObj.HttpListen
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
