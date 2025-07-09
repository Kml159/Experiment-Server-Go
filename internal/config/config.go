package config

type Config struct {
	ServerAddress                   string
	ThreadAmount                    int
	ClientSendUpdateStatusInSeconds int
	ReceivedOutputFilePath          string
	ExperimentBaseId                int
	SubtractCompletedExperiments    bool
	ExperimentDuplicate             int
	ProductGenerationPopulation     int
}

func Load() *Config {
	return &Config{
		ThreadAmount:                    4,
		ClientSendUpdateStatusInSeconds: 60 * 10,
		ReceivedOutputFilePath:          "received_output",
		ExperimentBaseId:                39168127,
		SubtractCompletedExperiments:    true,
		ExperimentDuplicate:             20,
		ProductGenerationPopulation:     1e6,
	}
}
