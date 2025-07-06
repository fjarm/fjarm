package certmanagerv1

import (
	"github.com/pulumi/pulumi-tls/sdk/v4/go/tls"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	exportedSelfSignedCert = "selfSignedCertificateCert"
	selfSignedCertName     = "certManagerClusterIssuerSelfSignedCert"
)

// deployPulumiRootCertificateSelfSignedCert deploys a self-signed certificate that will later be used by a
// ClusterIssuer.
func deployPulumiRootCertificateSelfSignedCert(
	ctx *pulumi.Context,
	key *tls.PrivateKey,
	deps []pulumi.Resource,
) (*tls.SelfSignedCert, error) {
	certArgs := newPulumiRootCertificateSelfSignedCertArgs(
		key,
	)
	cert, err := tls.NewSelfSignedCert(
		ctx,
		selfSignedCertName,
		certArgs,
		pulumi.DependsOn(deps),
	)
	if err != nil {
		return nil, err
	}
	ctx.Export(exportedSelfSignedCert, cert.CertPem)
	return cert, nil
}

// newPulumiRootCertificateSelfSignedCertArgs sets up the args used to deploy a self-signed certificate.
func newPulumiRootCertificateSelfSignedCertArgs(
	key *tls.PrivateKey,
) *tls.SelfSignedCertArgs {
	return &tls.SelfSignedCertArgs{
		PrivateKeyPem: key.PrivateKeyPem,
		AllowedUses: pulumi.StringArray{
			pulumi.String("cert_signing"),
			pulumi.String("client_auth"),
			pulumi.String("digital_signature"),
			pulumi.String("server_auth"),
		},
		IsCaCertificate: pulumi.Bool(true),
		Subject: &tls.SelfSignedCertSubjectArgs{
			Organization: pulumi.String("Fjarm"),
		},
		ValidityPeriodHours: pulumi.Int(807660),
	}
}
