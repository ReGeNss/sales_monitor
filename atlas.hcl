data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    # 1. ВИПРАВЛЕНО: "iternal" -> "internal" (typo)
    "--path", "./internal/models", 
    "--dialect", "mysql", 
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  
  # 2. ВИПРАВЛЕНО: Додано версію MariaDB. 
  # Atlas потребує формату docker://image/version/database_name
  dev = "docker://mariadb/latest/dev" 

  migration {
    dir = "file://migrations"
  }
  
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}