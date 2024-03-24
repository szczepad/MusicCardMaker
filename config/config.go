package config

import "github.com/spf13/viper"

type Configuration struct {
	Spotify SpotifyConfiguration
}

type SpotifyConfiguration struct {
	ClientID     string
	ClientSecret string
}

func CreateConfig() Configuration {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	var configuration Configuration

	err := viper.ReadInConfig()
	if err != nil {
		panic("Could not read config")
	}

	err = viper.Unmarshal(&configuration)
	if err != nil {
		panic("Could not unmarshal config")
	}

	return configuration
}
