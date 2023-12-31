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

package controller

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	operatorcodehorsecomv1beta1 "codehorse.com/mysql-operator/api/v1beta1"
	"codehorse.com/mysql-operator/common/logger"
)

// MysqlBackupReconciler reconciles a MysqlBackup object
type MysqlBackupReconciler struct {
	client.Client
	Scheme           *runtime.Scheme
	MysqlBackupQueue map[string]*operatorcodehorsecomv1beta1.MysqlBackup
	Lock             sync.Mutex
	Tickers          []*time.Ticker
	Wg               sync.WaitGroup
}

//+kubebuilder:rbac:groups=operator.codehorse.com,resources=mysqlbackups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=operator.codehorse.com,resources=mysqlbackups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=operator.codehorse.com,resources=mysqlbackups/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MysqlBackup object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *MysqlBackupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	// 查找mysql备份任务
	k8sMysqlBackup := operatorcodehorsecomv1beta1.NewMysqlBackup()
	err := r.Client.Get(ctx, req.NamespacedName, k8sMysqlBackup)
	if err != nil {
		if errors.IsNotFound(err) {
			// 查找不到mysql备份任务, 说明任务已经停止, 则删除停止任务
			logger.L().Error().Msgf("[%s] is already stopped!", k8sMysqlBackup.Name)
			r.DeleteMysqlBackupQueue(k8sMysqlBackup)
			return ctrl.Result{}, err
		}
		// mysql备份任务异常
		logger.L().Error().Msgf("[%s] is already abnormal!", k8sMysqlBackup.Name)
		return ctrl.Result{}, err
	}
	// 对比任务信息是否改变
	if lastMysqlBackup, ok := r.MysqlBackupQueue[k8sMysqlBackup.Name]; ok {
		if reflect.DeepEqual(k8sMysqlBackup.Spec, lastMysqlBackup.Spec) {
			return ctrl.Result{}, fmt.Errorf("[%s] information is not changed", k8sMysqlBackup.Name)
		}
	}
	// 添加备份任务到队列
	r.AddMysqlBackupQueue(k8sMysqlBackup)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MysqlBackupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorcodehorsecomv1beta1.MysqlBackup{}).
		Complete(r)
}
