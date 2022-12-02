package controllers

import (
	"context"
	"fmt"

	bgrouter "github.com/yashirook/bgrouter/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func nameWithSuffix(suffix string, bgr *bgrouter.BGRouter) string {
	return fmt.Sprintf("%s-%s", bgr.Name, suffix)
}

func fetchOject(client client.Client, namespace string, name string, obj client.Object) error {
	return client.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, obj)
}

func isOjectFound(client client.Client, namespace string, name string, obj client.Object) bool {
	return !apierrors.IsNotFound(fetchOject(client, namespace, name, obj))
}
