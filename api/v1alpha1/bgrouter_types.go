/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VirtualServiceConfig struct {
	// Name is the name for VirtualService managed by bgrouter controller.
	// +kubebuilder:validation:Required
	Name string `json:"name,omitempty"`

	// TargetServiceName is the FQDN for accessing target Service resource.
	// If you wanna access the service named "foo-service" in "foo" namespace,
	// you specifies "foo-service.foo.svc.cluster.local".
	// +kubebuilder:validation:Required
	TargetServiceName string `json:"targetServiceName,omitempty"`

	// HostsForInternalTraffic is list of hosts for in-cluster traffic.
	// +kubebuilder:validation:Required
	HostsForInClusterTraffic []string `json:"hostsForInClusterTraffic,omitempty"`

	// HostsForOutsideClusterTraffic is list of hosts for outside cluster
	// +optional
	HostsForOutsideClusterTraffic []string `json:"hostsForOutsideClusterTraffic,omitempty"`
}

// BGRouterSpec defines the desired state of BGRouter
type BGRouterSpec struct {
	// ActiveColor specifies the color label (blue or green) assigned to the pods to
	// witch you want to direct traffic.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=blue;green
	ActiveColor string `json:"activeColor,omitempty"`

	// ActiveReplicas is the number of pods with active color labels.
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=0
	// +optional
	ActiveReplicas int32 `json:"activeReplicas,omitempty"`

	// IdleReplicas is the number of pods with idle color labels.
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=0
	// +optional
	IdleReplicas int32 `json:"idleReplicas,omitempty"`

	// HpaBaseName is suffix using for name of hpa resource managed by bgrouter controller.
	// +optional
	HpaBaseName string `json:"hpaBaseName,omitempty"`

	// DeploymentBaseName is suffix using for name of deployment.
	// Deployment name have to be ${DeploymentBaseName}-blue or ${DeploymentBaseName}-green format.
	// +kubebuilder:validation:Required
	DeploymentBaseName string `json:"deploymentBaseName,omitempty"`

	// VirtualServiceConfig is set of information for building VirtualService resource.
	// +kubebuilder:validation:Required
	VirtualServiceConfig VirtualServiceConfig `json:"virtualServiceConfig,omitempty"`
}

// BGRouterStatus defines the observed state of BGRouter
type BGRouterStatus struct {
	// CurrentActiveColor represent for current active color.
	// +kubebuilder:validation:Enum=blue;green
	CurrentActiveColor string `json:"currentActiveColor,omitempty"`

	// Progress describe the progress of swithing operation.
	// +kubebuilder:validation:Enum=DONE;ONGOING
	Progress string `json:"progress,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="DESIRED_ACTIVE_COLOR",type="string",JSONPath=".spec.activeColor"
//+kubebuilder:printcolumn:name="PROGRESS",type="string",JSONPath=".status.progress"

// BGRouter is the Schema for the bgrouters API
type BGRouter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BGRouterSpec   `json:"spec,omitempty"`
	Status BGRouterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BGRouterList contains a list of BGRouter
type BGRouterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BGRouter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BGRouter{}, &BGRouterList{})
}
