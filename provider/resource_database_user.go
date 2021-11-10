package provider

import (
	"context"
	"fmt"

	"github.com/ecabiac/terraform-provider-mssqlserver/mssqlserver"
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
				Computed: true,
			},
		},
	}
}

func resourceDatabaseUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	dbServer := m.(*mssqlserver.MsSqlServerManager)

	databaseName := d.Get("database").(string)
	username := d.Get("username").(string)
	defaultSchema := d.Get("default_schema").(string)

	dbManager := dbServer.GetDatabaseManager(databaseName)
	userExists, err := dbManager.UserExists(username)

	if err != nil {
		return diag.FromErr(err)
	}

	if userExists == true {
		return diag.FromErr(fmt.Errorf("User already exists"))
	}

	userCreate := &mssqlserver.DatabaseUserCreate{
		Name:          username,
		DefaultSchema: defaultSchema,
	}

	_, err = dbManager.CreateUser(userCreate)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s.%s", databaseName, username))

	return nil
}

func resourceDatabaseUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	dbServer := m.(*mssqlserver.MsSqlServerManager)

	databaseName := d.Get("database").(string)
	username := d.Get("username").(string)
	database := dbServer.GetDatabaseManager(databaseName)
	dbUser, err := database.GetUser(username)

	if err != nil {
		return diag.FromErr(err)
	}

	if dbUser == nil {
		return nil
	}

	d.Set("username", dbUser.Name)
	d.Set("default_schema", dbUser.DefaultSchema)
	d.Set("login", dbUser.Login)

	return nil
}

func resourceDatabaseUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceDatabaseUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dbServer := m.(*mssqlserver.MsSqlServerManager)

	databaseName := d.Get("database").(string)
	username := d.Get("username").(string)

	database := dbServer.GetDatabaseManager(databaseName)

	dbExists, err := database.DbExists()

	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to determine Database status: %w", err))
	}

	if dbExists == false {
		return nil
	}

	userExists, err := database.UserExists(username)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to determine user status: %w", err))
	}

	if userExists == false {
		return nil
	}

	err = database.DropUser(username)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
