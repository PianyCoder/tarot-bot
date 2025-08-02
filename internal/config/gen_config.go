package config

type GenConfig struct {
	GenAPIToken     string `env:"GEN_API_TOKEN,required"`
	GenAPINetworkID string `env:"GEN_API_NETWORK_ID"`
}
