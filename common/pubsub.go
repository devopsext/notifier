package common

type PubSubAccount struct {
	Credentials string
	ProjectID   string
}

type PubOptions struct {
	Account *PubSubAccount
	Payload string
	Topic   string
}

type SubOptions struct {
	Account      *PubSubAccount
	Subscription string
	Gitlab       *GitlabOptions
}
