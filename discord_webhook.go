package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Returns a timestamp for the footer according to Discord's format (ISO8601)
func GetTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05-0700")
}

// Transforms hex to decimal (required for webhooks)
func GetColor(hexColor string) int {
	hexColor = strings.Replace(hexColor, "#", "", -1)
	decimalColor, err := strconv.ParseInt(hexColor, 16, 64)
	if err != nil {
		return 0
	}
	return int(decimalColor)
}

type Webhook struct {
	Content   string  `json:"content,omitempty"`
	Username  string  `json:"username,omitempty"`
	AvatarUrl string  `json:"avatar_url,omitempty"`
	Tts       bool    `json:"tts,omitempty"`
	Embeds    []Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Title       string         `json:"title,omitempty"`
	Description string         `json:"description,omitempty"`
	Url         string         `json:"url,omitempty"`
	Timestamp   string         `json:"timestamp,omitempty"`
	Color       int            `json:"color,omitempty"`
	Footer      EmbedFooter    `json:"footer,omitempty"`
	Image       EmbedImage     `json:"image,omitempty"`
	Thumbnail   EmbedThumbnail `json:"thumbnail,omitempty"`
	Video       EmbedVideo     `json:"video,omitempty"`
	Provider    EmbedProvider  `json:"provider,omitempty"`
	Author      EmbedAuthor    `json:"author,omitempty"`
	Fields      []EmbedFields  `json:"fields,omitempty"`
}

type EmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconUrl      string `json:"icon_url,omitempty"`
	ProxyIconUrl string `json:"proxy_icon_url,omitempty"`
}

type EmbedImage struct {
	Url      string `json:"url,omitempty"`
	ProxyUrl string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type EmbedThumbnail struct {
	Url      string `json:"url,omitempty"`
	ProxyUrl string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type EmbedVideo struct {
	Url      string `json:"url,omitempty"`
	ProxyUrl string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type EmbedProvider struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

type EmbedAuthor struct {
	Name         string `json:"name,omitempty"`
	Url          string `json:"url,omitempty"`
	IconUrl      string `json:"icon_url,omitempty"`
	ProxyIconUrl string `json:"proxy_icon_url,omitempty"`
}

type EmbedFields struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

// Execute the webhook request
func SendWebhook(webookUrl string, client *http.Client, webhook *Webhook, retryOnRateLimit bool) error {
	if webhook == nil {
		return errors.New("webhook is nil")
	}
	if webhook.Content == "" && len(webhook.Embeds) == 0 {
		return errors.New("You must attach atleast one of these: Content; Embeds")
	}
	if len(webhook.Embeds) > 10 {
		return errors.New("Maximum number of embeds per webhook is 10")
	}
	jsonData, err := json.Marshal(webhook)
	if err != nil {
		return err
	}
	for {
		res, err := client.Post(webookUrl, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
		switch res.StatusCode {
		case 204:
			res.Body.Close()
			return nil
		case 429:
			res.Body.Close()
			if !retryOnRateLimit {
				return errors.New("Webhook ratelimited")
			}
			timeout, err := strconv.Atoi(res.Header.Get("retry-after"))
			if err != nil {
				time.Sleep(5 * time.Second)
			} else {
				time.Sleep(time.Duration(timeout) * time.Millisecond)
			}
		default:
			res.Body.Close()
			return errors.New(fmt.Sprintf("Bad request (Status %d)", res.StatusCode))
		}
	}
}
