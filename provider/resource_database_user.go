package provider

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/ecabiac/terraform-provider-mssqlserver/mssqlserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatabaseUser() *schema.Resource {

	return &schema.Resource{
		Description: "Defines a database level User",

		CreateContext: resourceDatabaseUserCreate,
		ReadContext:   resourceDatabaseUserRead,
		UpdateContext: resourceDatabaseUserUpdate,
		DeleteContext: resourceDatabaseUserDelete,

		Schema: map[string]*schema.Schema{
			"database": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_schema": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "dbo",
			},
			"login": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDatabaseUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	db := m.(*sql.DB)
	database := d.Get("database").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	roles := d.Get("roles").(*schema.Set).List()

	row, err := checkLogin(db, username)
	if err == sql.ErrNoRows {

		_, err := db.Query(fmt.Sprintf("CREATE LOGIN \"%s\" WITH PASSWORD = '%s', CHECK_POLICY = OFF, CHECK_EXPIRATION = OFF", username, password))
		if err != nil {
			return diag.FromErr(errors.New(fmt.Sprint("Failed to create login", err)))
		}

		// TODO Schema?
		_, err = db.Query(fmt.Sprintf("exec('use %s; CREATE USER \"%s\" FOR LOGIN \"%s\" with default_schema = dbo')", database, username, username))
		//_, err = db.Query(fmt.Sprintf(  "CREATE USER \"%s\" FOR LOGIN \"%s\" with default_schema = dbo", username, username))
		if err != nil {
			return diag.FromErr(errors.New(fmt.Sprint("Failed to create user: ", err)))
		}

	}

	row, err = checkLogin(db, username)
	if err != nil {
		return diag.FromErr(errors.New(fmt.Sprint("Unknow error occured:", err)))
	}

	for _, role := range roles {
		_, err = db.Exec(fmt.Sprintf("exec('use %s; exec(''sp_addrolemember %s, %s '') ')", database, role, username))
		if err != nil {
			return diag.FromErr(errors.New(fmt.Sprint("Failed to add member to role:", err)))
		}
	}

	d.SetId(fmt.Sprint(row.principal_id))

	if err != nil {
		return diag.FromErr(err)
	}
	return nil

}

//SELECT [name]
//FROM [sys].[database_principals]
//WHERE [type] = N'S'

func resourceDatabaseUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	db := m.(*sql.DB)
	row := db.QueryRow(fmt.Sprintf("SELECT name FROM master.sys.server_principals WHERE principal_id = %s", d.Id()))
	var name string
	err := row.Scan(&name)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("username", name); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDatabaseUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceDatabaseUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	db := m.(*sql.DB)
	database := d.Get("database").(string)
	row := db.QueryRow(fmt.Sprintf("SELECT name FROM master.sys.server_principals WHERE principal_id = %s", d.Id()))
	var name string
	err := row.Scan(&name)

	if err != sql.ErrNoRows {
		_, err = db.Query(fmt.Sprintf("DROP LOGIN %s", name))
		if err != nil {
			return diag.FromErr(errors.New(fmt.Sprint("Failed to drop login: ", err)))
		}

		row, _ := checkDatabase(db, database)
		if row != nil {
			row, _ := checkUser(db, database, name)
			if row != nil {
				_, err = db.Query(fmt.Sprintf("exec('use %s; drop user %s');", database, name))
				if err != nil {
					return diag.FromErr(errors.New(fmt.Sprint("Failed to drop user: ", err)))
				}
			}
		}

		// check if user exists
		//SELECT [name]
		//FROM [Beratungsmappe].[sys].[database_principals]
		//WHERE [type] = N'S'

		//		row = db.Query(fmt.Sprintf("SELECT name FROM sys.server_principals WHERE principal_id = %s", d.Id()))

		//_, err = db.Query(fmt.Sprintf("exec('use %s; drop user %s');", database, name))
		//if err != nil {
		//	return errors.New(fmt.Sprint("Failed to drop user: ", err))
		//}

	}

	return nil
}
