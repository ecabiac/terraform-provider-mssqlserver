package mssqlserver

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"
)

type DatabaseManager struct {
	Db     *sql.DB
	Name   string
	exists bool
}

type DatabaseUser struct {
	Name  string
	Login string
}

type databaseCountRecord struct {
	DbCount int
}

type databaseNameRecord struct {
	Name string
}

func (dbManager *DatabaseManager) Drop() error {

	_, err := dbManager.Db.Exec(fmt.Sprintf("USE master; ALTER DATABASE %s SET SINGLE_USER WITH ROLLBACK IMMEDIATE;", dbManager.Name))
	if err != nil {
		return fmt.Errorf("Failed to set database %s to single user mode for dropping database\n:%w", dbManager.Name, err)
	}

	_, err = dbManager.Db.Exec(fmt.Sprintf("exec('USE master; DROP DATABASE %s')", dbManager.Name))
	if err != nil {
		time.Sleep(2 * time.Second)
		_, err = dbManager.Db.Exec(fmt.Sprintf("exec('USE master; DROP DATABASE %s')", dbManager.Name))
		if err != nil {
			return fmt.Errorf("Failed to drop database %s:\n%w", dbManager.Name, err)
		}
	}

	return nil
}

func (dbManager *DatabaseManager) DbExists() (bool, error) {
	var result countResult

	err := dbManager.Db.QueryRow("SELECT count(1) FROM sys.databases where name = @p1", dbManager.Name).Scan(&result.Count)
	if err != nil {
		return false, err
	}

	return (result.Count > 0), nil
}

func (dbManager *DatabaseManager) Restore(backupFile *DatabaseBackupFileInfo, restoreInfo *DatabaseRestoreInfo) error {

	exists, err := dbManager.DbExists()

	// only try to create database if it does not exist
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("Database %s already exists", dbManager.Name)
	}

	dataFileName := fmt.Sprintf("%s.mdf", dbManager.Name)
	logFileName := fmt.Sprintf("%s.ldf", dbManager.Name)
	restoreDataFilePath := filepath.Join(restoreInfo.DataFileDir, dataFileName)
	//fmt.Sprintf("%s/%s.mdf", restoreInfo.DataFileDir, dbManager.Name)
	restoreLogFilePath := filepath.Join(restoreInfo.LogFileDir, logFileName)
	//fmt.Sprintf("%s/%s.ldf", restoreInfo.LogFileDir, dbManager.Name)
	createCommand := fmt.Sprintf("RESTORE DATABASE %s FROM DISK = N'%s' WITH  FILE = 2, MOVE N'%s' TO N'%s', MOVE N'%s' TO N'%s', NOUNLOAD, STATS = 5",
		dbManager.Name,
		backupFile.Path,
		backupFile.DataFileName,
		restoreDataFilePath,
		backupFile.LogFileName,
		restoreLogFilePath,
	)

	_, err = dbManager.Db.Exec(createCommand)
	return err
}

func (dbManager *DatabaseManager) Create() error {

	exists, err := dbManager.DbExists()

	// only try to create database if it does not exist
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("Database %s already exists", dbManager.Name)
	}

	createCommand := fmt.Sprintf("CREATE DATABASE %s", dbManager.Name)
	_, err = dbManager.Db.Exec(createCommand)

	return err
}

func (dbManager *DatabaseManager) UserExists(username string) (bool, error) {
	var result countResult

	query := fmt.Sprintf("SELECT count(1) FROM [%s].[sys].[database_principals] where type = 'S' AND name = @p1", dbManager.Name)
	err := dbManager.Db.QueryRow(query, username).Scan(&result.Count)
	if err != nil {
		return false, err
	}

	return (result.Count > 0), nil
}

func (dbManager *DatabaseManager) GetUser(username string) (*DatabaseUser, error) {
	var result DatabaseUser

	queryString := `
    select dp.name as UserName, sp.name as LoginName
    from [%s].[sys].[database_principals] dp
    LEFT OUTER JOIN [%s].[sys].[server_principals] sp on dp.sid = sp.sid
    where dp.type = 'S' AND dp.name = @p1;
    `
	query := fmt.Sprintf(queryString, dbManager.Name, dbManager.Name)
	err := dbManager.Db.QueryRow(query, username).Scan(&result.Name, &result.Login)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (dbManager *DatabaseManager) AttachUser(userName string, loginName string) error {

	//exists, err := dbManager.DbExists()

	//// only try to create database if it does not exist
	//if err != nil {
	//	return err
	//}

	//if exists {
	//	return fmt.Errorf("Database %s already exists", dbManager.Name)
	//}

	alterCommand := fmt.Sprintf("USE [%s];ALTER USER [%s] WITH LOGIN = [%s];", dbManager.Name, userName, loginName)
	_, err := dbManager.Db.Exec(alterCommand)

	return err
}

//func (dbManager *DatabaseManager) DetachUser(userName string) error {
//
//	//exists, err := dbManager.DbExists()
//
//	//// only try to create database if it does not exist
//	//if err != nil {
//	//	return err
//	//}
//
//	//if exists {
//	//	return fmt.Errorf("Database %s already exists", dbManager.Name)
//	//}
//
//	alterCommand := fmt.Sprintf("USE [%s];ALTER USER [%s] WITH LOGIN = [%s];", dbManager.Name, userName, loginName)
//	_, err := dbManager.Db.Exec(alterCommand)
//
//	return err
//}
