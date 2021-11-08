package provider

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/ecabiac/terraform-provider-mssqlserver/mssqlserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The username used to login to the SQL Server Instance",
			},

			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The password for the user specified by username",
			},

			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The addressable hostname of the SQL Server Instance",
			},

			"port": {
				Type:        schema.TypeInt,
				Default:     1433,
				Optional:    true,
				Description: "The port of the SQL Server Instance",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"mssqlserver_database":   resourceDatabase(),
			"mssqlserver_user":       resourceDatabaseUser(),
			"mssqlserver_login":      resourceServerLogin(),
			"mssqlserver_user_login": resourceDatabaseUserLogin(),
		},

		DataSourcesMap: map[string]*schema.Resource{},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d",
		d.Get("host"), username, password, d.Get("port"))

	log.Printf(" connString:%s\n", connString)

	db, err := sql.Open("sqlserver", connString)

	if err != nil {
		err2 := fmt.Errorf("Open connection failed: %w", err)
		return nil, diag.FromErr(err2)
	}

	dbs := mssqlserver.NewMsSqlServerManager(db)

	return dbs, nil
}
