package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pratikkumar-mohite/kubesd/cmd/kubesd/cli"
	"github.com/pratikkumar-mohite/kubesd/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/klog"
)

// This variable is populated by goreleaser
var version string

var cf *genericclioptions.ConfigFlags

var rootCmd = &cobra.Command{
	Use:          "kubectl kubesd secret SECRET_NAME -o yaml",
	SilenceUsage: true, // for when RunE returns an error
	Short:        "Shows base64 decoded secret object",
	Example: "  kubectl kubesd secret my-secret -o yaml\n" +
		"  kubectl kubesd secret my-secret -o json\n",
	Args:               cobra.RangeArgs(1, 5),
	RunE:               run,
	Version:            versionString(),
	DisableFlagParsing: true,
}

// func main() {
// 	log := logger.NewLogger()
// 	var decodedData, err = cli.Decode()
// 	if err != nil {
// 		log.Error(err.Error())
// 		return
// 	}
// 	log.Info(decodedData)
// }

func versionString() string {
	if len(version) == 0 {
		return ""
	}
	return "v" + version
}

func init() {
	klog.InitFlags(nil)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	// hide all glog flags except for -v
	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		if f.Name != "v" {
			pflag.Lookup(f.Name).Hidden = true
		}
	})

	cf = genericclioptions.NewConfigFlags(true)

	cf.AddFlags(rootCmd.Flags())
	if err := flag.Set("logtostderr", "true"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to set logtostderr flag: %v\n", err)
		os.Exit(1)
	}
}

func run(command *cobra.Command, args []string) error {
	conn := cli.CreateConnection()
	//fmt.Println(args[0])
	log := logger.NewLogger()
	secretString, _ := conn.ReadSecret("secret1", "default")
	//secretString.WriteString(scanner.Text() + "\n")
	var decodedData, err = cli.Decode(secretString)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info(decodedData)
	return nil
}

func main() {
	defer klog.Flush()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
