package environment

type Environment struct {
	Production bool     `yaml:"production"`
	Database   Database `yaml:"database"`
}

type Database struct {
	Driver                 string `yaml:"driver"`
	Host                   string `yaml:"host"`
	Port                   int    `yaml:"port"`
	Username               string `yaml:"username"`
	Password               string `yaml:"password"`
	Database               string `yaml:"database"`
	MaxDatabaseConnections int    `yaml:"max_database_connections"`
	MaxIdleConnections     int    `yaml:"max_idle_connections"`
}
