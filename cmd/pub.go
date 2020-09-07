package cmd

import (
	"context"
	"os"
	"time"

	"github.com/devopsext/notifier/common"
	"github.com/devopsext/utils"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"

	pubsub "cloud.google.com/go/pubsub"
)

var pubEnv = utils.GetEnvironment()
var pubLog = utils.GetLog()

var pubOpts = common.PubOptions{

	Payload: pubEnv.Get("NOTIFIER_PUBSUB_PAYLOAD", "").(string),
	Topic:   pubEnv.Get("NOTIFIER_PUBSUB_TOPIC", "").(string),
}

func pub(cmd *cobra.Command, args []string) {

	t := pubLog.Info("Publishing message: %s into %s", pubOpts.Payload, pubOpts.Topic)

	ctx := context.Background()

	var o option.ClientOption

	if _, err := os.Stat(pubOpts.Account.Credentials); err == nil {
		o = option.WithCredentialsFile(pubOpts.Account.Credentials)
	} else {
		o = option.WithCredentialsJSON([]byte(pubOpts.Account.Credentials))
	}

	client, err := pubsub.NewClient(ctx, pubOpts.Account.ProjectID, o)
	if err != nil {
		pubLog.Panic(err)
	}

	topic := client.Topic(pubOpts.Topic)
	exists, err1 := topic.Exists(ctx)
	if err1 != nil {
		pubLog.Panic(err1)
	}

	if !exists {
		pubLog.Panic("Topic is not found")
	}

	serverID, err2 := topic.Publish(ctx, &pubsub.Message{Data: []byte(pubOpts.Payload)}).Get(ctx)

	if err2 != nil {
		pubLog.Error(err2)
	}

	var spent = (time.Now().UnixNano() - t) / 1000000

	pubLog.Info("Message published successfully with server id: %s, spent: %v", serverID, spent)
}

// GetPubCmd implements CLI interface to Publisher
func GetPubCmd(pubSubAccount *common.PubSubAccount) *cobra.Command {

	rootCmd := cobra.Command{
		Use:   "pub",
		Short: "Pub command",
		Run:   pub,
	}

	pubOpts.Account = pubSubAccount

	flags := rootCmd.PersistentFlags()

	flags.StringVar(&pubOpts.Payload, "pubsub-payload", pubOpts.Payload, "Pub/Sub payload")
	flags.StringVar(&pubOpts.Topic, "pubsub-topic", pubOpts.Topic, "Pub/Sub topic")

	return &rootCmd
}
