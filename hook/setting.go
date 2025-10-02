package hook

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	_ "embed"

	log "github.com/sirupsen/logrus"
	"github.com/ydzydzydz/pmail_telegram_push/db"
	"github.com/ydzydzydz/pmail_telegram_push/model"
)

var (
	//go:embed dist/index.html
	SettingHtml string
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (r *Response) Json() string {
	json, err := json.Marshal(r)
	if err != nil {
		log.Errorf("marshal response failed, err: %v", err)
		return ""
	}
	return string(json)
}

type TelegramPushBotInfo struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	BotLink   string `json:"bot_link"`
}

func (h *PmailTelegramPushHook) getSetting(id int) string {
	result, err := model.GetSetting(db.Instance, id)
	if err != nil {
		response := Response{
			Code:    -1,
			Message: "get setting failed",
		}
		return response.Json()
	}
	response := Response{
		Code:    0,
		Message: "success",
		Data:    result,
	}
	return response.Json()
}

func (h *PmailTelegramPushHook) getBotInfo() string {
	me, err := h.bot.GetMe(context.Background())
	if err != nil {
		response := Response{
			Code:    -1,
			Message: "get bot me failed",
		}
		return response.Json()
	}
	response := Response{
		Code:    0,
		Message: "success",
		Data: TelegramPushBotInfo{
			Username:  me.Username,
			FirstName: me.FirstName,
			BotLink:   fmt.Sprintf("https://t.me/%s", me.Username),
		},
	}
	fmt.Println(response.Json())
	return response.Json()
}

func (h *PmailTelegramPushHook) submitSetting(id int, requestData string) string {
	var setting model.TelegramPushSetting
	if err := json.Unmarshal([]byte(requestData), &setting); err != nil {
		log.Errorf("unmarshal setting request failed, err: %v", err)
		response := Response{
			Code:    -1,
			Message: "unmarshal setting request failed",
		}
		return response.Json()
	}
	setting.UserID = id
	setting.ChatID = strings.TrimSpace(setting.ChatID)
	if err := model.UpdateSetting(db.Instance, &setting); err != nil {
		log.Errorf("update setting failed, err: %v", err)
		response := Response{
			Code:    -1,
			Message: "update setting failed",
		}
		return response.Json()
	}
	response := Response{
		Code:    0,
		Message: "success",
	}
	return response.Json()
}