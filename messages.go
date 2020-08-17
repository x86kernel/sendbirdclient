package sendbirdclient

import (
	"net/url"

	"github.com/x86kernel/sendbirdclient/templates"
)

type baseMessage struct {
	MessageID  string `json:"message_id"`
	Type       string `json:"type"`
	User       User   `json:"user"`
	CustomType string `json:"custom_type"`
	ChannelURL string `json:"channel_url"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

type TextMessage struct {
	baseMessage

	Message string `json:"message"`
	Data    string `json:"data"`
	File    File   `json:"file"`
}

type FileMessage struct {
	baseMessage

	File     File   `json:"file"`
	FileName string `json:"file_name"`
	FileSize string `json:"file_size"`
	FileType string `json:"file_type"`
}

type File struct {
	URL  string `json:"url"`
	Data string `json:"data"`
}

type AdminMessage struct {
	baseMessage

	Message string `json:"message"`
	Data    string `json:"data"`
}

type messagesTemplateData struct {
	ChannelType string
	ChannelURL  string
	MessageID   string
}

type sendMessageRequest struct {
	MessageType string `json:"message_type"`
	UserId      string `json:"user_id"`
	Message     string `json:"message"`

	CustomType string `json:"custom_type"`
	Data       string `json:"data"`
}

func (r *sendMessageRequest) params() url.Values {
	q := make(url.Values)

	if r.MessageType != "" {
		q.Set("message_type", r.MessageType)
	}

	if r.UserId != "" {
		q.Set("user_id", string(r.UserId))
	}

	if r.Message != "" {
		q.Set("message", string(r.Message))
	}

	if r.CustomType != "" {
		q.Set("custom_type", string(r.CustomType))
	}

	return q
}

type sendMessageResponse struct {
	MessageType string `json:"message_type"`
	UserId      string `json:"user_id"`
	Message     string `json:"message"`
}

func (c *Client) SendMessage(channelType, channelURL string, r *sendMessageRequest) (sendMessageResponse, error) {
	pathString, err := templates.GetMessagesTemplate(messagesTemplateData{
		ChannelType: url.PathEscape(channelType),
		ChannelURL:  url.PathEscape(channelURL),
	}, templates.SendbirdURLMessagesMarkAsReadWithChannelTypeAndChannelURL)
	if err != nil {
		return sendMessageResponse{}, err
	}

	parsedURL := c.PrepareUrl(pathString)

	result := sendMessageResponse{}

	raw := r.params().Encode()
	if err := c.deleteAndReturnJSON(parsedURL, raw, &result); err != nil {
		return sendMessageResponse{}, err
	}

	return result, nil
}
