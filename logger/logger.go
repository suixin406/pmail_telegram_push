package logger

import "github.com/phuslu/log"

var BotLogger = log.Logger{
	Level:  log.DebugLevel,
	Caller: 1,
	Writer: &log.ConsoleWriter{
		ColorOutput:    true,
		EndWithMessage: true,
	},
}

var PluginLogger = log.Logger{
	Level:  log.InfoLevel,
	Caller: 1,
	Writer: &log.ConsoleWriter{
		ColorOutput:    true,
		EndWithMessage: true,
	},
}
