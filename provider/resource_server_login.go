package provider

import (
	"context"
	"fmt"

	"github.com/ecabiac/terraform-provider-mssqlserver/mssqlserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServerLogin() *schema.Resource {
	return &schema.Resource{
		CreateContext: serverLoginResourceCreate,
		ReadContext:   serverLoginResourceRead,
		UpdateContext: serverLoginResourceUpdate,
		DeleteContext: serverLoginResourceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"principal_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"roles": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

// Called by Terraform to create a new mssqlserver_login resource
func serverLoginResourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dbServer := m.(*mssqlserver.MsSqlServerManager)

	name := d.Get("name").(string)
	password := d.Get("password").(string)

	exists, err := dbServer.ServerLoginExists(name)
	if err != nil {
		return diag.FromErr(err)
	}

	if exists == false {
		loginCreate := &mssqlserver.ServerLoginCreate{
			Name:     name,
			Password: password,
		}

		_, err := dbServer.CreateLogin(loginCreate)
		//_, err := db.Query(fmt.Sprintf("CREATE LOGIN \"%s\" WITH PASSWORD = '%s', CHECK_POLICY = OFF, CHECK_EXPIRATION = OFF", name, password))
		if err != nil {
			//return diag.FromErr(errors.New(fmt.Sprint("Failed to create login", err)))
			return diag.FromErr(err)
		}
	}

	loginRecord, err := dbServer.GetLoginByName(name)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)
	d.Set("sid", loginRecord.Sid)
	d.Set("principal_id", loginRecord.PrincipalId)

	return nil
}

func serverLoginResourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dbServer := m.(*mssqlserver.MsSqlServerManager)

	id := d.Id()

	loginRecord, err := dbServer.GetLoginByName(id)

	if err != nil {
		err2 := fmt.Errorf("Error reading login %s\n%w", id, err)
		return diag.FromErr(err2)
	}

	d.SetId(id)
	d.Set("name", loginRecord.Name)
	d.Set("sid", loginRecord.Sid)
	d.Set("principal_id", loginRecord.PrincipalId)

	return nil
}

func serverLoginResourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func serverLoginResourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dbServer := m.(*mssqlserver.MsSqlServerManager)

	loginName := d.Id()
	err := dbServer.DropLogin(loginName)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
