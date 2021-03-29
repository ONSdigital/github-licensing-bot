package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ONSdigital/github-licensing-bot/pkg/github"
	"github.com/ONSdigital/github-licensing-bot/pkg/slack"
)

func main() {
	enterprise := ""
	if enterprise = os.Getenv("GITHUB_ENTERPRISE_NAME"); len(enterprise) == 0 {
		log.Fatal("Missing GITHUB_ENTERPRISE_NAME environmental variable")
	}

	token := ""
	if token = os.Getenv("GITHUB_TOKEN"); len(token) == 0 {
		log.Fatal("Missing GITHUB_TOKEN environmental variable")
	}

	overLicensedThreshold := ""
	if overLicensedThreshold = os.Getenv("OVER_LICENSED_THRESHOLD"); len(overLicensedThreshold) == 0 {
		log.Fatal("Missing OVER_LICENSED_THRESHOLD environment variable")
	}

	slackAlertsChannel := ""
	if slackAlertsChannel = os.Getenv("SLACK_ALERTS_CHANNEL"); len(slackAlertsChannel) == 0 {
		log.Fatal("Missing SLACK_ALERTS_CHANNEL environment variable")
	}

	slackWebHookURL := ""
	if slackWebHookURL = os.Getenv("SLACK_WEBHOOK"); len(slackWebHookURL) == 0 {
		log.Fatal("Missing SLACK_WEBHOOK environment variable")
	}

	underLicensedThreshold := ""
	if underLicensedThreshold = os.Getenv("UNDER_LICENSED_THRESHOLD"); len(underLicensedThreshold) == 0 {
		log.Fatal("Missing UNDER_LICENSED_THRESHOLD environment variable")
	}

	client := github.NewClient(token)
	data, err := client.GetEnterpriseLicensing(enterprise)
	if err != nil {
		log.Fatalf("Failed to fetch licensing: %v", err)
	}

	totalAvailableLicences := data.Enterprise.BillingInfo.TotalAvailableLicenses
	totalLicences := data.Enterprise.BillingInfo.TotalLicenses
	text := fmt.Sprintf("There are *%d* out of *%d* total GitHub licences available.", totalAvailableLicences, totalLicences)

	underLicensedThresholdCount, _ := strconv.Atoi(underLicensedThreshold)
	overLicensedThresholdCount, _ := strconv.Atoi(overLicensedThreshold)

	if totalAvailableLicences <= underLicensedThresholdCount {
		text = fmt.Sprintf("%s\nTime to order some more?", text)
	} else if totalAvailableLicences >= overLicensedThresholdCount {
		text = fmt.Sprintf("%s\nWe may be over-licensed.", text)
	}

	postSlackMessage(text, slackAlertsChannel, slackWebHookURL)
}

func postSlackMessage(text, channel, webHookURL string) {
	payload := slack.Payload{
		Text:      text,
		Username:  "GitHub Licensing Bot",
		Channel:   channel,
		IconEmoji: ":github:",
	}

	err := slack.Send(webHookURL, payload)
	if err != nil {
		log.Fatalf("Failed to send Slack message: %v", err)
	}
}
