package db

import (
	"time"

	"github.com/Jinnrry/pmail/config"
	"github.com/Jinnrry/pmail/utils/errors"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/ydzydzydz/pmail_telegram_push/model"
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

var Instance *xorm.Engine

func Init(config *config.Config) error {
	dsn := config.DbDSN
	var err error

	switch config.DbType {
	case "mysql":
		Instance, err = xorm.NewEngine("mysql", dsn)
		Instance.SetMaxOpenConns(100)
		Instance.SetMaxIdleConns(10)
	case "sqlite":
		Instance, err = xorm.NewEngine("sqlite", dsn)
		Instance.SetMaxOpenConns(1)
		Instance.SetMaxIdleConns(1)
	case "postgres":
		Instance, err = xorm.NewEngine("postgres", dsn)
		Instance.SetMaxOpenConns(100)
		Instance.SetMaxIdleConns(10)
	default:
		return errors.New("Database Type Error!")
	}
	if err != nil {
		log.Errorf("DB init Error! %s", err.Error())
		return errors.Wrap(err)
	}

	Instance.SetConnMaxLifetime(30 * time.Minute)
	Instance.ShowSQL(false)

	Instance.Sync2(new(model.TelegramPushSetting))
	return nil
}
