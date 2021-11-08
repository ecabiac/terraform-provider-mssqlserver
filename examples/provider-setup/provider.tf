terraform {
  required_providers {
    mssqlserver = {
      source  = "github.com/ecabiac/mssqlserver"
      version = ">= 0.1.0"
    }
  }
}

// Example SQL Server instance running in docker 
// but listening on the host's port 1433
provider "mssqlserver" {
  username  = "sa"
  password  = "yourStrong(!)Password"
  host  = "host.docker.internal"
  port = 1433
}
