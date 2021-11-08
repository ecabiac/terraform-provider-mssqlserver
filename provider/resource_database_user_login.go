package provider

import (
	"context"
	_ "errors"
	"fmt"
	"strings"

	"github.com/ecabiac/terraform-provider-mssqlserver/mssqlserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatabaseUserLogin() *schema.Resource {
	return &schema.Resource{
		Description: "Defines a Link between an instance level Login and a database level User",

		CreateContext: databaseUserLoginResourceCreate,
		ReadContext:   databaseUserLoginResourceRead,
		UpdateContext: databaseUserLoginResourceUpdate,
		DeleteContext: databaseUserLoginResourceDelete,

		Schema: map[string]*schema.Schema{

			"database": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"login": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// Called by Terraform to create a new mssqlserver_login resource
func databaseUserLoginResourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	dbName := d.Get("database").(string)
	userName := d.Get("username").(string)
	loginName := d.Get("login").(string)

	dbServer := m.(*mssqlserver.MsSqlServerManager)
	dbManager := dbServer.GetDatabaseManager(dbName)

	dbExists, err := dbManager.DbExists()
	if err != nil {
		return diag.FromErr(err)
	}

	if dbExists == false {
		return diag.FromErr(fmt.Errorf("Database doesn't exist"))
	}

	loginExists, err := dbServer.ServerLoginExists(loginName)
	if err != nil {
		return diag.FromErr(err)
	}

	if loginExists == false {
		return diag.FromErr(fmt.Errorf("Login doesn't exist"))
	}

	userExists, err := dbManager.UserExists(userName)
	if err != nil {
		return diag.FromErr(err)
	}

	if userExists == false {
		return diag.FromErr(fmt.Errorf("User doesn't exist"))
	}

	err = dbManager.AttachUser(userName, loginName)

	if err != nil {
		return diag.FromErr(err)
	}

	id := fmt.Sprintf("%s.%s.%s", dbName, userName, loginName)
	d.SetId(id)

	return nil
}

func databaseUserLoginResourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dbServer := m.(*mssqlserver.MsSqlServerManager)

	//dbName := d.Get("database").(string)
	//dbManager := dbServer.GetDatabaseManager(dbName)

	//exists, err := dbManager.DbExists()

	id := d.Id()
	v := strings.Split(id, ".")
	dbName := v[0]
	userName := v[1]

	dbManager := dbServer.GetDatabaseManager(dbName)
	userRecord, err := dbManager.GetUser(userName)

	if err != nil {
		return diag.FromErr(err)
	}

	if userRecord == nil {
		return diag.FromErr(fmt.Errorf("Couldn't resolve user"))
	}

	//loginRecord, err := dbServer.GetLoginByName(id)

	//if err != nil {
	//	err2 := fmt.Errorf("Error reading login %s\n%w", id, err)
	//	return diag.FromErr(err2)
	//}

	//d.SetId(id)
	d.Set("database", dbName)
	d.Set("username", userRecord.Name)
	d.Set("login", userRecord.Login)

	return nil
}

func databaseUserLoginResourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func databaseUserLoginResourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//dbServer := m.(*mssqlserver.MsSqlServerManager)
	//
	//loginName := d.Id()
	//err := dbServer.DropLogin(loginName)
	//
	//if err != nil {
	//	return diag.FromErr(err)
	//}

	return nil
}
