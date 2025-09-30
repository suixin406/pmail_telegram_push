package main

import (
	"github.com/Jinnrry/pmail/hooks/framework"
	"github.com/ydzydzydz/pmail_telegram_push/config"
	"github.com/ydzydzydz/pmail_telegram_push/db"
	"github.com/ydzydzydz/pmail_telegram_push/hook"
)

func main() {
	config := config.ReadConfig()
	db.Init(config.MainConfig)
	framework.CreatePlugin(
		hook.PLUGIN_NAME,
		hook.NewPmailTelegramPushHook(config),
	).Run()
}
