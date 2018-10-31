package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
	"github.com/xenolf/lego/emca/certificate"
	"github.com/xenolf/lego/emca/challenge"
)

// Configuration type from CLI and config files.
type Configuration struct {
	context *cli.Context
}

// NewConfiguration creates a new configuration from CLI data.
func NewConfiguration(c *cli.Context) *Configuration {
	return &Configuration{context: c}
}

// KeyType the type from which private keys should be generated
func (c *Configuration) KeyType() (certificate.KeyType, error) {
	switch strings.ToUpper(c.context.GlobalString("key-type")) {
	case "RSA2048":
		return certificate.RSA2048, nil
	case "RSA4096":
		return certificate.RSA4096, nil
	case "RSA8192":
		return certificate.RSA8192, nil
	case "EC256":
		return certificate.EC256, nil
	case "EC384":
		return certificate.EC384, nil
	}

	return "", fmt.Errorf("unsupported KeyType: %s", c.context.GlobalString("key-type"))
}

// ExcludedSolvers is a list of solvers that are to be excluded.
func (c *Configuration) ExcludedSolvers() (cc []challenge.Type) {
	for _, s := range c.context.GlobalStringSlice("exclude") {
		cc = append(cc, challenge.Type(s))
	}
	return
}

// ServerPath returns the OS dependent path to the data for a specific CA
func (c *Configuration) ServerPath() string {
	srv, _ := url.Parse(c.context.GlobalString("server"))
	return strings.NewReplacer(":", "_", "/", string(os.PathSeparator)).Replace(srv.Host)
}

// CertPath gets the path for certificates.
func (c *Configuration) CertPath() string {
	return filepath.Join(c.context.GlobalString("path"), "certificates")
}

// AccountsPath returns the OS dependent path to the
// local accounts for a specific CA
func (c *Configuration) AccountsPath() string {
	return filepath.Join(c.context.GlobalString("path"), "accounts", c.ServerPath())
}

// AccountPath returns the OS dependent path to a particular account
func (c *Configuration) AccountPath(acc string) string {
	return filepath.Join(c.AccountsPath(), acc)
}

// AccountKeysPath returns the OS dependent path to the keys of a particular account
func (c *Configuration) AccountKeysPath(acc string) string {
	return filepath.Join(c.AccountPath(acc), "keys")
}