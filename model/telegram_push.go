package model

import (
	"errors"

	"xorm.io/xorm"
)

type TelegramPushSetting struct {
	ID                 int    `xorm:"id pk autoincr comment('主键')" json:"-"`
	UserID             int    `xorm:"user_id int index('idx_uid') comment('用户id') unique('idx_uid')" json:"-"`
	ChatID             string `xorm:"chat_id varchar(255) index('idx_cid') index comment('聊天id')" json:"chat_id"`
	ShowContent        bool   `xorm:"show_content bool not null default 1 comment('是否显示邮件内容')" json:"show_content"`
	SpoilerContent     bool   `xorm:"spoiler_content bool not null default 1 comment('是否 spoiler 显示邮件内容')" json:"spoiler_content"`
	SendAttachments    bool   `xorm:"send_attachments bool not null default 1 comment('是否发送附件')" json:"send_attachments"`
	DisableLinkPreview bool   `xorm:"disable_link_preview bool not null default 1 comment('是否禁用链接预览')" json:"disable_link_preview"`
}

func (u *TelegramPushSetting) TableName() string {
	return "plugin_telegram_push_setting"
}

func GetSetting(db *xorm.Engine, userID int) (setting *TelegramPushSetting, err error) {
	setting = &TelegramPushSetting{}
	has, err := db.Where("user_id = ?", userID).Get(setting)
	if err != nil {
		return nil, err
	}
	if !has {
		if err = CreateDefaultSetting(db, userID); err != nil {
			return nil, err
		}
	}
	return setting, nil
}

func CreateDefaultSetting(db *xorm.Engine, userID int) (err error) {
	setting := &TelegramPushSetting{
		UserID:             userID,
		ChatID:             "",
		ShowContent:        true,
		SpoilerContent:     true,
		SendAttachments:    true,
		DisableLinkPreview: true,
	}
	_, err = db.Insert(setting)
	return err
}

func UpdateSetting(db *xorm.Engine, setting *TelegramPushSetting) (err error) {
	has, err := db.Exist(&TelegramPushSetting{UserID: setting.UserID})
	if err != nil {
		return err
	}
	if !has {
		return errors.New("setting not found")
	}

	_, err = db.Where("user_id = ?", setting.UserID).AllCols().Update(setting)
	return err
}
