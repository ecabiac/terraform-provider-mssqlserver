package provider

import (
	_ "context"
	_ "database/sql"
	_ "errors"
	_ "fmt"

	_ "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDatabaseBackup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDatabaseBackupRead,

		Schema: map[string]*schema.Schema{
			"filename": {
				Type:     schema.TypeString,
				Required: true,
			},
			"originaldb": {
				Type:     schema.TypeString,
				Required: true,
			},
			"datafile": {
				Type:     schema.TypeString,
				Required: true,
			},
			"logfile": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceDatabaseBackupRead(data *schema.ResourceData, meta interface{}) error {
	return nil
}
