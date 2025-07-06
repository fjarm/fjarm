package certmanagerv1

import (
	"github.com/fjarm/fjarm/infrastructure/internal/pkg/common"
	"github.com/pulumi/pulumi-tls/sdk/v4/go/tls"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	exportedSelfSignedPrivateKey = "selfSignedCertificatePrivateKey"
	privateKeyName               = "certManagerClusterIssuerPrivateKey"
)

// deployRootCertificatePrivateKey creates a TLS private key to be used by a root certificate.
//
// If deploying locally to KinD, it'll be a self-signed certificate. Otherwise, it'll be a certificate from the
// Infisical PKI.
func deployRootCertificatePrivateKey(
	ctx *pulumi.Context,
	kind bool,
) (*tls.PrivateKey, error) {
	if kind {
		keyArgs := newPulumiRootCertificatePrivateKeyArgs()
		key, err := tls.NewPrivateKey(
			ctx,
			privateKeyName,
			keyArgs,
		)
		if err != nil {
			return nil, err
		}

		ctx.Export(exportedSelfSignedPrivateKey, key.PrivateKeyPem)
		return key, nil
	}

	return nil, common.ErrUnimplemented
}

func newPulumiRootCertificatePrivateKeyArgs() *tls.PrivateKeyArgs {
	return &tls.PrivateKeyArgs{
		Algorithm: pulumi.String("RSA"),
	}
}
