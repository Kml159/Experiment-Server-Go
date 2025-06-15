package config

type Config struct {
	ServerAddress string
	ThreadAmount  int
}

func Load() *Config {
	return &Config{
		ServerAddress: "http://localhost:8080",
		ThreadAmount:  4,
	}
}
