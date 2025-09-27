package hook

import (
	"github.com/ydzydzydz/pmail_telegram_push/config"

	_ "embed"

	pconfig "github.com/Jinnrry/pmail/config"
	"github.com/Jinnrry/pmail/dto/parsemail"
	"github.com/Jinnrry/pmail/hooks/framework"
	"github.com/Jinnrry/pmail/models"
	"github.com/Jinnrry/pmail/utils/context"
	"github.com/go-telegram/bot"
	log "github.com/sirupsen/logrus"
)

const (
	PLUGIN_NAME = "pmail_telegram_push"
)

var (
	//go:embed setting.html
	SettingHtml string
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
		if u.UserID == 1 && u.IsRead == 0 && u.Status == 0 && email.MessageId > 0 {
			msg, err := h.sendNotification(email)
			if err != nil {
				log.Errorf("send notification failed, err: %v", err)
				continue
			}

			if h.pluginConfig.SendAttachment && len(email.Attachments) > 0 {
				if _, err = h.sendAttachments(msg.ID, email); err != nil {
					log.Errorf("send attachments failed, err: %v", err)
				} else {
					log.Infof("send attachments success, message id: %d", msg.ID)
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
	return SettingHtml
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
