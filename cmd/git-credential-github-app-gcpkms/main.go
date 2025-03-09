package main

import (
	"context"
	"os"

	"github.com/alecthomas/kong"

	"github.com/ericnorris/git-credential-github-app/internal"
	"github.com/ericnorris/git-credential-github-app/internal/gcpkms"
)

var args struct {
	KMSKey         string `name:"kms-key" help:"Google KMS key resource ID." required:""`
	ClientID       string `name:"client-id" help:"GitHub app client ID." required:""`
	InstallationID string `name:"installation-id" help:"GitHub app installation ID." required:""`
	Operation      string `arg:""`
}

func main() {
	k := kong.Parse(&args)

	if args.Operation != "get" {
		return
	}

	signer, err := gcpkms.NewSigner(context.Background(), args.KMSKey)

	if err != nil {
		k.Fatalf("%s\n", err)
	}

	credhelper := internal.NewAppCredentialHelper(args.ClientID, args.InstallationID, signer)
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
