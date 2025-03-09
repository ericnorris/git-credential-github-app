package gcpkms

import (
	"context"
	"crypto"
	"fmt"
	"io"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "cloud.google.com/go/kms/apiv1/kmspb"
)

type Signer struct {
	// not ideal, but crypto.Signer doesn't have an obvious way to pass in a context.
	// see https://github.com/golang/go/issues/28427
	ctx context.Context

	client     *cloudkms.KeyManagementClient
	keyVersion *kmspb.CryptoKeyVersion
}

func NewSigner(ctx context.Context, resourceID string) (*Signer, error) {
	client, err := cloudkms.NewKeyManagementClient(ctx)

	if err != nil {
		return nil, fmt.Errorf("could not construct KMS client: %w", err)
	}

	keyVersion, err := client.GetCryptoKeyVersion(ctx, &kmspb.GetCryptoKeyVersionRequest{
		Name: resourceID,
	})

	if err != nil {
		return nil, fmt.Errorf("could not fetch KMS key metadata: %w", err)
	}

	if keyVersion.Algorithm != kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_2048_SHA256 &&
		keyVersion.Algorithm != kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_3072_SHA256 &&
		keyVersion.Algorithm != kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_4096_SHA256 {

		return nil, fmt.Errorf("KMS key must be of type RSA_SIGN_PKCS1_*_SHA256")
	}

	return &Signer{ctx, client, keyVersion}, nil
}

func (s *Signer) Public() crypto.PublicKey {
	return nil
}

func (s *Signer) Sign(
	rand io.Reader,
	digest []byte,
	opts crypto.SignerOpts,
) (signature []byte, err error) {
	var kmspbDigest kmspb.Digest

	switch opts.HashFunc() {
	case crypto.SHA256:
		kmspbDigest = kmspb.Digest{Digest: &kmspb.Digest_Sha256{Sha256: digest}}

	case crypto.SHA384:
		kmspbDigest = kmspb.Digest{Digest: &kmspb.Digest_Sha384{Sha384: digest}}

	case crypto.SHA512:
		kmspbDigest = kmspb.Digest{Digest: &kmspb.Digest_Sha512{Sha512: digest}}

	default:
		return nil, fmt.Errorf("cannot convert hash function %v to KMS digest", opts.HashFunc())
	}

	signRequest := &kmspb.AsymmetricSignRequest{
		Name:   s.keyVersion.Name,
		Digest: &kmspbDigest,
	}

	signResponse, err := s.client.AsymmetricSign(s.ctx, signRequest)

	if err != nil {
		return nil, fmt.Errorf("error in AsymmetricSign(): %w", err)
	}

	return signResponse.Signature, nil
}
