data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./internal/models", 
    "--dialect", "mysql", 
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
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