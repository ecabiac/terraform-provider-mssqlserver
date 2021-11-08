package provider

import (
	"database/sql"
	"fmt"
)

func checkDatabase(db *sql.DB, name string) (*DatabaseSchemaRow, error) {
	var row DatabaseSchemaRow
	err := db.QueryRow(fmt.Sprintf("SELECT name FROM sys.databases where name = '%s'", name)).Scan(&row.name)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func checkLogin(db *sql.DB, username string) (*PrinicipalsRow, error) {

	var row PrinicipalsRow
	err := db.QueryRow(fmt.Sprintf("SELECT principal_id FROM master.sys.server_principals where name = '%s'", username)).Scan(&row.principal_id)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func checkUser(db *sql.DB, database string, username string) (*UsersRow, error) {
	var row UsersRow

	err := db.QueryRow(fmt.Sprintf("SELECT name FROM %s.sys.database_principals where name = '%s'", database, username)).Scan(&row.name)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

type DatabaseSchemaRow struct {
	name string
}

type PrinicipalsRow struct {
	principal_id int
}

type UsersRow struct {
	name string
}
