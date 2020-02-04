package cmd

import (
	"context"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"krabbitmqctl/kctl"
	"log"
	"os"
	"path/filepath"
)

const version = "0.1.0"

type Options struct {
	version     bool
	serviceName string
	kubeConfig  string
	context     string
	namespace   string
	podName     string
}

var opts = &Options{
	serviceName: "rabbitmq",
	namespace:   "default",
	podName:     "",
}

func Run() {
	cmd := &cobra.Command{}
	cmd.Short = "rabbitmqctl kubernetes interface command "
	cmd.Use = "rabbitmqctl command"
	cmd.Flags().StringVarP(&opts.serviceName, "service", "s", opts.serviceName, "RabbitMQ Service")
	cmd.Flags().StringVar(&opts.context, "context", opts.context, "Kubernetes context to use. Default to current context configured in kubeconfig.")
	cmd.Flags().StringVar(&opts.kubeConfig, "kubeconfig", opts.kubeConfig, "Path to kubeconfig file to use")
	cmd.Flags().StringVarP(&opts.namespace, "namespace", "n", opts.namespace, "Kubernetes namespace to use. Default to namespace configured in Kubernetes context")
	cmd.Flags().StringVarP(&opts.podName, "podname", "p", opts.podName, "Pod where execute the command. Default is \"\" pick one random")
	cmd.Flags().BoolVarP(&opts.version, "version", "v", opts.version, "Print the version and exit")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if opts.version {
			fmt.Printf("%s\n", version)
			return nil
		}

		narg := len(args)
		if narg == 0 {
			return cmd.Help()
		}
		config, err := parseConfig(args)
		if err != nil {
			log.Println(err)
			os.Exit(2)
		}

		ctx, cancel := context.WithCancel(context.Background())

		defer cancel()

		resultOk, _, err2 := kctl.Run(ctx, config)
		if err2 != nil {
			fmt.Println(err2)
			os.Exit(1)
		}
		fmt.Printf("%s", resultOk)

		return nil
	}

	cmd.Execute()
}

func parseConfig(args []string) (*kctl.Config, error) {
	kubeConfig, err := getKubeConfig()
	if err != nil {
		return nil, err
	}

	return &kctl.Config{
		KubeConfig:  kubeConfig,
		ContextName: opts.context,
		NameSpace:   opts.namespace,
		CtlCommand:  args,
		ServiceName: opts.serviceName,
		PodName:     opts.podName,
	}, nil
}

func getKubeConfig() (string, error) {
	var kubeconfig string
	if kubeconfig = opts.kubeConfig; kubeconfig != "" {
		return kubeconfig, nil
	}
	if kubeconfig = os.Getenv("KUBECONFIG"); kubeconfig != "" {
		return kubeconfig, nil
	}
	home, err := homedir.Dir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get user home directory")
	}

	kubeconfig = filepath.Join(home, ".kube/config")
	return kubeconfig, nil
}
