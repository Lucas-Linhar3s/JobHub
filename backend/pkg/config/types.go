package config

// Config is a struct that contains the configuration of the application
type Config struct {
	Env      string    `yaml:"env,required"`
	Http     *Http     `yaml:"http,required"`
	Security *Security `yaml:"security,required"`
	Data     *Data     `yaml:"data,required"`
	Log      *Log      `yaml:"log,required"`
}

// Http is a struct that contains the host and port of the http server
type Http struct {
	Host string `yaml:"http.host,required"`
	Port string `yaml:"http.port,required"`
}

// Security is a struct that contains the security configuration of the application
type Security struct {
	ApiSign *ApiSign `yaml:"apiSign,required"`
	Jwt     *Jwt     `yaml:"jwt,required"`
	Oauth2  *Oauth2  `yaml:"oauth2,required"`
}

type Oauth2 struct {
	Google *Google `yaml:"google,required"`
	Github *Github `yaml:"github,required"`
}

// ApiSign is a struct that contains the app key and app security
type ApiSign struct {
	AppKey      string `yaml:"security.api_sign.app_key,required,environment"`
	AppSecurity string `yaml:"security.api_sign.app_security,required,environment"`
}

// Jwt is a struct that contains the key of the jwt
type Jwt struct {
	ExpiresAt int    `yaml:"security.jwt.expire_at,required"`
	Key       string `yaml:"security.jwt.key,required,environment"`
}

type Google struct {
	ClientId     string   `yaml:"security.oauth2.google.client_id,required,environment"`
	ClientSecret string   `yaml:"security.oauth2.google.client_secret,required,environment"`
	RedirectUrl  string   `yaml:"security.oauth2.google.redirect_url,required"`
	Scopes       []string `yaml:"security.oauth2.google.scopes,required"`
}

type Github struct {
	ClientId     string   `yaml:"security.oauth2.github.client_id,required,environment"`
	ClientSecret string   `yaml:"security.oauth2.github.client_secret,required,environment"`
	RedirectUrl  string   `yaml:"security.oauth2.github.redirect_url,required"`
	Scopes       []string `yaml:"security.oauth2.github.scopes,required"`
}

// Data is a struct that contains the database configuration
type Data struct {
	DB *Db `yaml:"db,required"`
}

// Db is a struct that contains the user configuration of the database
type Db struct {
	User *User `yaml:"user,required"`
}

// User is a struct that contains the user configuration of the database
type User struct {
	Driver            string `yaml:"data.db.user.driver,required"`
	Nick              string `yaml:"data.db.user.nick,required"`
	Name              string `yaml:"data.db.user.name,required"`
	Username          string `yaml:"data.db.user.username,required"`
	Password          string `yaml:"data.db.user.password,required"`
	HostName          string `yaml:"data.db.user.hostname,required"`
	Port              string `yaml:"data.db.user.port,required"`
	MaxConn           int    `yaml:"data.db.user.max_conn,required"`
	MaxIdle           int    `yaml:"data.db.user.max_idle,required"`
	TransationTimeout int    `yaml:"data.db.user.transaction_timeout,required"`
	Dsn               string `yaml:"data.db.user.dsn"`
}

// Log is a struct that contains the log configuration
type Log struct {
	LogLevel    string `yaml:"log.log_level,required"`
	Enconding   string `yaml:"log.encoding,required"`
	LogFileName string `yaml:"log.log_file_name,required"`
	MaxBackups  int    `yaml:"log.max_backups,required"`
	MaxAge      int    `yaml:"log.max_age,required"`
	MaxSize     int    `yaml:"log.max_size,required"`
	Compress    bool   `yaml:"log.compress,required"`
}
