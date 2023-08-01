package controller

import (
	"strconv"
	"strings"
	"time"
)

// 计算备份时间
func (r *MysqlBackupReconciler) GetMysqlBackupDelay(startTime string) time.Duration {
	startTimeArr := strings.Split(startTime, ":")
	expectedHour, _ := strconv.Atoi(startTimeArr[0])
	expectedMin, _ := strconv.Atoi(startTimeArr[1])
	now := time.Now().Truncate(time.Second)
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	otherStartDate := startDate.Add(time.Hour * time.Duration(24))
	expectedStartTime := time.Hour*time.Duration(expectedHour) + time.Minute*time.Duration(expectedMin)
	currentTime := time.Hour*time.Duration(now.Hour()) + time.Minute*time.Duration(now.Minute())
	var seconds int64
	if expectedStartTime >= currentTime {
		seconds = int64(startDate.Add(expectedStartTime).Sub(now).Seconds())
	} else {
		seconds = int64(otherStartDate.Add(expectedStartTime).Sub(now).Seconds())
	}
	return time.Second * time.Duration(seconds)
}
