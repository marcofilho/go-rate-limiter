package configs

import "github.com/spf13/viper"

type Config struct {
	RateLimiterMaxRequests   int    `mapstructure:"RATE_LIMITER_MAX_REQUESTS"`
	RateLimiterBlockDuration int    `mapstructure:"RATE_LIMITER_BLOCK_DURATION"`
	RateLimiterType          string `mapstructure:"RATE_LIMITER_TYPE"`
	RedisAddress             string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword            string `mapstructure:"REDIS_PASSWORD"`
	RedisDB                  int    `mapstructure:"REDIS_DB"`
	WebServerPort            string `mapstructure:"WEB_SERVER_PORT"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
