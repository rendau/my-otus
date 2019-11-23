package config

type Config struct {
	HttpListen string

	LogFile  string
	LogLevel string // error | warn | info | debug
}
