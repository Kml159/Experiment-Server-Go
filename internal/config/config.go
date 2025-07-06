package config

type Config struct {
	ServerAddress 					string
	ThreadAmount  					int
	ClientSendUpdateStatusInSeconds int
}

func Load() *Config {
	return &Config{
		ServerAddress: "http://localhost:8080",
		ThreadAmount:  4,
		ClientSendUpdateStatusInSeconds: 60*10,
	}
}
