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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the Login. Must be unique within the datbase server instance.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The password to use for the Login",
				Sensitive:   true,
			},
			"default_database": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "master",
				Description: "The default database to assign to a Login",
			},
			"drop_on_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "A flag indicatintg that destroying the resource should drop the Login.",
			},
			"sid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The sid assigned to the Login",
			},
			"principal_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The principal id assigned to the Login",
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
	defaultdb := d.Get("default_database").(string)

	exists, err := dbServer.ServerLoginExists(name)
	if err != nil {
		return diag.FromErr(err)
	}

	if exists == false {
		loginCreate := &mssqlserver.ServerLoginCreate{
			Name:            name,
			Password:        password,
			DefaultDatabase: defaultdb,
		}

		_, err := dbServer.CreateLogin(loginCreate)

		if err != nil {
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
	dropOnDestroy := d.Get("drop_on_destroy").(bool)

	if dropOnDestroy {
		dbServer := m.(*mssqlserver.MsSqlServerManager)
		loginName := d.Id()
		_ = dbServer.KillLogins(loginName)

		err := dbServer.DropLogin(loginName)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}
