package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/teamhephy/k8s-claimer/cli/commands"
)

// This value is overwritten by the linker during build.
var version = "dev"

func main() {
	app := cli.NewApp()
	app.Name = "k8s-claimer"
	app.Version = version
	app.Usage = "This CLI can be used against a k8s-claimer server to acquire and release leases"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "server",
			Value: "",
			Usage: "The k8s-claimer server to talk to",
		},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name: "lease",
			Subcommands: []cli.Command{
				cli.Command{
					Name: "create",
					Usage: `Creates a new lease and returns 'export' statements to set the lease values as environment variables. Set the 'env-prefix' flag to prefix the environment variable names. If you pass that flag, a '_' character will separate the prefix with the rest of the environment variable name. Below are the basic environment variable names:

- IP - the IP address of the Kubernetes master server
- TOKEN - contains the lease token. Use this when you run 'k8s-claimer-cli lease delete'
- CLUSTER_NAME - contains the name of the cluster. For informational purposes only

The Kubeconfig file will be written to kubeconfig-file
`,
					Action: commands.CreateLease,
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "duration",
							Value: 10,
							Usage: "The duration of the lease in seconds",
						},
						cli.StringFlag{
							Name:  "env-prefix",
							Value: "",
							Usage: "The prefix for all environment variables that this command sets",
						},
						cli.StringFlag{
							Name:  "kubeconfig-file",
							Value: "./kubeconfig.yaml",
							Usage: "The location of the resulting Kubeconfig file",
						},
						cli.StringFlag{
							Name:  "cluster-regex",
							Value: "",
							Usage: "A regular expression that will be used to match which cluster you lease",
						},
						cli.StringFlag{
							Name:  "cluster-version",
							Value: "",
							Usage: "A version string that will be used to find a cluster to lease",
						},
						cli.StringFlag{
							Name:  "provider",
							Value: "",
							Usage: "Which cloud provider to use when creating a cluster lease. Acceptable values are azure and google. If a value is not provided it will return an error.",
						},
					},
				},
				cli.Command{
					Name:   "delete",
					Action: commands.DeleteLease,
					Usage: `Releases a currently held lease. Pass the lease token as the first and only parameter to this command. For example:

k8s-claimer-cli lease delete $TOKEN
`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "provider",
							Value: "",
							Usage: "Which cloud provider to use when deleting a cluster lease. Acceptable values are azure and google. If a value is not provided it will return an error.",
						},
					},
				},
			},
		},
	}
	app.Run(os.Args)
}
