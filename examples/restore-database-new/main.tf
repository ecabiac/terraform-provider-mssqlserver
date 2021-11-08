
// Generate a pet name for our new Database instance
resource "random_pet" "foo_dbname" {
  prefix = "foo"
  separator = "_"
}

// restore the database from a backup file that
// already exists on the DB Server machine
resource "mssqlserver_database" "foo" {  
  name            = random_pet.foo_dbname.id
  drop_on_destroy = true

  backup_restore {
      filename   = "/var/opt/mssql/data/foo_backup.bak"
      originaldb = "<not needed>" // todo: remove this
      datafile   = "foo"
      logfile    = "foo_log"
  }
}

// Create a new password for a new SQL Server Login
resource "random_password" "foo_password" {
  length           = 16
  special          = true
  override_special = "_&"
}

// Create a new SQL Login to use for this database
resource "mssqlserver_login" "foo_login" {
  provider     = mssqlserver
  name         = "${random_pet.foo_dbname.id}_login"
  password     = random_password.foo_password.result
}

// The backup database includes a user "foouser"
// that needs to be mapped to the login we created
resource "mssqlserver_user_login" "foo_user_login" {
  database = mssqlserver_database.foo.name
  username = "foouser"
  login = mssqlserver_login.foo_login.name
}

output "foo_dbname" {
  value = mssqlserver_database.foo.name
}

output "foo_login" {
  value = mssqlserver_login.foo_login.name
}

output "foo_password" {
  value = random_password.foo_password.result
}