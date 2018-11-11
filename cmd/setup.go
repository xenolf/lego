package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"
	"github.com/xenolf/lego/acme"
	"github.com/xenolf/lego/certcrypto"
	"github.com/xenolf/lego/log"
	"github.com/xenolf/lego/registration"
)

const filePerm os.FileMode = 0600

func setup(ctx *cli.Context, accountsStorage *AccountsStorage) (*Account, *acme.Client) {
	privateKey := accountsStorage.GetPrivateKey()

	var account *Account
	if accountsStorage.ExistsAccountFilePath() {
		account = accountsStorage.LoadAccount(privateKey)
	} else {
		account = &Account{Email: accountsStorage.GetUserID(), key: privateKey}
	}

	client := newClient(ctx, account)

	return account, client
}

func newClient(ctx *cli.Context, acc registration.User) *acme.Client {
	keyType := getKeyType(ctx)

	config := acme.NewDefaultConfig(acc).
		WithKeyType(keyType).
		WithCADirURL(ctx.GlobalString("server")).
		WithUserAgent(fmt.Sprintf("lego-cli/%s", ctx.App.Version))

	if ctx.GlobalIsSet("http-timeout") {
		config.HTTPClient.Timeout = time.Duration(ctx.GlobalInt("http-timeout")) * time.Second
	}

	client, err := acme.NewClient(config)
	if err != nil {
		log.Fatalf("Could not create client: %v", err)
	}

	setupChallenges(ctx, client)

	if client.GetExternalAccountRequired() && !ctx.GlobalIsSet("eab") {
		log.Fatal("Server requires External Account Binding. Use --eab with --kid and --hmac.")
	}

	return client
}

// getKeyType the type from which private keys should be generated
func getKeyType(ctx *cli.Context) certcrypto.KeyType {
	keyType := ctx.GlobalString("key-type")
	switch strings.ToUpper(keyType) {
	case "RSA2048":
		return certcrypto.RSA2048
	case "RSA4096":
		return certcrypto.RSA4096
	case "RSA8192":
		return certcrypto.RSA8192
	case "EC256":
		return certcrypto.EC256
	case "EC384":
		return certcrypto.EC384
	}

	log.Fatalf("Unsupported KeyType: %s", keyType)
	return ""
}

func getEmail(ctx *cli.Context) string {
	email := ctx.GlobalString("email")
	if len(email) == 0 {
		log.Fatal("You have to pass an account (email address) to the program using --email or -m")
	}
	return email
}

func createNonExistingFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0700)
	} else if err != nil {
		return err
	}
	return nil
}
