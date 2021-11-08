package provider

import (
	_ "context"
	_ "database/sql"
	_ "errors"
	_ "fmt"

	_ "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDatabaseServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDatabaseBackupRead,

		Schema: map[string]*schema.Schema{
			"data_file_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_file_path": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceDatabaseServerRead(data *schema.ResourceData, meta interface{}) error {
	return nil
}
