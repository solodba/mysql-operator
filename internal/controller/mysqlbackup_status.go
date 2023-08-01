package controller

import (
	"context"

	operatorcodehorsecomv1beta1 "codehorse.com/mysql-operator/api/v1beta1"
	"codehorse.com/mysql-operator/common/logger"
	"k8s.io/apimachinery/pkg/types"
)

// 更新mysql备份任务状态
func (r *MysqlBackupReconciler) UpdateMysqlBackupStatus(mysqlBackup *operatorcodehorsecomv1beta1.MysqlBackup) {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	// 查找mysql备份任务
	ctx := context.TODO()
	namespacedName := types.NamespacedName{
		Namespace: mysqlBackup.Namespace,
		Name:      mysqlBackup.Name,
	}
	k8sMysqlBackup := operatorcodehorsecomv1beta1.NewMysqlBackup()
	err := r.Client.Get(ctx, namespacedName, k8sMysqlBackup)
	if err != nil {
		logger.L().Error().Msgf("[%s] mysql backup task is not found!", mysqlBackup.Name)
		return
	}
	// 更新状态
	k8sMysqlBackup.Status = mysqlBackup.Status
	err = r.Status().Update(ctx, k8sMysqlBackup)
	if err != nil {
		logger.L().Error().Msgf("update [%s] mysql backup task status failed, err: %s!", mysqlBackup.Name, err.Error())
	}
	return
}
