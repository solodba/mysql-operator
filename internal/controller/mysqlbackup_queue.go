package controller

import (
	operatorcodehorsecomv1beta1 "codehorse.com/mysql-operator/api/v1beta1"
)

// 从mysql备份队列中删除备份任务
func (r *MysqlBackupReconciler) DeleteMysqlBackupQueue(mysqlBackup *operatorcodehorsecomv1beta1.MysqlBackup) {
	delete(r.MysqlBackupQueue, mysqlBackup.Name)
	r.StopTask(mysqlBackup)
	go r.StartTask(mysqlBackup)
}

// 添加mysql备份任务到备份队列中
func (r *MysqlBackupReconciler) AddMysqlBackupQueue(mysqlBackup *operatorcodehorsecomv1beta1.MysqlBackup) {
	if r.MysqlBackupQueue == nil {
		r.MysqlBackupQueue = make(map[string]*operatorcodehorsecomv1beta1.MysqlBackup)
	}
	r.MysqlBackupQueue[mysqlBackup.Name] = mysqlBackup
	r.StopTask(mysqlBackup)
	go r.StartTask(mysqlBackup)
}
