package controllers

import (
	"context"
	"fmt"

	bgrouter "github.com/yashirook/bgrouter/api/v1alpha1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logr "sigs.k8s.io/controller-runtime/pkg/log"
)

func newHorizontalPodAutoscaler(bgr *bgrouter.BGRouter) *autoscalingv1.HorizontalPodAutoscaler {
	return &autoscalingv1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bgr.Name,
			Namespace: bgr.Namespace,
		},
	}
}

func newHorizontalPodAutoscalerWithName(name string, bgr *bgrouter.BGRouter) *autoscalingv1.HorizontalPodAutoscaler {
	hpa := newHorizontalPodAutoscaler(bgr)
	hpa.ObjectMeta.Name = name

	lbls := hpa.ObjectMeta.Labels
	lbls["ArgoCDKeyName"] = name
	hpa.ObjectMeta.Labels = lbls

	return hpa
}

func newHorizontalPodAutoscalerWithSuffix(suffix string, bgr *bgrouter.BGRouter) *autoscalingv1.HorizontalPodAutoscaler {
	return newHorizontalPodAutoscalerWithName(nameWithSuffix(suffix, bgr), bgr)
}

func (r *BGRouterReconciler) reconcileHPA(ctx context.Context, bgr *bgrouter.BGRouter) error {
	logger := logr.FromContext(ctx)
	colors := []string{"blue", "green"}

	hpaBaseName := bgr.ObjectMeta.Name + "-hpa"
	if bgr.Spec.HpaBaseName != "" {
		hpaBaseName = bgr.Spec.HpaBaseName
	}

	for _, color := range colors {
		hpa := &autoscalingv1.HorizontalPodAutoscaler{}
		hpa.SetName(hpaBaseName + "-" + color)
		hpa.SetNamespace(bgr.Namespace)

		op, err := ctrl.CreateOrUpdate(ctx, r.Client, hpa, func() error {
			hpa.Spec.MaxReplicas = bgr.Spec.ActiveReplicas
			hpa.Spec.MinReplicas = &bgr.Spec.ActiveReplicas

			var tcup int32 = 50
			hpa.Spec.TargetCPUUtilizationPercentage = &tcup

			hpa.Spec.ScaleTargetRef = autoscalingv1.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       fmt.Sprintf("%s-%s", bgr.Spec.DeploymentBaseName, color),
			}
			return ctrl.SetControllerReference(bgr, hpa, r.Scheme)
		})

		if err != nil {
			logger.Error(err, "unable to create or update HorizontalPodAutoscaler.")
			return err
		}

		if op != controllerutil.OperationResultNone {
			logger.Info("reconcile HorizontalPodAutoscaler successfully", "op", op)
		}
	}
	return nil
}
