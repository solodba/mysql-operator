package controller

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	operatorcodehorsecomv1beta1 "codehorse.com/mysql-operator/api/v1beta1"
	"codehorse.com/mysql-operator/common/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// mysql备份方法
func (r *MysqlBackupReconciler) MysqlDataDumpAndUpload(mysqlBackup *operatorcodehorsecomv1beta1.MysqlBackup) error {
	mysqlBackupDir, err := r.CreateMysqlBackupDir(mysqlBackup)
	if err != nil {
		return err
	}
	mysqlBackupFileName, err := r.MysqlDataDump(mysqlBackupDir, mysqlBackup)
	if err != nil {
		return err
	}
	err = r.UploadMysqlDumpFileToMinIO(mysqlBackupFileName, mysqlBackup)
	if err != nil {
		return err
	}
	return nil
}

// 创建备份目录
func (r *MysqlBackupReconciler) CreateMysqlBackupDir(mysqlBackup *operatorcodehorsecomv1beta1.MysqlBackup) (string, error) {
	now := time.Now()
	monthDay := fmt.Sprintf("%02d_%02d", now.Month(), now.Day())
	mysqlBackupDir := fmt.Sprintf("/tmp/%s/%s", mysqlBackup.Name, monthDay)
	if _, err := os.Stat(mysqlBackupDir); err != nil {
		if err := os.MkdirAll(mysqlBackupDir, 0700); err != nil {
			return "", fmt.Errorf("[%s] mysql backup dir create failed, err: %s", mysqlBackupDir, err.Error())
		}
		logger.L().Info().Msgf("[%s] mysql backup dir create successful!", mysqlBackupDir)
	}
	return mysqlBackupDir, nil
}

// mysql备份任务
func (r *MysqlBackupReconciler) MysqlDataDump(mysqlBackupDir string, mysqlBackup *operatorcodehorsecomv1beta1.MysqlBackup) (string, error) {
	fileList, err := os.ReadDir(mysqlBackupDir)
	if err != nil {
		return "", fmt.Errorf("read mysql backup [%s] failed, err: %s", mysqlBackupDir, err.Error())
	}
	mysqlBackupFileName := fmt.Sprintf("%s/%d.sql", mysqlBackupDir, len(fileList)+1)
	oper := fmt.Sprintf("mysqldump -u%s -p%s -h%s -P%d > %s",
		mysqlBackup.Spec.MysqlSource.Username,
		mysqlBackup.Spec.MysqlSource.Password,
		mysqlBackup.Spec.MysqlSource.Host,
		mysqlBackup.Spec.MysqlSource.Port,
		mysqlBackupFileName)
	cmd := exec.Command("bash", "-c", oper)
	_, err = cmd.Output()
	if err != nil {
		return "", fmt.Errorf("execute command [%s] failed, err: %s", oper, err.Error())
	}
	logger.L().Info().Msgf("execute command [%s] successful!", oper)
	return mysqlBackupFileName, nil
}

// 上传Mysql备份文件到MinIO
func (r *MysqlBackupReconciler) UploadMysqlDumpFileToMinIO(mysqlBackupFileName string, mysqlBackup *operatorcodehorsecomv1beta1.MysqlBackup) error {
	minioClient, err := minio.New(mysqlBackup.Spec.BackupDestination.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(mysqlBackup.Spec.BackupDestination.AccessKey,
			mysqlBackup.Spec.BackupDestination.AccessSecret, ""),
		Secure: false,
	})
	if err != nil {
		return fmt.Errorf("create minio client failed, err: %s", err.Error())
	}
	object, err := os.Open(mysqlBackupFileName)
	if err != nil {
		return fmt.Errorf("open file [%s] failed, err: %s", mysqlBackupFileName, err.Error())
	}
	ctx := context.TODO()
	_, err = minioClient.PutObject(ctx,
		mysqlBackup.Spec.BackupDestination.BucketName,
		mysqlBackupFileName,
		object,
		-1,
		minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("put mysql backup file [%s] failed, err: %s", mysqlBackupFileName, err.Error())
	}
	return nil
}
