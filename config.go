package ivy

import "github.com/plimble/envconfig"

var (
	config *Config
)

type Config struct {
	Addr              string `default:":20000" required:"true"`
	SourceAwsId       string
	SourceAwsSecret   string
	SourceAwsS3Bucket string
	SourceAwsS3Region string
}

func GetConfig() (*Config, error) {
	config := &Config{}
	err := envconfig.Process("ivy", config)

	return config, err
}
