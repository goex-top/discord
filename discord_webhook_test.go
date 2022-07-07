package discord

import (
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var (
	proxy      = ""
	webhookUrl = "https://discord.com/api/webhooks/xxxxx/xxxxxxx"
)

func TestDiscord(t *testing.T) {
	client := http.DefaultClient
	if proxy != "" {
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					return url.Parse(proxy)
				},
				Dial: (&net.Dialer{
					Timeout: 10 * time.Second,
				}).Dial,
			},
			Timeout: 10 * time.Second,
		}
	}
	newWebhook := &Webhook{
		Content:   "Hello",
		Username:  "I'm the King",
		AvatarUrl: "https://golang.org/lib/godoc/images/footer-gopher.jpg",
		Embeds: []Embed{
			{
				Title:       "市场异动机器人配置",
				Description: "This is the embed's description",
				Url:         "https://github.com/etaaa/go-webhooks",
				Timestamp:   GetTimestamp(),      // RETURNS NEW TIMESTAMP ACCORDING TO DISCORD'S FORMAT
				Color:       GetColor("#00ff00"), // RETURNS COLOR ACCORDING TO DISCORD'S FORMAT
				//Footer: EmbedFooter{
				//	Text: "Sent via github.com/etaaa/go-webhooks",
				//},
				Footer: EmbedFooter{
					Text:         "powered by binance",
					IconUrl:      "https://raw.githubusercontent.com/coinrust/crex/master/images/binance.jpg",
					ProxyIconUrl: "https://raw.githubusercontent.com/coinrust/crex/master/images/binance.jpg",
				},
				Thumbnail: EmbedThumbnail{
					Url: "https://raw.githubusercontent.com/coinrust/crex/master/images/binance.jpg",
				},
				Fields: []EmbedFields{
					{
						Name:   "I'm richer",
						Value:  "Win $100,000,000",
						Inline: true,
					},
					{
						Name:   "I'm richer",
						Value:  "Win $100,000,000 again",
						Inline: true,
					},
				},
			},
		},
	}

	err := SendWebhook(webhookUrl,
		client, newWebhook, false)
	if err != nil {
		t.Fatal(err)
	}
}
