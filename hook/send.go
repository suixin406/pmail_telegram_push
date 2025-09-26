package hook

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/Jinnrry/pmail/dto/parsemail"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	log "github.com/sirupsen/logrus"
)

const TEXT_MAX_SIZE = 4096

func (h *PmailTelegramPushHook) getWebButton() *models.InlineKeyboardMarkup {
	var url string
	if h.httpsEnabled > 1 {
		url = "http://" + h.webDomain
	} else {
		url = "https://" + h.webDomain
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

func removeHTMLTags(text string) string {
	re := regexp.MustCompile("<.*?>")
	return re.ReplaceAllString(text, "")
}

func (h *PmailTelegramPushHook) getText(email *parsemail.Email) (text string) {
	text = "📧 有新邮件\n"
	text += fmt.Sprintf("主题：<b>%s</b>\n", email.Subject)
	text += fmt.Sprintf("发件：&#60;%s&#62;\n", email.From.EmailAddress)
	if len(email.To) > 0 {
		text += "收件："
		for _, to := range email.To {
			text += fmt.Sprintf("&#60;%s&#62; ", to.EmailAddress)
		}
		text += "\n"
	}
	if len(email.Cc) > 0 {
		text += "抄送："
		for _, cc := range email.Cc {
			text += fmt.Sprintf("&#60;%s&#62; ", cc.EmailAddress)
		}
		text += "\n"
	}
	if len(email.Bcc) > 0 {
		text += "密送："
		for _, bcc := range email.Bcc {
			text += fmt.Sprintf("&#60;%s&#62; ", bcc.EmailAddress)
		}
		text += "\n"
	}
	if len(email.Attachments) > 0 {
		text += fmt.Sprintf("附件：%d 个\n", len(email.Attachments))
	}

	if h.showContent {
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
				emailContent = fmt.Sprintf("%s...", removeHTMLTags(string(email.HTML))[:size])
			} else {
				emailContent = removeHTMLTags(string(email.HTML))
			}
		}
		if len(emailContent) > 0 && h.spoilerContent {
			emailContent = fmt.Sprintf("<tg-spoiler>%s</tg-spoiler>", emailContent)
		}
		if len(emailContent) > 0 {
			text += fmt.Sprintf("内容：\n%s\n", emailContent)
		}
	}

	return
}

func (h *PmailTelegramPushHook) sendNotification(email *parsemail.Email) (msg *models.Message, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()
	return h.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      h.chatId,
		Text:        h.getText(email),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: h.getWebButton(),
	})
}

func (h *PmailTelegramPushHook) sendAttachments(id int, email *parsemail.Email) (msg *models.Message, err error) {
	count := 0
	for _, attachment := range email.Attachments {
		count++
		ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
		defer cancel()
		msg, err = h.bot.SendDocument(ctx, &bot.SendDocumentParams{
			ChatID:  h.chatId,
			Caption: fmt.Sprintf("附件 %d", count),
			Document: &models.InputFileUpload{
				Filename: filepath.Base(attachment.Filename),
				Data:     bytes.NewReader(attachment.Content),
			},
			ReplyParameters: &models.ReplyParameters{
				MessageID: id,
			},
		})
		if err != nil {
			err = errors.Join(err, fmt.Errorf("send document failed, err: %w", err))
			continue
		}
	}
	return
}

func (h *PmailTelegramPushHook) SendMessage() (msg *models.Message, err error) {
	return h.bot.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: h.chatId,
		Text:   "测试消息",
	})
}
