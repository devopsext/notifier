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

### Export environment variable

```sh
export NOTIFIER_PUBSUB_CREDENTIALS=account.json
export NOTIFIER_PUBSUB_PROJECT_ID=some-pubsub-project-id
```


### Run in Sub mode

```sh
./notifier sub --pubsub-subscription some-pubsub-subscription \
               --gitlab-base-url some-gitlab-base-url --gitlab-variable some-gitlab-variable \
               --gitlab-project-id some-gitlab-project-id --gitlab-project-ref some-gitlab-project-ref \
               --gitlab-trigger-token some-gitlab-trigger-toke \
               --log-format stdout --log-level debug --log-template '{{.msg}}' 
```

### Run in Pub mode

```sh
./notifier pub --pubsub-topic some-pubsub-topic --pubsub-payload some-pubsub-payload \
               --log-format stdout --log-level debug --log-template '{{.msg}}' 
```

## Usage

```
Notifier command

Usage:
  notifier [flags]
  notifier [command]

Available Commands:
  help        Help about any command
  pub         Pub command
  sub         Sub command
  version     Print the version number

Flags:
      --gitlab-base-url string        Gitlab base URL
      --gitlab-project-id string      Gitlab project ID
      --gitlab-project-ref string     Gitlab project ref
      --gitlab-token string           Gitlab token
      --gitlab-trigger-token string   Gitlab trigger token
      --gitlab-variable string        Gitlab variable
  -h, --help                          help for notifier
      --log-format string             Log format: json, text, stdout (default "text")
      --log-level string              Log level: info, warn, error, debug, panic (default "info")
      --log-template string           Log template (default "{{.func}} [{{.line}}]: {{.msg}}")
      --prometheus-listen string      Prometheus listen (default "127.0.0.1:8080")
      --prometheus-url string         Prometheus endpoint url (default "/metrics")
      --pubsub-credentials string     Pub/Sub credentials
      --pubsub-project-id string      Pub/Sub project ID
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
