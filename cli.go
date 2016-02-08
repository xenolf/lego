// Let's Encrypt client to go!
// CLI application for generating Let's Encrypt certificates using the ACME package.
package main

import (
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/xenolf/lego/acme"
	"github.com/xenolf/lego/dns_provider"
)

// Logger is used to log errors; if nil, the default log.Logger is used.
var Logger *log.Logger

// logger is an helper function to retrieve the available logger
func logger() *log.Logger {
	if Logger == nil {
		Logger = log.New(os.Stderr, "", log.LstdFlags)
	}
	return Logger
}

var gittag string

// Return a string describing all the DNS providers together, like
// "\n\t{name}: {config description}", suitable for use in a CLI flag description
func dnsProviderHelpString() string {
	lines := []string{}

	for _, entry := range dns_provider.Registry.Entries() {
		lines = append(lines, "\n\t"+entry.Name+": "+entry.ConfigDescription)
	}

	sort.Strings(lines)
	return strings.Join(lines, "")
}

func main() {
	app := cli.NewApp()
	app.Name = "lego"
	app.Usage = "Let's encrypt client to go!"

	version := "0.2.0"
	if strings.HasPrefix(gittag, "v") {
		version = gittag
	}

	app.Version = version

	acme.UserAgent = "lego/" + app.Version

	cwd, err := os.Getwd()
	if err != nil {
		logger().Fatal("Could not determine current working directory. Please pass --path.")
	}
	defaultPath := path.Join(cwd, ".lego")

	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "Register an account, then create and install a certificate",
			Action: run,
		},
		{
			Name:   "revoke",
			Usage:  "Revoke a certificate",
			Action: revoke,
		},
		{
			Name:   "renew",
			Usage:  "Renew a certificate",
			Action: renew,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "days",
					Value: 0,
					Usage: "The number of days left on a certificate to renew it.",
				},
				cli.BoolFlag{
					Name:  "reuse-key",
					Usage: "Used to indicate you want to reuse your current private key for the new certificate.",
				},
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "domains, d",
			Usage: "Add domains to the process",
		},
		cli.StringFlag{
			Name:  "server, s",
			Value: "https://acme-v01.api.letsencrypt.org/directory",
			Usage: "CA hostname (and optionally :port). The server certificate must be trusted in order to avoid further modifications to the client.",
		},
		cli.StringFlag{
			Name:  "email, m",
			Usage: "Email used for registration and recovery contact.",
		},
		cli.IntFlag{
			Name:  "rsa-key-size, B",
			Value: 2048,
			Usage: "Size of the RSA key.",
		},
		cli.StringFlag{
			Name:  "path",
			Usage: "Directory to use for storing the data",
			Value: defaultPath,
		},
		cli.StringSliceFlag{
			Name:  "exclude, x",
			Usage: "Explicitly disallow solvers by name from being used. Solvers: \"http-01\", \"tls-sni-01\".",
		},
		cli.StringFlag{
			Name:  "http",
			Usage: "Set the port and interface to use for HTTP based challenges to listen on. Supported: interface:port or :port",
		},
		cli.StringFlag{
			Name:  "tls",
			Usage: "Set the port and interface to use for TLS based challenges to listen on. Supported: interface:port or :port",
		},
		cli.StringFlag{
			Name: "dns",
			Usage: "Solve a DNS challenge using the specified provider. Disables all other challenges." +
				"\n\tCredentials for providers have to be passed through environment variables." +
				"\n\tFor a more detailed explanation of the parameters, please see the online docs." +
				"\n\tValid providers:" + dnsProviderHelpString(),
		},
	}

	app.Run(os.Args)
}
