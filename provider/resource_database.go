package provider

import (
	"context"
	"database/sql"
	"log"
	"strconv"

	//"mssqlserver"
	"github.com/ecabiac/terraform-provider-mssqlserver/mssqlserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DatabaseBackup struct {
	FileName   string
	OriginalDb string
	DataFile   string
	LogFile    string
}

func resourceDatabase() *schema.Resource {
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
			"backup_restore": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filename": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The full path to the .bak file",
						},
						"originaldb": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"datafile": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"logfile": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"user_logins": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"login": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceDatabaseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dbServer := m.(*mssqlserver.MsSqlServerManager)
	name := d.Get("name").(string)

	dbManager := dbServer.GetDatabaseManager(name)
	exists, err := dbManager.DbExists()

	if err != nil {
		return diag.FromErr(err)
	}

	if exists {

		return nil
	}

	backupData, backupOk := d.GetOk("backup_restore")
	if backupOk {
		dbData := backupData.(*schema.Set).List()
		dbDataItem := dbData[0].(map[string]interface{})
		//dbBackup := &DatabaseBackup{
		//	FileName:   dbDataItem["filename"].(string),
		//	OriginalDb: dbDataItem["originaldb"].(string),
		//	DataFile:   dbDataItem["datafile"].(string),
		//	LogFile:    dbDataItem["logfile"].(string),
		//}

		//restoreDataFilePath := fmt.Sprintf("/var/opt/mssql/data%s.mdf", name)
		//restoreLogFilePath := fmt.Sprintf("/var/opt/mssql/data%s.ldf", name)

		backupInfo := &mssqlserver.DatabaseBackupFileInfo{
			Path:         dbDataItem["filename"].(string),
			DataFileName: dbDataItem["datafile"].(string),
			LogFileName:  dbDataItem["logfile"].(string),
		}

		restoreInfo := &mssqlserver.DatabaseRestoreInfo{
			DataFileDir: "/var/opt/mssql/",
			LogFileDir:  "/var/opt/mssql/",
		}

		err = dbManager.Restore(backupInfo, restoreInfo)

		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(name)
		return nil
		//createCommand = fmt.Sprintf("RESTORE DATABASE %s FROM DISK = N'%s' WITH  FILE = 2, MOVE N'%s' TO N'%s', MOVE N'%s' TO N'%s', NOUNLOAD, STATS = 5",
		//	name,
		//	dbBackup.FileName,
		//	dbBackup.DataFile,
		//	restoreDataFilePath,
		//	dbBackup.LogFile,
		//	restoreLogFilePath,
		//)
	}

	err = dbManager.Create()

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)
	return nil
}

func resourceDatabaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Id()

	dbServer := m.(*mssqlserver.MsSqlServerManager)
	row, err := dbServer.CheckDatabaseX(name)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", row.Name); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDatabaseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func databaseResourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dropOnDestroy := d.Get("drop_on_destroy").(bool)

	log.Printf("[DEBUG] dropOnDestroy was: %s", strconv.FormatBool(dropOnDestroy))

	if dropOnDestroy {

		name := d.Id()
		dbServer := m.(*mssqlserver.MsSqlServerManager)
		mgr := dbServer.GetDatabaseManager(name)

		exists, err := mgr.DbExists()
		if err != nil {
			return diag.FromErr(err)
		}

		if exists {
			err = mgr.Drop()
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}

func getBackupFromData(data *schema.ResourceData) (*DatabaseBackup, error) {
	dbBackup := &DatabaseBackup{
		FileName:   data.Get("filename").(string),
		OriginalDb: data.Get("originaldb").(string),
		DataFile:   data.Get("datafile").(string),
		LogFile:    data.Get("logfile").(string),
	}

	return dbBackup, nil
}

func getBackupFromDataItem(data map[string]interface{}) (*DatabaseBackup, error) {
	dbBackup := &DatabaseBackup{
		FileName:   data["filename"].(string),
		OriginalDb: data["originaldb"].(string),
		DataFile:   data["datafile"].(string),
		LogFile:    data["logfile"].(string),
	}

	return dbBackup, nil
}
