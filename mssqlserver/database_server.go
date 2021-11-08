package mssqlserver

import (
	"database/sql"
	"fmt"
)

type SelfExistsFunc func() (bool, error)
type ChildExistsFunc func(interface{}) (bool, error)

type MsSqlServerManager struct {
	Db *sql.DB
}

type DatabaseSchemaRecord struct {
	Name string
}

type PrincipalRecord struct {
	PrincipalId int
}

type UserRecord struct {
	Name             string
	DatabaseName     string
	PrincipalId      int
	OwnerPrincipalId int
	Sid              string
}

type countResult struct {
	Count int
}

// Represents a simple database backup file which contains a
// sing data (mdf) file and a single log (ldf) file
type DatabaseBackupFileInfo struct {
	Path         string
	DataFileName string
	LogFileName  string
}

type DatabaseRestoreInfo struct {
	DataFileDir string
	LogFileDir  string
}

type dbLoginQueryArgs struct {
	Name string
}

func (dbServer *MsSqlServerManager) CheckDatabaseX(name string) (*DatabaseSchemaRecord, error) {
	var row DatabaseSchemaRecord

	err := dbServer.Db.QueryRow(fmt.Sprintf("SELECT name FROM sys.databases where name = '%s'", name)).Scan(&row.Name)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (dbServer *MsSqlServerManager) DatabaseExists(name string) (bool, error) {

	var result countResult

	err := dbServer.Db.QueryRow("SELECT count(1) FROM sys.databases where name = @p1", name).Scan(&result.Count)
	if err != nil {
		return false, err
	}

	return (result.Count > 0), nil
}

func (dbServer *MsSqlServerManager) CheckUserX(database string, username string) (*UserRecord, error) {
	var row UserRecord

	err := dbServer.Db.QueryRow(fmt.Sprintf("SELECT name FROM %s.sys.database_principals where name = '%s'", database, username)).Scan(&row.Name)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func NewMsSqlServerManager(db *sql.DB) *MsSqlServerManager {
	return &MsSqlServerManager{
		Db: db,
	}
}

func (dbServer *MsSqlServerManager) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return dbServer.Db.Query(query, args...)
}

func (dbServer *MsSqlServerManager) GetDatabaseManager(dbName string) *DatabaseManager {

	return &DatabaseManager{
		Db:     dbServer.Db,
		Name:   dbName,
		exists: true,
	}
}

func (dbServer *MsSqlServerManager) restoreDatabase(backupFile *DatabaseBackupFileInfo, restoreInfo *DatabaseRestoreInfo, name string) error {

	_, err := dbServer.CheckDatabaseX(name)
	// only try to create database if it not exists

	if err == sql.ErrNoRows {

		restoreDataFilePath := fmt.Sprintf("%s/%s.mdf", restoreInfo.DataFileDir, name)
		restoreLogFilePath := fmt.Sprintf("%s/%s.ldf", restoreInfo.LogFileDir, name)
		createCommand := fmt.Sprintf("RESTORE DATABASE %s FROM DISK = N'%s' WITH  FILE = 2, MOVE N'%s' TO N'%s', MOVE N'%s' TO N'%s', NOUNLOAD, STATS = 5",
			name,
			backupFile.Path,
			backupFile.DataFileName,
			restoreDataFilePath,
			backupFile.LogFileName,
			restoreLogFilePath,
		)

		_, err := dbServer.Db.Exec(createCommand)
		return err
	}

	return nil
}
