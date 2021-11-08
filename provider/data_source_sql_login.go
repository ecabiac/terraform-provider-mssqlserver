package provider

import (
	_ "context"
	"database/sql"
	_ "database/sql"
	_ "errors"
	_ "fmt"

	//"mssqlserver"
	"github.com/ecabiac/terraform-provider-mssqlserver/mssqlserver"

	_ "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSqlLogin() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDatabaseBackupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSqlLoginRead(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)

	dbs := mssqlserver.NewMsSqlServerManager(db)
	login, err := dbs.GetLoginByName(name)

	if err != nil {
		data.Set("id", login.PrincipalId)
		data.Set("sid", login.Sid)
	}

	return err
}
