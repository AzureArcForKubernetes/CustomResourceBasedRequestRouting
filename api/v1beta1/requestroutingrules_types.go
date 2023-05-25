/*
Copyright 2023.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RequestRoutingRulesSpec defines the desired state of RequestRoutingRules
type RequestRoutingRulesSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of RequestRoutingRules. Edit requestroutingrules_types.go to remove/update
	DNSName                string `json:"dnsName,omitempty"`
	IsPublicEndpoint       bool   `json:"isPublicEndpoint,omitempty"`
	IsPortForwardingNeeded bool   `json:"isPortForwardingNeeded,omitempty"`
	KubeConfigSecretName   string `json:"kubeConfigSecretName,omitempty"`
	ResourceNameSubstring  string `json:"resourceNameSubstring,omitempty"`
}

// RequestRoutingRulesStatus defines the observed state of RequestRoutingRules
type RequestRoutingRulesStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RequestRoutingRules is the Schema for the requestroutingrules API
type RequestRoutingRules struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RequestRoutingRulesSpec   `json:"spec,omitempty"`
	Status RequestRoutingRulesStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RequestRoutingRulesList contains a list of RequestRoutingRules
type RequestRoutingRulesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RequestRoutingRules `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RequestRoutingRules{}, &RequestRoutingRulesList{})
}
