package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nlopes/slack"
)

func main() {
	slackToken := flag.String("token", "", "slack token")
	slackChannel := flag.String("channel", "", "slack post message channelID or channelName")
	userName := flag.String("name", "", "slack bot name")
	iconURL := flag.String("icon", "", "slack bot icon url")
	dryRun := flag.Bool("dry-run", false, "dry-run")
	flag.Parse()

	if *slackToken == "" {
		*slackToken = os.Getenv("SLACK_TOKEN")
	}
	if *slackChannel == "" {
		*slackChannel = os.Getenv("SLACK_CHANNEL")
	}

	s := Slack{
		Slack:       slack.New(*slackToken),
		channelName: *slackChannel,
		userName:    *userName,
		iconURL:     *iconURL,
		dryRun:      *dryRun,
	}

	NewCoverageService(s).Send(os.Stdin)
}

type Slack struct {
	*slack.Slack
	channelName string
	userName    string
	iconURL     string
	dryRun      bool
}

func (s Slack) NewDefaultPostMessageParams() slack.PostMessageParameters {
	params := slack.NewPostMessageParameters()
	if s.userName != "" {
		params.Username = s.userName
	}
	if s.iconURL != "" {
		params.IconURL = s.iconURL
	}
	return params
}

func (s Slack) PostDefaultMessage(message string) error {
	return s.PostMessage(message, s.NewDefaultPostMessageParams())
}

func (s Slack) PostMessage(message string, params slack.PostMessageParameters) error {
	if s.dryRun {
		fmt.Println(message)
		return nil
	}
	_, _, err := s.Slack.PostMessage(s.channelName, message, params)
	return err
}
