package models

type Configurations struct {
	Env        string
	LogLevel   string
	LogPath    string
	ServerConf ServerConfig
	RedisConf  RedisConfig
	MySqlConf  MySqlConfig
}

type ServerConfig struct {
	Host         string
	Port         string
	Prefix       string
	ReadTimeout  int
	WriteTimeout int
}

type RedisConfig struct {
	Host                    string `mapstructure:"redis_host"`
	Port                    string `mapstructure:"redis_port"`
	UserName                string `mapstructure:"username"`
	Password                string `mapstructure:"password"`
	RedisDB                 int    `mapstructure:"redis_db"`
	RedisUploadedDataKey    string `mapstructure:"redis_uploaded_data_key"`
	RedisUploadedDataKeyTtl int    `mapstructure:"redis_uploaded_data_key_ttl"`
	UpdateFileStatKeyTtl    int    `mapstructure:"update_file_stat_key_ttl"`
}

type MySqlConfig struct {
	Host     string `mapstructure:"db_host"`
	Port     string `mapstructure:"db_port"`
	UserName string `mapstructure:"db_username"`
	Password string `mapstructure:"db_password"`
	Database string `mapstructure:"db_database"`
}
