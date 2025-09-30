package config

import (
	"encoding/json"
	"os"

	pconfig "github.com/Jinnrry/pmail/config"
)

const (
	MAIN_CONFIG_FILE   = "./config/config.json"
	PLUGIN_CONFIG_FILE = "./config/pmail_telegram_push_config.json"
)

type PluginConfig struct {
	TelegramBotToken string `json:"telegram_bot_token"`
	Proxy            string `json:"proxy" default:""`
	Timeout          int    `json:"timeout" default:"30"`
	Debug            bool   `json:"debug" default:"false"`
}

type Config struct {
	PluginConfig *PluginConfig
	MainConfig   *pconfig.Config
}

func readMainConfig() *pconfig.Config {
	content, err := os.ReadFile(MAIN_CONFIG_FILE)
	if err != nil {
		panic(err)
	}
	var mainConfig pconfig.Config
	if err := json.Unmarshal(content, &mainConfig); err != nil {
		panic(err)
	}
	return &mainConfig
}

func readPluginConfig() *PluginConfig {
	content, err := os.ReadFile(PLUGIN_CONFIG_FILE)
	if err != nil {
		panic(err)
	}
	var pluginConfig PluginConfig
	if err := json.Unmarshal(content, &pluginConfig); err != nil {
		panic(err)
	}
	if pluginConfig.TelegramBotToken == "" {
		panic("telegram bot token is empty")
	}
	return &pluginConfig
}

func ReadConfig() *Config {
	return &Config{
		PluginConfig: readPluginConfig(),
		MainConfig:   readMainConfig(),
	}
}
