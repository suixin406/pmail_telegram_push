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
	TelegramChatID   string `json:"telegram_chat_id"`
	ShowContent      bool   `json:"show_content" default:"true"`
	SpoilerContent   bool   `json:"spoiler_content" default:"true"`
	SendAttachment   bool   `json:"send_attachment" default:"true"`
	Debug            bool   `json:"debug" default:"false"`
	Proxy            string `json:"proxy" default:""`
	Timeout          int    `json:"timeout" default:"30"`
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
	if pluginConfig.TelegramBotToken == "" || pluginConfig.TelegramChatID == "" {
		panic("telegram bot token or chat id is empty")
	}
	return &pluginConfig
}

func ReadConfig() *Config {
	return &Config{
		PluginConfig: readPluginConfig(),
		MainConfig:   readMainConfig(),
	}
}
