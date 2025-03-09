package main

import (
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/alecthomas/kong"
	"github.com/ericnorris/git-credential-github-app/internal"
)

var args struct {
	PrivateKey     string `name:"private-key" help:"path to the private key." type:"path" required:""`
	ClientID       string `name:"client-id" help:"GitHub app client ID." required:""`
	InstallationID string `name:"installation-id" help:"GitHub app installation ID." required:""`
	Operation      string `arg:""`
}

func main() {
	k := kong.Parse(&args)

	if args.Operation != "get" {
		return
	}

	key, err := os.ReadFile(args.PrivateKey)

	if err != nil {
		k.Fatalf("could not read private key: %s\n", err)
	}

	block, _ := pem.Decode(key)

	if block == nil {
		k.Fatalf("failed to decode PEM block containing private key\n")
	}

	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		k.Fatalf("could not parse private key: %s\n", err)
	}

	credhelper := internal.NewAppCredentialHelper(args.ClientID, args.InstallationID, private)
	attrs, err := internal.ReadCredentialAttributes(os.Stdin)

	if err != nil {
		k.Fatalf("%s\n", err)
	}

	outAttrs, err := credhelper.Get(attrs)

	if err != nil {
		k.Fatalf("%s\n", err)
	}

	internal.WriteCredentialAttributes(os.Stdout, outAttrs)
}
