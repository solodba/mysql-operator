package v1beta1

// MysqlBackup结构体初始化
func NewMysqlBackup() *MysqlBackup {
	return &MysqlBackup{
		Spec:   NewMysqlBackupSpec(),
		Status: NewMysqlBackupStatus(),
	}
}

// MysqlBackupSpec结构体初始化
func NewMysqlBackupSpec() *MysqlBackupSpec {
	return &MysqlBackupSpec{
		MysqlSource:       NewMysqlSource(),
		BackupDestination: NewBackupDestination(),
	}
}

// MysqlSource结构体初始化
func NewMysqlSource() *MysqlSource {
	return &MysqlSource{}
}

// BackupDestination结构体初始化
func NewBackupDestination() *BackupDestination {
	return &BackupDestination{}
}

// MysqlBackupStatus结构体初始化
func NewMysqlBackupStatus() *MysqlBackupStatus {
	return &MysqlBackupStatus{}
}

// MysqlBackupList结构体初始化
func NewMysqlBackupList() *MysqlBackupList {
	return &MysqlBackupList{
		Items: make([]MysqlBackup, 0),
	}
}
