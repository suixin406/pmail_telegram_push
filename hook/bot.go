package hook

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/ydzydzydz/pmail_telegram_push/config"
	"golang.org/x/net/proxy"
)

func NewBot(config *config.Config) (*bot.Bot, error) {
	opts := []bot.Option{
		bot.WithCheckInitTimeout(time.Duration(config.PluginConfig.Timeout) * time.Second),
	}
	if config.PluginConfig.Debug {
		opts = append(opts, bot.WithDebug())
	}

	if config.PluginConfig.Proxy == "" {
		return newBotWithOutProxy(config, opts...)
	}
	parsedURL, err := url.Parse(config.PluginConfig.Proxy)
	if err != nil {
		panic(err)
	}
	switch strings.ToLower(parsedURL.Scheme) {
	case "socks5":
		return newBotWithSocks5Proxy(config, parsedURL, opts...)
	case "http", "https":
		return newBotWithHTTPProxy(config, parsedURL, opts...)
	default:
		return newBotWithOutProxy(config, opts...)
	}
}

func newBotWithOutProxy(config *config.Config, options ...bot.Option) (*bot.Bot, error) {
	return bot.New(config.PluginConfig.TelegramBotToken, options...)
}

func newBotWithSocks5Proxy(config *config.Config, proxyURL *url.URL, options ...bot.Option) (*bot.Bot, error) {
	var auth *proxy.Auth
	if proxyURL.User != nil {
		password, _ := proxyURL.User.Password()
		auth = &proxy.Auth{
			User:     proxyURL.User.Username(),
			Password: password,
		}
	}
	dialer, err := proxy.SOCKS5(
		"tcp",
		proxyURL.Host,
		auth,
		proxy.Direct,
	)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		},
	}
	opts := append(options, bot.WithHTTPClient(time.Duration(config.PluginConfig.Timeout)*time.Second, httpClient))
	return bot.New(config.PluginConfig.TelegramBotToken, opts...)
}

func newBotWithHTTPProxy(config *config.Config, proxyURL *url.URL, options ...bot.Option) (*bot.Bot, error) {
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	options = append(options, bot.WithHTTPClient(time.Duration(config.PluginConfig.Timeout)*time.Second, httpClient))
	return bot.New(config.PluginConfig.TelegramBotToken, options...)
}
