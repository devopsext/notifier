package cmd

import (
	"context"
	"os"
	"sync"
	"time"

	pubsub "cloud.google.com/go/pubsub"
	"github.com/devopsext/notifier/common"
	"github.com/devopsext/utils"
	"github.com/spf13/cobra"
	gitlab "github.com/xanzy/go-gitlab"
	"google.golang.org/api/option"
)

var subEnv = utils.GetEnvironment()
var subLog = utils.GetLog()

var subOpts = common.SubOptions{

	Subscription: subEnv.Get("NOTIFIER_PUBSUB_SUBSCRIPTION", "").(string),
}

func startSubscriber(wg *sync.WaitGroup) {

	wg.Add(1)

	go func(wg *sync.WaitGroup) {

		defer wg.Done()

		subLog.Info("Start subscriber...")

		ctx := context.Background()

		var o option.ClientOption

		if _, err := os.Stat(subOpts.Account.Credentials); err == nil {
			o = option.WithCredentialsFile(subOpts.Account.Credentials)
		} else {
			o = option.WithCredentialsJSON([]byte(subOpts.Account.Credentials))
		}

		client, err := pubsub.NewClient(ctx, subOpts.Account.ProjectID, o)
		if err != nil {
			subLog.Panic(err)
		}

		if client == nil {
			subLog.Panic("Client is not found")
		}

		sub := client.Subscription(subOpts.Subscription)

		subLog.Info("Subscriber is up. Listening...")

		err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {

			t := subLog.Info("Got message from %s: %s", subOpts.Subscription, m.Data)

			if !utils.IsEmpty(subOpts.Gitlab.BaseURL) && !utils.IsEmpty(subOpts.Gitlab.ProjectID) {

				git, err := gitlab.NewClient(subOpts.Gitlab.Token, gitlab.WithBaseURL(subOpts.Gitlab.BaseURL))
				if err != nil {
					subLog.Error("Failed to create gitlab client: %s", err.Error())
				}

				ref := subOpts.Gitlab.ProjectRef
				if utils.IsEmpty(ref) {
					ref = "master"
				}

				vars := make(map[string]string)
				vars[subOpts.Gitlab.Variable] = string(m.Data[:])

				opt := &gitlab.RunPipelineTriggerOptions{Ref: &ref, Token: &subOpts.Gitlab.TriggerToken, Variables: vars}
				pipeline, _, err := git.PipelineTriggers.RunPipelineTrigger(subOpts.Gitlab.ProjectID, opt)

				if err != nil {
					subLog.Error("Failed to create pipeline: %s", err.Error())
				}

				if pipeline != nil {

					var spent = (time.Now().UnixNano() - t) / 1000000
					subLog.Info("Pipeline started with id: %d, spent: %v", pipeline.ID, spent)
				}
			}

			m.Ack()
		})

		if err != nil {
			subLog.Panic(err)
		}
	}(wg)
}

func sub(cmd *cobra.Command, args []string) {

	var wg sync.WaitGroup

	//startMetrics(&wg)
	startSubscriber(&wg)

	wg.Wait()
}

// GetSubCmd implements CLI interface to Subscription
func GetSubCmd(pubSubAccount *common.PubSubAccount, gitlabOptions *common.GitlabOptions) *cobra.Command {

	rootCmd := cobra.Command{
		Use:   "sub",
		Short: "Sub command",
		Run:   sub,
	}

	subOpts.Account = pubSubAccount
	subOpts.Gitlab = gitlabOptions

	flags := rootCmd.PersistentFlags()

	flags.StringVar(&subOpts.Subscription, "pubsub-subscription", subOpts.Subscription, "Pub/Sub subscription")

	return &rootCmd
}
