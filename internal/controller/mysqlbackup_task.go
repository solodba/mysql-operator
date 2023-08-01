package controller

import (
	"time"

	operatorcodehorsecomv1beta1 "codehorse.com/mysql-operator/api/v1beta1"
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
		delay := r.GetMysqlBackupDelay(mysqlBackup.Spec.StartTime)
		if delay.Hours() >= 1 {
			logger.L().Info().Msgf("[%s] mysql backup start after [%.1f] hours start!", mysqlBackup.Name, delay.Hours())
		} else {
			logger.L().Info().Msgf("[%s] mysql backup start after [%.1f] minute start!", mysqlBackup.Name, delay.Minutes())
		}
		// 更新备份任务状态
		mysqlBackup.Status.Active = true
		mysqlBackup.Status.NextTime = r.GetMysqlBackupNextTime(delay).Unix()
		r.UpdateMysqlBackupStatus(mysqlBackup)
		// 开始定时启动备份任务
		ticker := time.NewTicker(delay)
		if r.Tickers == nil {
			r.Tickers = make([]*time.Ticker, 0)
		}
		r.Tickers = append(r.Tickers, ticker)
		r.Wg.Add(1)
		go func(mysqlBackup *operatorcodehorsecomv1beta1.MysqlBackup) {
			defer r.Wg.Done()
			// 循环执行执行任务
			for {
				<-ticker.C
				// 重置定时器
				ticker.Reset(time.Minute * time.Duration(mysqlBackup.Spec.Period))
				// 更新状态
				mysqlBackup.Status.Active = true
				mysqlBackup.Status.NextTime = r.GetMysqlBackupNextTime(time.Minute * time.Duration(mysqlBackup.Spec.Period)).Unix()
				// 开始备份
				// todo

				// 更新状态
				r.UpdateMysqlBackupStatus(mysqlBackup)
			}
		}(mysqlBackup)
		r.Wg.Wait()
	}

}

// 停止任务
func (r *MysqlBackupReconciler) StopTask() {

}
