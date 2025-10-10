package hook

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	_ "embed"

	"github.com/ydzydzydz/pmail_telegram_push/db"
	"github.com/ydzydzydz/pmail_telegram_push/logger"
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
		logger.PluginLogger.Error().Err(err).Msg("marshal response failed")
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
	logger.PluginLogger.Info().Int("user_id", id).Msg("获取Telegram Push设置")
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
	logger.PluginLogger.Info().Msg("获取Telegram Bot信息")
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
	return response.Json()
}

func (h *PmailTelegramPushHook) submitSetting(id int, requestData string) string {
	logger.PluginLogger.Info().Int("user_id", id).Msg("更新Telegram Push设置")
	var setting model.TelegramPushSetting
	if err := json.Unmarshal([]byte(requestData), &setting); err != nil {
		logger.PluginLogger.Error().Err(err).Msg("unmarshal setting request failed")
		response := Response{
			Code:    -1,
			Message: "unmarshal setting request failed",
		}
		return response.Json()
	}
	setting.UserID = id
	setting.ChatID = strings.TrimSpace(setting.ChatID)
	if err := model.UpdateSetting(db.Instance, &setting); err != nil {
		logger.PluginLogger.Error().Err(err).Msg("update setting failed")
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
