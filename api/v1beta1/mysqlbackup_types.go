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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MysqlBackupSpec defines the desired state of MysqlBackup
type MysqlBackupSpec struct {
	// 是否开启备份
	Enable bool `json:"enable"`
	// 备份开始时间
	StartTime string `json:"startTime"`
	// 备份时间间隔
	Period int64 `json:"period"`
	// mysql备份信息
	MysqlSource *MysqlSource `json:"mysqlSource"`
	// 备份目标地址
	BackupDestination *BackupDestination `json:"backupDestination"`
}

// mysql信息
type MysqlSource struct {
	// mysql用户名
	Username string `json:"username"`
	// mysql密码
	Password string `json:"password"`
	// mysql地址
	Host string `json:"host"`
	// mysql端口
	Port int32 `json:"port"`
}

// BackupDestination信息
type BackupDestination struct {
	// 备份目的地地址
	Endpoint string `json:"endpoint"`
	// 访问Key
	AccessKey string `json:"accessKey"`
	// 访问秘钥
	AccessSecret string `json:"accessSecret"`
	// 访问桶名称
	BucketName string `json:"bucketName"`
}

// MysqlBackupStatus defines the observed state of MysqlBackup
type MysqlBackupStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MysqlBackup is the Schema for the mysqlbackups API
type MysqlBackup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MysqlBackupSpec   `json:"spec,omitempty"`
	Status MysqlBackupStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MysqlBackupList contains a list of MysqlBackup
type MysqlBackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MysqlBackup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MysqlBackup{}, &MysqlBackupList{})
}
