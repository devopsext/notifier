# Notifier

Utility that allows to publish events into Google Pub/Sub from one side and receive events from the other side by triggering Gitlab CI/CD pipeline with a variable. A content of an event will be passed to a variable.

[![GoDoc](https://godoc.org/github.com/devopsext/notifier?status.svg)](https://godoc.org/github.com/devopsext/notifier)
[![build status](https://img.shields.io/travis/devopsext/notifier/master.svg?style=flat-square)](https://travis-ci.org/devopsext/notifier)

## Features

- Publish messages with any content
- Subscribe and listen messages
- Trigger pipeline with a content of message

## Build

```sh
git clone https://github.com/devopsext/notifier.git
cd notifier/
go build
```

## Example

To try this example it's neccesary to prepare Google Pub/Sub to allow Notifier publish and receive messsages. To do so, create Service Account (under Google Gloud Console) with proper rights (read and write), download it in Json format to account.json as well as to make a topic and subscription in Pub/Sub. Gitlab repository should also have trigger token and branch.   

### Run in Sub mode

```sh
./notifier sub --pubsub-credentials account.json --pubsub-project-id some-pubsub-project-id \
               --pubsub-subscription some-pubsub-subscription \
               --gitlab-base-url some-gitlab-base-url \
               --gitlab-project-id some-gitlab-project-id --gitlab-project-ref some-gitlab-project-ref --gitlab-variable some-gitlab-variable \
               --gitlab-trigger-token some-gitlab-trigger-toke \
               --log-format stdout --log-level debug --log-template '{{.msg}}' 
```

### Run in Pub mode

```sh
./notifier pub --pubsub-credentials account.json --pubsub-project-id some-pubsub-project-id \
               --pubsub-topic some-pubsub-topic --pubsub-payload some-pubsub-payload \
               --log-format stdout --log-level debug --log-template '{{.msg}}' 
```

## Environment variables

For containerization purpose all command switches have environment variables analogs.

- NOTIFIER_LOG_FORMAT
- NOTIFIER_LOG_LEVEL
- NOTIFIER_LOG_TEMPLATE
- NOTIFIER_PROMETHEUS_URL
- NOTIFIER_PROMETHEUS_LISTEN
- NOTIFIER_PUBSUB_CREDENTIALS
- NOTIFIER_PUBSUB_PROJECT_ID
- NOTIFIER_PUBSUB_PAYLOAD
- NOTIFIER_PUBSUB_TOPIC
- NOTIFIER_PUBSUB_SUBSCRIPTION
- NOTIFIER_GITLAB_TOKEN
- NOTIFIER_GITLAB_BASE_URL
- NOTIFIER_GITLAB_PROJECT_ID
- NOTIFIER_GITLAB_PROJECT_REF
- NOTIFIER_GITLAB_VARIABLE
- NOTIFIER_GITLAB_TRIGGER_TOKEN
