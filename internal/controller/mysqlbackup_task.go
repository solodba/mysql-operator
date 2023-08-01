package controller

import (
	"codehorse.com/mysql-operator/common/logger"
)

// 开启任务
func (r *MysqlBackupReconciler) StartTask() {
	// 遍历mysql备份队列
	for _, mysqlBackup := range r.MysqlBackupQueue {
		if !mysqlBackup.Spec.Enable {
			logger.L().Info().Msgf("[%s] is not enable!", mysqlBackup.Name)
			mysqlBackup.Status.Active = false
			// 更新备份任务状态
			r.UpdateMysqlBackupStatus(mysqlBackup)
			continue
		}
		// 计算备份时间

	}

}

// 停止任务
func (r *MysqlBackupReconciler) StopTask() {

}
