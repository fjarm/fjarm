package main

import (
	"github.com/fjarm/fjarm/infrastructure/internal/pkg/certmanager/certmanagerv1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	k8sProviderLogicalNamePrefix = "kubernetes-cluster-provider"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		k8sProvider, err := kubernetes.NewProvider(ctx, k8sProviderLogicalNamePrefix, &kubernetes.ProviderArgs{
			//RenderYamlToDirectory: pulumi.String("yaml"),
		})
		if err != nil {
			return err
		}
		_, err = certmanagerv1.DeployCertManager(
			ctx,
			k8sProvider,
		)
		if err != nil {
			return err
		}
		return nil
	})
}
