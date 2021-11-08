package provider

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatabaseBackup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseCreate,
		ReadContext:   resourceDatabaseRead,
		UpdateContext: resourceDatabaseUpdate,
		DeleteContext: databaseResourceDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"drop_on_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceDatabaseBackupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	db := m.(*sql.DB)
	name := d.Get("name").(string)

	row, err := checkDatabase(db, name)

	// only try to create database if it not exists
	if err == sql.ErrNoRows {
		_, err := db.Query(fmt.Sprintf("CREATE DATABASE %s", name))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	row, err = checkDatabase(db, name)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(row.name)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDatabaseBackupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	db := m.(*sql.DB)
	name := d.Id()
	row, err := checkDatabase(db, name)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", row.name); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDatabaseBackupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceDatabaseBackupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dropOnDestroy := d.Get("drop_on_destroy").(bool)

	if dropOnDestroy {
		name := d.Id()
		db := m.(*sql.DB)
		row, err := checkDatabase(db, name)
		if err != nil && err != sql.ErrNoRows {
			return diag.FromErr(errors.New(fmt.Sprint("Failed to check if database exists", err)))
		}

		if row != nil {
			_, err = db.Query(fmt.Sprintf("exec('USE master; ALTER DATABASE %s SET SINGLE_USER WITH ROLLBACK IMMEDIATE')", name))
			if err != nil {
				return diag.FromErr(errors.New(fmt.Sprint("Failed to set database to single user mode for dropping database: ", err)))
			}

			_, err := db.Exec(fmt.Sprintf("exec('USE master; DROP DATABASE %s')", name))
			if err != nil {
				return diag.FromErr(errors.New(fmt.Sprint("Failed to drop database: ", err)))
			}
		}
	}

	return nil
}
