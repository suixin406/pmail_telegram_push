package hook

import (
	"strings"

	"github.com/ydzydzydz/pmail_telegram_push/config"
	"github.com/ydzydzydz/pmail_telegram_push/logger"
	"github.com/ydzydzydz/pmail_telegram_push/model"

	pconfig "github.com/Jinnrry/pmail/config"
	"github.com/Jinnrry/pmail/dto/parsemail"
	"github.com/Jinnrry/pmail/hooks/framework"
	"github.com/Jinnrry/pmail/models"
	"github.com/Jinnrry/pmail/utils/context"
	"github.com/go-telegram/bot"

	"github.com/ydzydzydz/pmail_telegram_push/db"
)

const (
	PLUGIN_NAME = "pmail_telegram_push"
)

type PmailTelegramPushHook struct {
	bot          *bot.Bot
	mainConfig   *pconfig.Config
	pluginConfig *config.PluginConfig
}

var _ framework.EmailHook = (*PmailTelegramPushHook)(nil)

func (h *PmailTelegramPushHook) GetName(ctx *context.Context) string {
	return PLUGIN_NAME
}

func (h *PmailTelegramPushHook) ReceiveSaveAfter(ctx *context.Context, email *parsemail.Email, ue []*models.UserEmail) {
	for _, u := range ue {
		if u.IsRead == 0 && u.Status == 0 && email.MessageId > 0 {
			setting, err := model.GetSetting(db.Instance, u.UserID)
			if err != nil || setting.ChatID == "" {
				continue
			}

			msg, err := h.sendNotification(email, setting)
			if err != nil {
				logger.PluginLogger.Error().Err(err).Msg("发送通知失败")
				continue
			}
			logger.PluginLogger.Info().Int("message_id", msg.ID).Msg("发送通知成功")

			if setting.SendAttachments && len(email.Attachments) > 0 {
				if _, err = h.sendAttachments(msg.ID, email, setting); err != nil {
					logger.PluginLogger.Error().Err(err).Msg("发送附件失败")
				} else {
					logger.PluginLogger.Info().Int("message_id", msg.ID).Msg("发送附件成功")
				}
			}
		}
	}
}

func (h *PmailTelegramPushHook) ReceiveParseBefore(ctx *context.Context, email *[]byte) {

}

func (h *PmailTelegramPushHook) ReceiveParseAfter(ctx *context.Context, email *parsemail.Email) {

}

func (h *PmailTelegramPushHook) SendAfter(ctx *context.Context, email *parsemail.Email, err map[string]error) {

}

func (h *PmailTelegramPushHook) SendBefore(ctx *context.Context, email *parsemail.Email) {

}

func (h *PmailTelegramPushHook) SettingsHtml(ctx *context.Context, url string, requestData string) string {
	switch {
	case strings.Contains(url, "setting"):
		return h.getSetting(ctx.UserID)
	case strings.Contains(url, "bot"):
		return h.getBotInfo()
	case strings.Contains(url, "submit"):
		return h.submitSetting(ctx.UserID, requestData)
	default:
		return SettingHtml
	}
}

func NewPmailTelegramPushHook(cfg *config.Config) *PmailTelegramPushHook {
	bot, err := NewBot(cfg)
	if err != nil {
		panic(err)
	}
	return &PmailTelegramPushHook{
		bot:          bot,
		mainConfig:   cfg.MainConfig,
		pluginConfig: cfg.PluginConfig,
	}
}
