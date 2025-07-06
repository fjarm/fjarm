package certmanagerv1

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	chartNamespace = "cert-manager"
)

// deployCertManagerNamespace creates a Namespace for cert-manager to be installed in.
func deployCertManagerNamespace(
	ctx *pulumi.Context,
	provider *kubernetes.Provider,
	deps []pulumi.Resource,
) (*corev1.Namespace, error) {
	args := newCertManagerNamespaceArgs()
	ns, err := corev1.NewNamespace(
		ctx,
		chartNamespace,
		args,
		pulumi.Provider(provider),
		pulumi.DependsOn(deps),
	)
	if err != nil {
		return nil, err
	}
	return ns, nil
}

// newCertManagerNamespaceArgs creates corev1.NamespaceArgs required for the creation of the cert-manager namespace.
func newCertManagerNamespaceArgs() *corev1.NamespaceArgs {
	return &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(chartNamespace),
			Labels: pulumi.StringMap{
				"app": pulumi.String(appLabel),
			},
		},
	}
}
