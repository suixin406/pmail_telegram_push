module github.com/ydzydzydz/pmail_telegram_push

go 1.25.1

replace github.com/Jinnrry/pmail v2.8.7+incompatible => github.com/Jinnrry/pmail/server v0.0.0-20250830122511-778e279b4fae

require (
	github.com/Jinnrry/pmail v2.8.7+incompatible
	github.com/go-telegram/bot v1.17.0
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/net v0.39.0
)

require (
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/emersion/go-message v0.18.2 // indirect
	github.com/emersion/go-msgauth v0.6.8 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/microcosm-cc/bluemonday v1.0.27 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
)
