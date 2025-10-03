package hook

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/Jinnrry/pmail/dto/parsemail"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	log "github.com/sirupsen/logrus"
	"github.com/ydzydzydz/pmail_telegram_push/model"
)

const TEXT_MAX_SIZE = 4096

func (h *PmailTelegramPushHook) getWebButton() *models.InlineKeyboardMarkup {
	var url string
	if h.mainConfig.HttpsEnabled > 1 {
		url = "http://" + h.mainConfig.WebDomain
	} else {
		url = "https://" + h.mainConfig.WebDomain
	}

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text: "查收邮件",
					URL:  url,
				},
			},
		},
	}
}

// func removeHTMLTags(text string) string {
// 	re := regexp.MustCompile("<.*?>")
// 	return re.ReplaceAllString(text, " ")
// }

func (h *PmailTelegramPushHook) getText(email *parsemail.Email, setting *model.TelegramPushSetting) (text string) {
	text = "📧 有新邮件\n"
	text += fmt.Sprintf("🔖 主题：<b>%s</b>\n", email.Subject)
	text += fmt.Sprintf("📤 发件：&#60;%s&#62;\n", email.From.EmailAddress)
	if len(email.To) > 0 {
		text += "📥 收件："
		for _, to := range email.To {
			text += fmt.Sprintf("&#60;%s&#62; ", to.EmailAddress)
		}
		text += "\n"
	}
	if len(email.Cc) > 0 {
		text += "📋 抄送："
		for _, cc := range email.Cc {
			text += fmt.Sprintf("&#60;%s&#62; ", cc.EmailAddress)
		}
		text += "\n"
	}
	if len(email.Bcc) > 0 {
		text += "🕵️ 密送："
		for _, bcc := range email.Bcc {
			text += fmt.Sprintf("&#60;%s&#62; ", bcc.EmailAddress)
		}
		text += "\n"
	}
	if len(email.Attachments) > 0 {
		text += fmt.Sprintf("📎 附件：%d 个\n", len(email.Attachments))
	}

	if setting.ShowContent {
		size := TEXT_MAX_SIZE - len(text) - 100
		if size <= 0 {
			log.Warnf("text size too large: %s", text)
			return
		}

		var emailContent string
		if len(email.Text) > 0 {
			if len(email.Text) > size {
				emailContent = fmt.Sprintf("%s...", string(email.Text[:size]))
			} else {
				emailContent = string(email.Text)
			}
		} else if len(email.HTML) > 0 {
			if len(email.HTML) > size {
				emailContent = fmt.Sprintf("%s...", RemoveHTMLTag(string(email.HTML))[:size])
			} else {
				emailContent = RemoveHTMLTag(string(email.HTML))
			}
		}
		if len(emailContent) > 0 && setting.SpoilerContent {
			emailContent = fmt.Sprintf("<tg-spoiler>%s</tg-spoiler>", emailContent)
		}
		if len(emailContent) > 0 {
			text += fmt.Sprintf("%s\n", emailContent)
		}
	}

	return
}

func (h *PmailTelegramPushHook) sendNotification(email *parsemail.Email, setting *model.TelegramPushSetting) (msg *models.Message, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.pluginConfig.Timeout)*time.Second)
	defer cancel()

	parmas := &bot.SendMessageParams{
		ChatID:      setting.ChatID,
		Text:        h.getText(email, setting),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: h.getWebButton(),
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: &setting.DisableLinkPreview,
		},
	}

	return h.bot.SendMessage(ctx, parmas)
}

func (h *PmailTelegramPushHook) sendAttachments(id int, email *parsemail.Email, setting *model.TelegramPushSetting) (msg *models.Message, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.pluginConfig.Timeout)*time.Second)
	defer cancel()

	params := &bot.SendDocumentParams{
		ChatID: setting.ChatID,
		ReplyParameters: &models.ReplyParameters{
			MessageID: id,
			Quote:     fmt.Sprintf("📎 附件：%d 个", len(email.Attachments)),
		},
	}

	for i, attachment := range email.Attachments {
		params.Caption = fmt.Sprintf("📎 附件 %d", i+1)
		params.Document = &models.InputFileUpload{
			Filename: filepath.Base(attachment.Filename),
			Data:     bytes.NewReader(attachment.Content),
		}

		if msg, err = h.bot.SendDocument(ctx, params); err != nil {
			err = errors.Join(err, fmt.Errorf("send document failed, err: %w", err))
			continue
		}
	}
	return
}

// TODO: 合并多个附件为一个消息发送
// func (h *PmailTelegramPushHook) sendAttachmentsCombine(id int, email *parsemail.Email) (msg []*models.Message, err error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.pluginConfig.Timeout)*time.Second)
// 	defer cancel()
// 	params := &bot.SendMediaGroupParams{
// 		ChatID: h.pluginConfig.TelegramChatID,
// 		ReplyParameters: &models.ReplyParameters{
// 			MessageID: id,
// 			Quote:     fmt.Sprintf("📎 附件：%d 个", len(email.Attachments)),
// 		},
// 	}
// 	for i, attachment := range email.Attachments {
// 		params.Media = append(params.Media, &models.InputMediaDocument{
// 			Media:           filepath.Base(attachment.Filename),
// 			Caption:         fmt.Sprintf("📎 附件 %d", i+1),
// 			MediaAttachment: bytes.NewReader(attachment.Content),
// 		})
// 	}
// 	return h.bot.SendMediaGroup(ctx, params)
// }
