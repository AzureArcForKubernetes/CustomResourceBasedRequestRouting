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

package controllers

import (
	"context"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	arcv1beta1 "github.com/AzureArcForKubernetes/CustomResourceBasedRequestRouting/api/v1beta1"
)

// RequestRoutingRulesReconciler reconciles a RequestRoutingRules object
type RequestRoutingRulesReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var (
	RoutingRules map[string]*arcv1beta1.RequestRoutingRules
)

//+kubebuilder:rbac:groups=arc.azure.com,resources=requestroutingrules,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=arc.azure.com,resources=requestroutingrules/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=arc.azure.com,resources=requestroutingrules/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RequestRoutingRules object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *RequestRoutingRulesReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	instance := &arcv1beta1.RequestRoutingRules{}
	err := r.Client.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		log.Log.Info("RequestRoutingRules resource not found. Ignoring since object must be deleted.")
		if RoutingRules[req.Name+"-"+req.Namespace] != nil {
			delete(RoutingRules, req.Name)
		}
		return ctrl.Result{}, nil
	}
	RoutingRules[req.Name+"-"+req.Namespace] = instance.DeepCopy()
	return ctrl.Result{}, nil
}

func ResolveProxyEndpoint(incomingUrl string) string {
	for _, rule := range RoutingRules {
		if strings.Contains(incomingUrl, rule.Spec.ResourceNameSubstring) {
			if rule.Spec.IsPublicEndpoint {
				return rule.Spec.DNSName
			}
		}
	}
	return ""
}

// SetupWithManager sets up the controller with the Manager.
func (r *RequestRoutingRulesReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&arcv1beta1.RequestRoutingRules{}).
		Complete(r)
}
