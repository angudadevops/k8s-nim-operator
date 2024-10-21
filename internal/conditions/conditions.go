/*
Copyright 2024.

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

package conditions

import (
	"context"
	"fmt"

	"github.com/NVIDIA/k8s-nim-operator/api/apps/v1alpha1"
	appsv1alpha1 "github.com/NVIDIA/k8s-nim-operator/api/apps/v1alpha1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// Ready indicates that the service is ready
	Ready = "Ready"
	// NotReady indicates that the service is not yet ready
	NotReady = "NotReady"
	// Failed indicates that the service has failed
	Failed = "Failed"
	// ReasonServiceAccountFailed indicates that the creation of serviceaccount has failed
	ReasonServiceAccountFailed = "ServiceAccountFailed"
	// ReasonRoleFailed indicates that the creation of serviceaccount has failed
	ReasonRoleFailed = "RoleFailed"
	// ReasonRoleBindingFailed indicates that the creation of rolebinding has failed
	ReasonRoleBindingFailed = "RoleBindingFailed"
	// ReasonServiceFailed indicates that the creation of service has failed
	ReasonServiceFailed = "ServiceFailed"
	// ReasonIngressFailed indicates that the creation of ingress has failed
	ReasonIngressFailed = "IngressFailed"
	// ReasonHPAFailed indicates that the creation of hpa has failed
	ReasonHPAFailed = "HPAFailed"
	// ReasonSCCFailed indicates that the creation of scc has failed
	ReasonSCCFailed = "SCCFailed"
	// ReasonServiceMonitorFailed indicates that the creation of Service Monitor has failed
	ReasonServiceMonitorFailed = "ServiceMonitorFailed"
	// ReasonDeploymentFailed indicates that the creation of deployment has failed
	ReasonDeploymentFailed = "DeploymentFailed"
	// ReasonStatefulSetFailed indicates that the creation of statefulset has failed
	ReasonStatefulSetFailed = "StatefulsetFailed"
)

// Updater is the condition updater

type Updater interface {
	SetConditionsReady(ctx context.Context, cr client.Object, reason, message string) error
	SetConditionsNotReady(ctx context.Context, cr client.Object, reason, message string) error
	SetConditionsFailed(ctx context.Context, cr client.Object, reason, message string) error
}

type updater struct {
	client client.Client
}

// NewUpdater returns an instance of updater
func NewUpdater(c client.Client) Updater {
	return &updater{client: c}
}

func (u *updater) SetConditionsReady(ctx context.Context, obj client.Object, reason, message string) error {
	if cr, ok := obj.(*appsv1alpha1.NIMService); ok {
		return u.SetConditionsReadyNIMService(ctx, cr, reason, message)
	} else if gr, ok := obj.(*appsv1alpha1.NemoGuardrail); ok {
		return u.SetConditionsReadyNemoGuardrail(ctx, gr, reason, message)
	}
	return fmt.Errorf("unknown CRD type for %v", obj)
}

func (u *updater) SetConditionsReadyNIMService(ctx context.Context, cr *appsv1alpha1.NIMService, reason, message string) error {
	meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
		Type:    Ready,
		Status:  metav1.ConditionTrue,
		Reason:  reason,
		Message: message,
	})

	meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
		Type:   Failed,
		Status: metav1.ConditionFalse,
		Reason: Ready,
	})
	cr.Status.State = v1alpha1.NIMServiceStatusReady
	return u.updateNIMServiceStatus(ctx, cr)
}

func (u *updater) SetConditionsReadyNemoGuardrail(ctx context.Context, gr *appsv1alpha1.NemoGuardrail, reason, message string) error {
	meta.SetStatusCondition(&gr.Status.Conditions, metav1.Condition{
		Type:    Ready,
		Status:  metav1.ConditionTrue,
		Reason:  reason,
		Message: message,
	})

	meta.SetStatusCondition(&gr.Status.Conditions, metav1.Condition{
		Type:   Failed,
		Status: metav1.ConditionFalse,
		Reason: Ready,
	})
	gr.Status.State = v1alpha1.NIMServiceStatusReady
	return u.updateNemoGuardrailStatus(ctx, gr)
}

func (u *updater) SetConditionsNotReadyNIMService(ctx context.Context, cr *appsv1alpha1.NIMService, reason, message string) error {
	meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
		Type:    Ready,
		Status:  metav1.ConditionFalse,
		Reason:  reason,
		Message: message,
	})

	meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
		Type:    Failed,
		Status:  metav1.ConditionFalse,
		Reason:  Ready,
		Message: message,
	})
	cr.Status.State = v1alpha1.NIMServiceStatusNotReady
	return u.updateNIMServiceStatus(ctx, cr)
}

func (u *updater) SetConditionsNotReadyNemoGuardrail(ctx context.Context, gr *appsv1alpha1.NemoGuardrail, reason, message string) error {
	meta.SetStatusCondition(&gr.Status.Conditions, metav1.Condition{
		Type:    Ready,
		Status:  metav1.ConditionFalse,
		Reason:  reason,
		Message: message,
	})

	meta.SetStatusCondition(&gr.Status.Conditions, metav1.Condition{
		Type:    Failed,
		Status:  metav1.ConditionFalse,
		Reason:  Ready,
		Message: message,
	})
	gr.Status.State = v1alpha1.NemoGuardrailStatusNotReady
	return u.updateNemoGuardrailStatus(ctx, gr)
}

func (u *updater) SetConditionsNotReady(ctx context.Context, obj client.Object, reason, message string) error {
	if cr, ok := obj.(*appsv1alpha1.NIMService); ok {
		return u.SetConditionsNotReadyNIMService(ctx, cr, reason, message)
	} else if gr, ok := obj.(*appsv1alpha1.NemoGuardrail); ok {
		return u.SetConditionsNotReadyNemoGuardrail(ctx, gr, reason, message)
	}
	return fmt.Errorf("unknown CRD type for %v", obj)
}

func (u *updater) SetConditionsFailed(ctx context.Context, obj client.Object, reason, message string) error {
	if cr, ok := obj.(*appsv1alpha1.NIMService); ok {
		return u.SetConditionsFailedNIMService(ctx, cr, reason, message)
	} else if gr, ok := obj.(*appsv1alpha1.NemoGuardrail); ok {
		return u.SetConditionsFailedNemoGuardrail(ctx, gr, reason, message)
	}
	return fmt.Errorf("unknown CRD type for %v", obj)
}

func (u *updater) SetConditionsFailedNIMService(ctx context.Context, cr *appsv1alpha1.NIMService, reason, message string) error {
	meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
		Type:   Ready,
		Status: metav1.ConditionFalse,
		Reason: Failed,
	})

	meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
		Type:    Failed,
		Status:  metav1.ConditionTrue,
		Reason:  reason,
		Message: message,
	})
	cr.Status.State = v1alpha1.NIMServiceStatusFailed
	return u.updateNIMServiceStatus(ctx, cr)
}

func (u *updater) SetConditionsFailedNemoGuardrail(ctx context.Context, cr *appsv1alpha1.NemoGuardrail, reason, message string) error {
	meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
		Type:   Ready,
		Status: metav1.ConditionFalse,
		Reason: Failed,
	})

	meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
		Type:    Failed,
		Status:  metav1.ConditionTrue,
		Reason:  reason,
		Message: message,
	})
	cr.Status.State = v1alpha1.NemoGuardrailStatusFailed
	return u.updateNemoGuardrailStatus(ctx, cr)
}

func (u *updater) updateNIMServiceStatus(ctx context.Context, cr *appsv1alpha1.NIMService) error {

	obj := &appsv1alpha1.NIMService{}
	errGet := u.client.Get(ctx, types.NamespacedName{Name: cr.Name, Namespace: cr.GetNamespace()}, obj)
	if errGet != nil {
		return errGet
	}
	obj.Status = cr.Status
	if err := u.client.Status().Update(ctx, obj); err != nil {
		return err
	}
	return nil
}

func (u *updater) updateNemoGuardrailStatus(ctx context.Context, cr *appsv1alpha1.NemoGuardrail) error {
	obj := &appsv1alpha1.NemoGuardrail{}
	errGet := u.client.Get(ctx, types.NamespacedName{Name: cr.Name, Namespace: cr.GetNamespace()}, obj)
	if errGet != nil {
		return errGet
	}
	obj.Status = cr.Status
	if err := u.client.Status().Update(ctx, obj); err != nil {
		return err
	}
	return nil
}

// UpdateCondition updates the given condition into the conditions list
func UpdateCondition(conditions *[]metav1.Condition, conditionType string, status metav1.ConditionStatus, reason, message string) {
	for i := range *conditions {
		if (*conditions)[i].Type == conditionType {
			// existing condition
			(*conditions)[i].Status = status
			(*conditions)[i].LastTransitionTime = metav1.Now()
			(*conditions)[i].Reason = reason
			(*conditions)[i].Message = message
			// condition updated
			return
		}
	}
	// new condition
	*conditions = append(*conditions, metav1.Condition{
		Type:               conditionType,
		Status:             status,
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	})
	// condition updated
}

func IfPresentUpdateCondition(conditions *[]metav1.Condition, conditionType string, status metav1.ConditionStatus, reason, message string) {
	for i := range *conditions {
		if (*conditions)[i].Type == conditionType {
			// existing condition
			(*conditions)[i].Status = status
			(*conditions)[i].LastTransitionTime = metav1.Now()
			(*conditions)[i].Reason = reason
			(*conditions)[i].Message = message
			// condition updated
			return
		}
	}
}
