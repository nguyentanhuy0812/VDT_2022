package main

import (
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

type ChartSpec struct {

	ReleaseName string `json:"release"`
	ChartName string `json: "chart"`
	Namespace string `json: "namespace"`

}

type HelmClient struct {
		// Settings defines the environment settings of a client.
		Settings  *cli.EnvSettings
		Providers getter.Providers
	
		// ActionConfig is the helm action configuration.
		ActionConfig *action.Configuration
	
		DebugLog action.DebugLog
		// contains filtered or unexported fields
}

func main() {
	
	fmt.Println("Uninstall done!")
}

func (c *HelmClient) uninstallReleaseByName(name string) error {
	client := action.NewUninstall(c.ActionConfig)

	resp, err := client.Run(name)
	if err != nil {
		return err
	}

	c.DebugLog("release uninstalled, response: %v", resp)

	return nil
}

func (c *HelmClient) UninstallReleaseByName(name string) error {
	return c.uninstallReleaseByName(name)
}

// func (c *HelmClient) UninstallRelease(spec *ChartSpec) error {

// 	client := action.NewUninstall(c.ActionConfig)

// 	mergeUninstallReleaseOptions(spec, client)

// 	resp, err := client.Run(spec.ReleaseName)

// 	if err != nil {
// 		return err
// 	}

// 	c.DebugLog("Release uninstalled, response: %v", resp)

// 	return nil
// }

// func mergeUninstallReleaseOptions(chartSpec *ChartSpec, uninstallReleaseOptions *action.Uninstall) {
// 	uninstallReleaseOptions.DisableHooks = chartSpec.DisableHooks
// 	uninstallReleaseOptions.Timeout = chartSpec.Timeout
// }


