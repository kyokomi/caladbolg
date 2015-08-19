package slack

import "github.com/nlopes/slack"

type Client interface {
	NewAttachment() slack.Attachment
	NewAttachmentField() slack.AttachmentField
	NewDefaultPostMessageParams() slack.PostMessageParameters
	PostDefaultMessage(message string) error
	PostMessage(message string, params slack.PostMessageParameters) error
}

type client struct {
	*slack.Client
	channelName string
	userName    string
	iconURL     string
}

func New(token string, channelName string, userName string, iconURL string) Client {
	return client{
		Client:      slack.New(token),
		channelName: channelName,
		userName:    userName,
		iconURL:     iconURL,
	}
}

func (s client) NewAttachment() slack.Attachment {
	return slack.Attachment{}
}

func (s client) NewAttachmentField() slack.AttachmentField {
	return slack.AttachmentField{}
}

func (s client) NewDefaultPostMessageParams() slack.PostMessageParameters {
	params := slack.NewPostMessageParameters()
	if s.userName != "" {
		params.Username = s.userName
	}
	if s.iconURL != "" {
		params.IconURL = s.iconURL
	}
	return params
}

func (s client) PostDefaultMessage(message string) error {
	return s.PostMessage(message, s.NewDefaultPostMessageParams())
}

func (s client) PostMessage(message string, params slack.PostMessageParameters) error {
	_, _, err := s.Client.PostMessage(s.channelName, message, params)
	return err
}

var _ Client = (*client)(nil)
