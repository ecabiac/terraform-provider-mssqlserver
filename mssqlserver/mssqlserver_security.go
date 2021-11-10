package mssqlserver

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"time"
)

type ServerLoginRecord struct {
	PrincipalId int
	Name        string
	Sid         string
}

type ServerLoginCreate struct {
	Name            string
	Password        string
	DefaultDatabase string
}

type whoRecord struct {
	Spid      int
	Ecid      int
	Status    string
	LoginName string
	HostName  string
	Blk       string
	DbName    sql.NullString
	Cmd       sql.NullString
	RequestId int
}

// GetLoginByName retrieves a Sql Server instance login principal
//
// If the login does not exist, nil is returned
func (dbServer *MsSqlServerManager) GetLoginByName(loginName string) (*ServerLoginRecord, error) {
	var record ServerLoginRecord
	var sidVal *sql.RawBytes
	row := dbServer.Db.QueryRow("SELECT name, principal_id, sid FROM master.sys.server_principals where name = @p1", loginName)
	err := row.Scan(&record.Name, &record.PrincipalId, &sidVal)

	if err != nil {
		return nil, err
	}

	record.Sid = base64.StdEncoding.EncodeToString(*sidVal)
	return &record, nil
}

func (dbServer *MsSqlServerManager) ServerLoginExists(username string) (bool, error) {

	var result countResult

	query := "SELECT count(1) FROM master.sys.server_principals where name = @p1"
	err := dbServer.Db.QueryRow(query, username).Scan(&result.Count)
	if err != nil {
		return false, err
	}

	return (result.Count > 0), nil
}

// Creates a SqlServer Login with the provided password
//
// This function does not guard against SQL Injection so consumers
// should not allow arbitrary values for either the login name
// or the password
func (dbServer *MsSqlServerManager) CreateLogin(login *ServerLoginCreate) (*ServerLoginRecord, error) {

	cmd := fmt.Sprintf("CREATE LOGIN [%s] WITH PASSWORD = '%s', DEFAULT_DATABASE = [%s], CHECK_POLICY = OFF, CHECK_EXPIRATION = OFF", login.Name, login.Password, login.DefaultDatabase)
	_, err := dbServer.Db.Exec(cmd)

	if err != nil {
		err2 := fmt.Errorf("Failed to create login %s\nOriginal query\n%s\n%w", login.Name, cmd, err)
		return nil, err2
	}

	return dbServer.GetLoginByName(login.Name)
}

func (dbServer *MsSqlServerManager) DropLogin(loginName string) error {
	cmd := fmt.Sprintf("DROP LOGIN [%s]", loginName)
	retryFunc := func() error {
		_, err := dbServer.Db.Exec(cmd)
		return err
	}

	err := retry(3, time.Duration(time.Second), retryFunc)

	return err
}

func (dbServer *MsSqlServerManager) KillLogins(loginName string) error {

	var resultRow whoRecord

	// https://docs.microsoft.com/en-us/sql/t-sql/language-elements/kill-transact-sql?view=sql-server-ver15
	// Requires the ALTER ANY CONNECTION permission. ALTER ANY CONNECTION is
	// included with membership in the sysadmin or processadmin fixed server roles.
	cmd := fmt.Sprintf("EXEC sp_who '%s'", loginName)
	rows, err := dbServer.Db.Query(cmd)
	for more := rows.Next(); more; more = rows.Next() {
		err = rows.Scan(&resultRow.Spid, &resultRow.Ecid, &resultRow.Status, &resultRow.LoginName, &resultRow.HostName, &resultRow.Blk, &resultRow.DbName, &resultRow.Cmd, &resultRow.RequestId)

		if err == nil {
			cmd = fmt.Sprintf("Kill %d;", resultRow.Spid)
			dbServer.Db.Exec(cmd)
		} else {
			log.Printf(err.Error())
		}

	}

	return err
}
