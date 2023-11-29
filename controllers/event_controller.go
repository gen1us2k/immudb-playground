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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	immudbURL = "https://vault.immudb.io/ics/api/v1/ledger/default/collection/default/document"
)

type (
	// EventReconciler reconciles a Event object
	EventReconciler struct {
		client.Client
		Scheme       *runtime.Scheme
		immudbAPIKey string
		httpClient   *http.Client
	}
	Event struct {
		Event *corev1.Event `json:"event"`
		ID    string        `json:"id"`
	}
)

//+kubebuilder:rbac:groups=immudb.com,resources=events,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=immudb.com,resources=events/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=immudb.com,resources=events/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *EventReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	e := &corev1.Event{}
	err := r.Get(ctx, req.NamespacedName, e)
	if err != nil {
		// NotFound cannot be fixed by requeuing so ignore it. During background
		// deletion, we receive delete events from cluster's dependents after
		// cluster is deleted.
		if err = client.IgnoreNotFound(err); err != nil {
			logger.Error(err, "unable to fetch event")
		}
		return reconcile.Result{}, err
	}
	if err := r.storeEvent(e); err != nil {
		logger.Error(err, "failed storing event")
	}

	return ctrl.Result{}, nil
}

func (r *EventReconciler) storeEvent(e *corev1.Event) error {
	ie := &Event{
		ID:    string(e.ObjectMeta.UID),
		Event: e,
	}
	data, err := json.Marshal(ie)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, immudbURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("X-API-Key", r.immudbAPIKey)
	req.Header.Set("Content-type", "application/json")
	res, err := r.httpClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		fmt.Println(string(body))
		fmt.Println(string(data))
		return fmt.Errorf("failed saving event. Received %d code", res.StatusCode)
	}
	return nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *EventReconciler) SetupWithManager(mgr ctrl.Manager) error {
	key := os.Getenv("IMMUDB_API_KEY")
	if key == "" {
		return errors.New("IMMUDB_API_KEY is not set")
	}
	r.immudbAPIKey = key
	r.httpClient = &http.Client{}
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Event{}).
		Complete(r)
}
