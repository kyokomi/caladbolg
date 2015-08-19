package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kyokomi/caladbolg/slack"
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

	c := slack.New(*slackToken, *slackChannel, *userName, *iconURL)
	s := NewCoverageService(c)
	s.DryRun = *dryRun
	if err := s.Send(os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "error : %s\n", err.Error())
		os.Exit(2)
	}
}
