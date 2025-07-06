package config

type Config struct {
	ServerAddress 					string
	ThreadAmount  					int
	ClientSendUpdateStatusInSeconds int
}

func Load() *Config {
	return &Config{
		ThreadAmount:  4,
		ClientSendUpdateStatusInSeconds: 60*10,
	}
}
