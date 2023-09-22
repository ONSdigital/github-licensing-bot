package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/ONSdigital/github-licensing-bot/pkg/github"
	"github.com/ONSdigital/github-licensing-bot/pkg/slack"
)

func main() {
	apiBase := ""
	if apiBase = os.Getenv("GITHUB_API_BASE_URI"); len(apiBase) == 0 {
		log.Fatal("Missing GITHUB_API_BASE_URI environment variable")
	}

	enterprise := ""
	if enterprise = os.Getenv("GITHUB_ENTERPRISE_NAME"); len(enterprise) == 0 {
		log.Fatal("Missing GITHUB_ENTERPRISE_NAME environmental variable")
	}

	token := ""
	if token = os.Getenv("GITHUB_TOKEN"); len(token) == 0 {
		log.Fatal("Missing GITHUB_TOKEN environmental variable")
	}

	monitoringProject := ""
	if monitoringProject = os.Getenv("MONITORING_PROJECT"); len(monitoringProject) == 0 {
		log.Fatal("Missing MONITORING_PROJECT environment variable")
	}

	overLicensedThreshold := ""
	if overLicensedThreshold = os.Getenv("OVER_LICENSED_THRESHOLD"); len(overLicensedThreshold) == 0 {
		log.Fatal("Missing OVER_LICENSED_THRESHOLD environment variable")
	}

	slackAlertsChannel := ""
	if slackAlertsChannel = os.Getenv("SLACK_ALERTS_CHANNEL"); len(slackAlertsChannel) == 0 {
		log.Fatal("Missing SLACK_ALERTS_CHANNEL environment variable")
	}

	slackPubSubTopic := ""
	if slackPubSubTopic = os.Getenv("SLACK_PUBSUB_TOPIC"); len(slackPubSubTopic) == 0 {
		log.Fatal("Missing SLACK_PUBSUB_TOPIC environment variable")
	}

	system := ""
	if system = os.Getenv("SYSTEM"); len(system) == 0 {
		log.Fatal("Missing SYSTEM environment variable")
	}

	underLicensedThreshold := ""
	if underLicensedThreshold = os.Getenv("UNDER_LICENSED_THRESHOLD"); len(underLicensedThreshold) == 0 {
		log.Fatal("Missing UNDER_LICENSED_THRESHOLD environment variable")
	}

	client := github.NewClient(apiBase, token)
	data, err := client.GetEnterpriseLicensing(enterprise)
	if err != nil {
		log.Fatalf("Failed to fetch licensing: %v", err)
	}

	totalAvailableLicences := data.Enterprise.BillingInfo.TotalAvailableLicenses
	totalLicences := data.Enterprise.BillingInfo.TotalLicenses
	licencesUsed := totalLicences - totalAvailableLicences
	text := fmt.Sprintf("*%s* is using *%d* out of *%d* total GitHub licences available.", system, licencesUsed, totalLicences)

	postSlackMessage(text, monitoringProject, slackAlertsChannel, slackPubSubTopic)
}

func postSlackMessage(text, monitoringProject, slackAlertsChannel, slackPubSubTopic string) {
	context := context.Background()
	client, err := pubsub.NewClient(context, monitoringProject)
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}

	topic := client.Topic(slackPubSubTopic)

	payload := slack.Payload{
		Text:      text,
		Username:  "GitHub Licensing Bot",
		Channel:   slackAlertsChannel,
		IconEmoji: ":github:",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to serialise Slack message payload into JSON: %v", err)
	}

	result := topic.Publish(context, &pubsub.Message{
		Data: []byte(jsonPayload),
	})

	_, err = result.Get(context)
	if err != nil {
		log.Fatalf("Failed to publish Pub/Sub message: %v", err)
	}
}
