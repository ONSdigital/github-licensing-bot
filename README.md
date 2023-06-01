# GitHub Licensing Bot
This repository contains a [Go](https://golang.org/) application that consumes the [GitHub GraphQL API](https://docs.github.com/en/graphql) and posts a Slack message containing GitHub licensing information about the number of remaining licences and the the total number of licences for the GitHub enterprise.

## Building
Use `make` to compile binaries for macOS and Linux.

## Running
### Environment Variables
The environment variables below are required:

```
GITHUB_ENTERPRISE_NAME  # Name of the GitHub Enterprise
GITHUB_TOKEN            # GitHub personal access token
MONITORING_PROJECT      # Google project containing the Cloud Pub/Sub topic to post alerts to
OVER_LICENSED_THRESHOLD # If the number of available licences is equal to or greater than this threshold then an additional message is displayed
SLACK_ALERTS_CHANNEL    # Name of the Slack channel to post alerts to
SLACK_PUBSUB_TOPIC      # Name of the Cloud Pub/Sub topic to post alerts to
UNDER_LICENSED_THRESOLD # If the number of available licences is equal to or less than this threshold then an additional message is displayed
```

### Token Scopes
The GitHub personal access token for using this application requires the following scope:

- `manage_billing:enterprise`

## Copyright
Copyright (C) 2021 Crown Copyright (Office for National Statistics)