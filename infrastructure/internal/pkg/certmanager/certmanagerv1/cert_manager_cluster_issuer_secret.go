package certmanagerv1

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi-tls/sdk/v4/go/tls"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	certManagerCACertSecretName = "cert-manager-ca-cert"
)

// deployClusterIssuerTLSSecret deploys a TLS secret containing the key and certificate pair.
func deployClusterIssuerTLSSecret(
	ctx *pulumi.Context,
	provider *kubernetes.Provider,
	namespace *corev1.Namespace,
	key *tls.PrivateKey,
	cert *tls.SelfSignedCert,
	deps []pulumi.Resource,
) (*corev1.Secret, error) {
	secretArgs := newClusterIssuerTLSSecretArgs(
		namespace,
		key,
		cert,
	)
	secret, err := corev1.NewSecret(
		ctx,
		certManagerCACertSecretName,
		secretArgs,
		pulumi.Provider(provider),
		pulumi.DependsOn(deps),
	)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func newClusterIssuerTLSSecretArgs(
	namespace *corev1.Namespace,
	key *tls.PrivateKey,
	cert *tls.SelfSignedCert,
) *corev1.SecretArgs {
	return &corev1.SecretArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Secret"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(certManagerCACertSecretName),
			Namespace: namespace.Metadata.Name(),
		},
		Type: pulumi.String("kubernetes.io/tls"),
		StringData: pulumi.StringMap{
			"tls.key": key.PrivateKeyPem,
			"tls.crt": cert.CertPem,
		},
	}
}
