package certmanagerv1

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	InternalClusterIssuerName = "internal-cluster-issuer"
)

// deployCertManagerInternalClusterIssuer deploys an internal ClusterIssuer resource.
func deployCertManagerInternalClusterIssuer(
	ctx *pulumi.Context,
	provider *kubernetes.Provider,
	namespace *corev1.Namespace,
	secret *corev1.Secret,
	deps []pulumi.Resource,
) (*apiextensions.CustomResource, error) {
	clusterIssuerArgs := newCertManagerInternalClusterIssuerArgs(
		namespace,
		secret,
	)
	clusterIssuer, err := apiextensions.NewCustomResource(
		ctx,
		InternalClusterIssuerName,
		clusterIssuerArgs,
		pulumi.Provider(provider),
		pulumi.DependsOn(deps),
	)
	if err != nil {
		return nil, err
	}
	return clusterIssuer, nil
}

// newCertManagerInternalClusterIssuerArgs returns a pointer to apiextensions.CustomResourceArgs that sets up a
// ClusterIssuer with a reference to self-signed certificates created by Pulumi.
func newCertManagerInternalClusterIssuerArgs(
	namespace *corev1.Namespace,
	secret *corev1.Secret,
) *apiextensions.CustomResourceArgs {
	cra := apiextensions.CustomResourceArgs{
		ApiVersion: pulumi.String("cert-manager.io/v1"),
		Kind:       pulumi.String("ClusterIssuer"),
		Metadata: metav1.ObjectMetaArgs{
			Name:      pulumi.String(InternalClusterIssuerName),
			Namespace: namespace.Metadata.Name(),
			Annotations: pulumi.StringMap{
				"config.kubernetes.io/depends-on": pulumi.String(
					"/namespaces/cert-manager/Deployment/cert-manager-webhook",
				),
			},
		},
		OtherFields: kubernetes.UntypedArgs{
			"spec": kubernetes.UntypedArgs{
				"ca": kubernetes.UntypedArgs{
					"secretName": secret.Metadata.Name(),
				},
			},
		},
	}
	return &cra
}
