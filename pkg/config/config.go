package config

import (
	"github.com/spf13/viper"
)

type Point struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

type Bounds struct {
	BottomLeft Point `json:"bottomLeft"`
	TopRight   Point `json:"topRight"`
}

type Filters struct {
	Banks      []string `json:"banks"`
	Currencies []string `json:"currencies"`
}

type TRequest struct {
	Bounds  Bounds  `json:"bounds"`
	Filters Filters `json:"filters"`
	Zoom    int     `json:"zoom"`
}

type Limit struct {
	Currency string `json:"currency"`
	Amount   int    `json:"amount"`
}

type Config struct {
	Redis string `json:"redis"`

	Telegram struct {
		Token string `mapstructure:"token"`
		Chat  int64  `mapstructure:"chat"`
	} `mapstructure:"telegram"`

	Bounds struct {
		BottomLeft Point `mapstructure:"bottom_left"`
		TopRight   Point `mapstructure:"top_right"`
	} `mapstructure:"bounds"`

	Zoom int `mapstructure:"zoom"`
}

func New(configFile string) (*Config, error) {
	var cfg Config
	if err := LoadConfig(configFile, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func LoadConfig(configFile string, config interface{}) error {

	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return err
	}

	return err
}
