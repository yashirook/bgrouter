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

package controllers

import (
	"context"

	autoscalingv1 "k8s.io/api/autoscaling/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logr "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	bgrouter "github.com/yashirook/bgrouter/api/v1alpha1"
)

// BGRouterReconciler reconciles a BGRouter object
type BGRouterReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Recourder record.EventRecorder
}

var log = logr.Log.WithName("controller_bgrouter")

//+kubebuilder:rbac:groups=bgrouter.yashirook.github.io,resources=bgrouters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=bgrouter.yashirook.github.io,resources=bgrouters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=bgrouter.yashirook.github.io,resources=bgrouters/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch
//+kubebuilder:rbac:groups=autoscaling,resources=horizontalpodautoscalers,verbs=*

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the BGRouter object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *BGRouterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logr.FromContext(ctx, "namespace", req.Namespace, "name", req.Name)
	log.Info("Reconciling BGRouter")

	bgrouter := &bgrouter.BGRouter{}
	if err := r.Get(ctx, req.NamespacedName, bgrouter); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}
	// TODO(user): your logic here

	if err := r.reconcileHPA(ctx, bgrouter); err != nil {
		return reconcile.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BGRouterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bgrouter.BGRouter{}).
		Owns(&autoscalingv1.HorizontalPodAutoscaler{}).
		Complete(r)
}
