package valkeyv1

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	helmv4 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v4"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	clusterAppLabel = "valkey"
	chartName       = "valkey"
	chartRepo       = "oci://registry-1.docker.io/bitnamicharts/valkey"
	chartVersion    = "3.0.16"
)

// DeployValkeyCluster sets up the required resources needed to run Valkey including a namespace, TLS certificate, and
// Helm chart.
func DeployValkeyCluster(
	ctx *pulumi.Context,
	provider *kubernetes.Provider,
	deps []pulumi.Resource,
) ([]pulumi.Resource, error) {
	namespace, err := deployValkeyClusterNamespace(
		ctx,
		provider,
		[]pulumi.Resource{},
	)
	if err != nil {
		return nil, err
	}

	cert, err := deployValkeyClusterCertificate(
		ctx,
		namespace,
		provider,
		append(deps, namespace),
	)
	if err != nil {
		return nil, err
	}

	// TODO(2025-07-06): Use better (Infisical) method to supply ACL credentials.
	commonConfig := &valkeyConfig{
		DefaultUserCredentials:  "somepassword",
		SentinelUserCredentials: "somepassword",
		ReplicaUserCredentials:  "somepassword",
		Users: []*valkeyUser{
			{
				Username:        "test",
				Password:        "test",
				EnabledCommands: []string{"+AUTH", "+INFO", "+ACL", "+PING", "+GET", "+SET", "~*"},
			},
		},
	}
	configContent, err := newValkeyCommonConfig(commonConfig)
	if err != nil {
		return nil, err
	}

	chart, err := deployValkeyClusterHelmChart(
		ctx,
		namespace,
		commonConfig,
		configContent,
		provider,
		append(deps, namespace, cert),
	)
	if err != nil {
		return nil, err
	}

	return []pulumi.Resource{
		namespace,
		cert,
		chart,
	}, nil
}

// deployValkeyClusterHelmChart deploys a Valkey cluster using the bitnami Helm chart.
func deployValkeyClusterHelmChart(
	ctx *pulumi.Context,
	namespace *corev1.Namespace,
	commonConfig *valkeyConfig,
	aclContent string,
	provider *kubernetes.Provider,
	deps []pulumi.Resource,
) (pulumi.Resource, error) {
	args := newValkeyClusterHelmChartArgs(
		namespace,
		commonConfig,
		aclContent,
	)

	chart, err := helmv4.NewChart(
		ctx,
		chartName,
		args,
		pulumi.Provider(provider),
		pulumi.DependsOn(deps),
	)
	if err != nil {
		return nil, err
	}
	return chart, nil
}

// newValkeyClusterHelmChartArgs constructs the Helm chart values needed to deploy Valkey to k8s.
func newValkeyClusterHelmChartArgs(
	namespace *corev1.Namespace,
	commonConfig *valkeyConfig,
	aclContent string,
) *helmv4.ChartArgs {
	chartArgs := &helmv4.ChartArgs{
		Chart:     pulumi.String(chartRepo),
		Namespace: namespace.Metadata.Name(),
		Version:   pulumi.String(chartVersion),
		Values: pulumi.Map{
			"fullnameOverride": pulumi.String("valkey"),
			"architecture":     pulumi.String("replication"),
			"auth": pulumi.Map{
				"enabled": pulumi.Bool(true),
				// The password for the default user
				// TODO(2025-07-05): Enable supplying a custom password for the default user.
				"password": pulumi.String(commonConfig.DefaultUserCredentials),
			},
			"commonConfiguration": pulumi.String(aclContent),
			"image": pulumi.Map{
				"digest": pulumi.String("sha256:0384ca2eec63789450b2e07a00f377c2c9d0b548c2e346e1003bc0dd629fa71a"),
			},
			"sentinel": pulumi.Map{
				"enabled": pulumi.Bool(true),
				"image": pulumi.Map{
					"digest": pulumi.String("sha256:071cb353bc17f27492655c710386d5e3afc5c36d8057ea5d48c5886da6f1bc3a"),
				},
				"primarySet": pulumi.String("valkey-primary"),
				"resources": pulumi.Map{
					"limits": pulumi.Map{
						"cpu":    pulumi.String("600m"),
						"memory": pulumi.String("750Mi"),
					},
					"requests": pulumi.Map{
						"cpu":    pulumi.String("500m"),
						"memory": pulumi.String("500Mi"),
					},
				},
			},
			"tls": pulumi.Map{
				"enabled":         pulumi.Bool(true),
				"existingSecret":  pulumi.String(clusterCertificateSecretName),
				"certFilename":    pulumi.String("tls.crt"),
				"certKeyFilename": pulumi.String("tls.key"),
				"certCAFilename":  pulumi.String("ca.crt"),
			},
		},
	}
	return chartArgs
}
