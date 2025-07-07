package certmanagerv1

import (
	"fmt"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	helmv4 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v4"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

const (
	appLabel     = "cert-manager"
	chartName    = "cert-manager"
	chartRepo    = "https://charts.jetstack.io"
	chartVersion = "1.17.2"
	configKind   = "kind"
)

// DeployCertManager deploys the cert-manager Helm chart and a ClusterIssuer.
func DeployCertManager(ctx *pulumi.Context, provider *kubernetes.Provider) ([]pulumi.Resource, error) {
	kind := config.GetBool(ctx, configKind)

	namespace, err := deployCertManagerNamespace(
		ctx,
		provider,
		[]pulumi.Resource{},
	)
	if err != nil {
		return nil, err
	}

	chart, err := deployCertManagerHelmChart(
		ctx,
		namespace,
		provider,
		[]pulumi.Resource{namespace},
	)
	if err != nil {
		return nil, err
	}

	privateKey, err := deployRootCertificatePrivateKey(
		ctx,
		kind,
	)
	if err != nil {
		return nil, err
	}

	cert, err := deployPulumiRootCertificateSelfSignedCert(
		ctx,
		privateKey,
		[]pulumi.Resource{privateKey},
	)
	if err != nil {
		return nil, err
	}

	clusterIssuerCertBundleSecret, err := deployClusterIssuerTLSSecret(
		ctx,
		provider,
		namespace,
		privateKey,
		cert,
		[]pulumi.Resource{namespace, privateKey, cert},
	)
	if err != nil {
		return nil, err
	}

	clusterIssuer, err := deployCertManagerInternalClusterIssuer(
		ctx,
		provider,
		namespace,
		clusterIssuerCertBundleSecret,
		[]pulumi.Resource{namespace, chart, clusterIssuerCertBundleSecret},
	)
	if err != nil {
		return nil, err
	}

	return []pulumi.Resource{
		namespace,
		chart,
		privateKey,
		cert,
		clusterIssuerCertBundleSecret,
		clusterIssuer,
	}, nil
}

// deployCertManagerHelmChart deploys `cert-manager` using its Helm chart.
func deployCertManagerHelmChart(
	ctx *pulumi.Context,
	namespace *corev1.Namespace,
	provider *kubernetes.Provider,
	deps []pulumi.Resource,
) (*helmv4.Chart, error) {
	kind := config.GetBool(ctx, configKind)

	chartArgs := newCertManagerHelmChartArgs(
		namespace,
		kind,
	)
	certManager, err := helmv4.NewChart(
		ctx,
		chartName,
		chartArgs,
		pulumi.Provider(provider),
		pulumi.DependsOn(deps),
	)
	if err != nil {
		return nil, err
	}
	return certManager, nil
}

// newCertManagerHelmChartArgs creates a Helm chart arguments. The Helm chart args can then be used by a Pulumi program
// to deploy cert-manager.
//
// [kind] controls the [enableGatewayAPI] value by disabling the GatewayAPI if the chart is deployed locally.
func newCertManagerHelmChartArgs(
	namespace *corev1.Namespace,
	kind bool,
) *helmv4.ChartArgs {
	return &helmv4.ChartArgs{
		Chart: pulumi.String(chartName),
		RepositoryOpts: &helmv4.RepositoryOptsArgs{
			Repo: pulumi.String(chartRepo),
		},
		Namespace: namespace.Metadata.Name(),
		Version:   pulumi.String(chartVersion),
		Values: pulumi.Map{
			"global": pulumi.Map{
				"leaderElection": pulumi.Map{
					"namespace": namespace.Metadata.Name(),
				},
			},
			"image": pulumi.Map{
				"tag": pulumi.String(fmt.Sprintf("v%s", chartVersion)),
			},
			"replicaCount": pulumi.Int(3),
			"podDisruptionBudget": pulumi.Map{
				"enabled": pulumi.Bool(true),
			},
			"strategy": pulumi.Map{
				"type": pulumi.String("RollingUpdate"),
				"rollingUpdate": pulumi.Map{
					"maxSurge":       pulumi.Int(1),
					"maxUnavailable": pulumi.Int(1),
				},
			},
			"config": pulumi.Map{
				"apiVersion": pulumi.String("controller.config.cert-manager.io/v1alpha1"),
				"kind":       pulumi.String("ControllerConfiguration"),
				"logging": pulumi.Map{
					"format":    pulumi.String("json"),
					"verbosity": pulumi.Int(5), // Debug
				},
				"leaderElectionConfig": pulumi.Map{
					"namespace": namespace.Metadata.Name(),
				},
				"enableGatewayAPI": pulumi.Bool(!kind),
			},
			"crds": pulumi.Map{
				"enabled": pulumi.Bool(true),
				"keep":    pulumi.Bool(false),
			},
			"prometheus": pulumi.Map{
				"enabled": pulumi.Bool(true),
			},
			"cainjector": pulumi.Map{
				"replicaCount": pulumi.Int(3),
				"config": pulumi.Map{
					"apiVersion": pulumi.String("cainjector.config.cert-manager.io/v1alpha1"),
					"kind":       pulumi.String("CAInjectorConfiguration"),
					"logging": pulumi.Map{
						"format":    pulumi.String("json"),
						"verbosity": pulumi.Int(5), // Debug
					},
					"leaderElectionConfig": pulumi.Map{
						"namespace": namespace.Metadata.Name(),
					},
				},
				"strategy": pulumi.Map{
					"type": pulumi.String("RollingUpdate"),
					"rollingUpdate": pulumi.Map{
						"maxSurge":       pulumi.Int(0),
						"maxUnavailable": pulumi.Int(1),
					},
				},
				"podDisruptionBudget": pulumi.Map{
					"enabled": pulumi.Bool(true),
				},
			},
			"webhook": pulumi.Map{
				"replicaCount": pulumi.Int(3),
				"config": pulumi.Map{
					"apiVersion": pulumi.String("webhook.config.cert-manager.io/v1alpha1"),
					"kind":       pulumi.String("WebhookConfiguration"),
					"logging": pulumi.Map{
						"format":    pulumi.String("json"),
						"verbosity": pulumi.Int(5), // Debug
					},
				},
				"strategy": pulumi.Map{
					"type": pulumi.String("RollingUpdate"),
					"rollingUpdate": pulumi.Map{
						"maxSurge":       pulumi.Int(0),
						"maxUnavailable": pulumi.Int(1),
					},
				},
				"podDisruptionBudget": pulumi.Map{
					"enabled": pulumi.Bool(true),
				},
			},
		},
	}
}
