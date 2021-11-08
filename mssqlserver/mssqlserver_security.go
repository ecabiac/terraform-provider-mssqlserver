package mssqlserver

import (
	"database/sql"
	"encoding/base64"
	"fmt"
)

type ServerLoginRecord struct {
	PrincipalId int
	Name        string
	Sid         string
}

type ServerLoginCreate struct {
	Name     string
	Password string
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

	cmd := fmt.Sprintf("CREATE LOGIN [%s] WITH PASSWORD = '%s', CHECK_POLICY = OFF, CHECK_EXPIRATION = OFF", login.Name, login.Password)
	_, err := dbServer.Db.Exec(cmd)

	if err != nil {
		err2 := fmt.Errorf("Failed to create login %s\nOriginal query\n%s\n%w", login.Name, cmd, err)
		return nil, err2
	}

	return dbServer.GetLoginByName(login.Name)
}

func (dbServer *MsSqlServerManager) DropLogin(loginName string) error {

	cmd := fmt.Sprintf("DROP LOGIN [%s]", loginName)
	_, err := dbServer.Db.Exec(cmd)

	return err
}
