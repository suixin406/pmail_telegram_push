package main

import (
	"github.com/Jinnrry/pmail/hooks/framework"
	"github.com/ydzydzydz/pmail_telegram_push/config"
	"github.com/ydzydzydz/pmail_telegram_push/hook"
)

func main() {
	config := config.ReadConfig()
	hook := hook.NewPmailTelegramPushHook(config)
	framework.CreatePlugin("telegram_push", hook).Run()
}
