package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

)

const (
	// bien status

)

type HelmChartSpec struct {

	ClusterSelector metav1.LabelSelector `json:"clusterSelector"`


	Resources []ResourceRef `json:"resources,omitempty"`
}

type HelmChartStatus struct {

	Status string `json:"status,omitempty"`

}

type ResourceRef struct {

	TargetNamespace string `json:"targetNamespace,omitempty"`
	URL string `json:"URL"`
	Chart string `json:"chart,omitempty"`
	Version string `json:"version,omitempty"`

	Repo string `json:"repo,omitempty"`
	Values string `json:"values,omitempty"`
}

type HelmChart struct {
		
	metav1.TypeMeta `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec HelmChartSpec `json:"spec,omitempty"`
	Status HelmChartStatus `json:"status,omitempty"`
}

type HelmChartList struct {

	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items			[]HelmChart `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HelmChart{}, &HelmChartList{})
}