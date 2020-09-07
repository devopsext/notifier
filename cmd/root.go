package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/devopsext/notifier/common"
	"github.com/devopsext/utils"
	"github.com/spf13/cobra"
)

// VERSION of the app
var VERSION = "unknown"

var rootLog = utils.GetLog()
var rootEnv = utils.GetEnvironment()

type rootOptions struct {
	LogFormat   string
	LogLevel    string
	LogTemplate string

	PrometheusURL    string
	PrometheusListen string

	PubSubAccount common.PubSubAccount
	GitlabOptions common.GitlabOptions
}

var rootOpts = rootOptions{

	LogFormat:   rootEnv.Get("NOTIFIER_LOG_FORMAT", "text").(string),
	LogLevel:    rootEnv.Get("NOTIFIER_LOG_LEVEL", "info").(string),
	LogTemplate: rootEnv.Get("NOTIFIER_LOG_TEMPLATE", "{{.func}} [{{.line}}]: {{.msg}}").(string),

	PrometheusURL:    rootEnv.Get("NOTIFIER_PROMETHEUS_URL", "/metrics").(string),
	PrometheusListen: rootEnv.Get("NOTIFIER_PROMETHEUS_LISTEN", "127.0.0.1:8080").(string),

	PubSubAccount: common.PubSubAccount{
		Credentials: rootEnv.Get("NOTIFIER_PUBSUB_CREDENTIALS", "").(string),
		ProjectID:   rootEnv.Get("NOTIFIER_PUBSUB_PROJECT_ID", "").(string),
	},

	GitlabOptions: common.GitlabOptions{
		Token:        rootEnv.Get("NOTIFIER_GITLAB_TOKEN", "").(string),
		BaseURL:      rootEnv.Get("NOTIFIER_GITLAB_BASE_URL", "").(string),
		ProjectID:    rootEnv.Get("NOTIFIER_GITLAB_PROJECT_ID", "").(string),
		ProjectRef:   rootEnv.Get("NOTIFIER_GITLAB_PROJECT_REF", "").(string),
		Variable:     rootEnv.Get("NOTIFIER_GITLAB_VARIABLE", "").(string),
		TriggerToken: rootEnv.Get("NOTIFIER_GITLAB_TRIGGER_TOKEN", "").(string),
	},
}

func interceptSyscall() {

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	go func() {
		<-c
		rootLog.Info("Exiting...")
		os.Exit(1)
	}()
}

// Execute of the cmd
func Execute() {

	rootCmd := &cobra.Command{
		Use:   "notifier",
		Short: "Notifier",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

			rootLog.CallInfo = true
			rootLog.Init(rootOpts.LogFormat, rootOpts.LogLevel, rootOpts.LogTemplate)

		},
		Run: func(cmd *cobra.Command, args []string) {

			rootLog.Info("Booting...")
		},
	}

	flags := rootCmd.PersistentFlags()

	flags.StringVar(&rootOpts.LogFormat, "log-format", rootOpts.LogFormat, "Log format: json, text, stdout")
	flags.StringVar(&rootOpts.LogLevel, "log-level", rootOpts.LogLevel, "Log level: info, warn, error, debug, panic")
	flags.StringVar(&rootOpts.LogTemplate, "log-template", rootOpts.LogTemplate, "Log template")

	interceptSyscall()

	flags.StringVar(&rootOpts.PrometheusURL, "prometheus-url", rootOpts.PrometheusURL, "Prometheus endpoint url")
	flags.StringVar(&rootOpts.PrometheusListen, "prometheus-listen", rootOpts.PrometheusListen, "Prometheus listen")

	flags.StringVar(&rootOpts.PubSubAccount.Credentials, "pubsub-credentials", rootOpts.PubSubAccount.Credentials, "Pub/Sub credentials")
	flags.StringVar(&rootOpts.PubSubAccount.ProjectID, "pubsub-project-id", rootOpts.PubSubAccount.ProjectID, "Pub/Sub project ID")

	flags.StringVar(&rootOpts.GitlabOptions.Token, "gitlab-token", rootOpts.GitlabOptions.Token, "Gitlab token")
	flags.StringVar(&rootOpts.GitlabOptions.BaseURL, "gitlab-base-url", rootOpts.GitlabOptions.BaseURL, "Gitlab base URL")
	flags.StringVar(&rootOpts.GitlabOptions.ProjectID, "gitlab-project-id", rootOpts.GitlabOptions.ProjectID, "Gitlab project ID")
	flags.StringVar(&rootOpts.GitlabOptions.ProjectRef, "gitlab-project-ref", rootOpts.GitlabOptions.ProjectRef, "Gitlab project ref")
	flags.StringVar(&rootOpts.GitlabOptions.TriggerToken, "gitlab-trigger-token", rootOpts.GitlabOptions.TriggerToken, "Gitlab trigger token")
	flags.StringVar(&rootOpts.GitlabOptions.Variable, "gitlab-variable", rootOpts.GitlabOptions.Variable, "Gitlab variable")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(VERSION)
		},
	})

	rootCmd.AddCommand(GetPubCmd(&rootOpts.PubSubAccount))
	rootCmd.AddCommand(GetSubCmd(&rootOpts.PubSubAccount, &rootOpts.GitlabOptions))

	if err := rootCmd.Execute(); err != nil {
		rootLog.Error(err)
		os.Exit(1)
	}
}
