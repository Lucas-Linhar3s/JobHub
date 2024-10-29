package config

// Config is a struct that contains the configuration of the application
type Config struct {
	Env      string
	Http     *Http
	Security *Security
	Data     *Data
	Log      *Log
}

// Http is a struct that contains the host and port of the http server
type Http struct {
	Host string
	Port string
}

// Security is a struct that contains the security configuration of the application
type Security struct {
	ApiSign *ApiSign
	Jwt     *Jwt
	Oauth2  *Oauth2
}

type Oauth2 struct {
	Google   *Google
	Github   *Github
	Linkedin *Linkedin
}

// ApiSign is a struct that contains the app key and app security
type ApiSign struct {
	AppKey      string
	AppSecurity string
}

// Jwt is a struct that contains the key of the jwt
type Jwt struct {
	ExpiresAt int
	Key       string
}

type Google struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
	Scopes       []string
}

type Github struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
	Scopes       []string
}

type Linkedin struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
	Scopes       []string
}

// Data is a struct that contains the database configuration
type Data struct {
	DB *Db
}

// Db is a struct that contains the user configuration of the database
type Db struct {
	User *User
}

// User is a struct that contains the user configuration of the database
type User struct {
	Driver            string
	Nick              string
	Name              string
	Username          string
	Password          string
	HostName          string
	Port              string
	MaxConn           int
	MaxIdle           int
	TransationTimeout int
	Dsn               string
}

// Log is a struct that contains the log configuration
type Log struct {
	LogLevel    string
	Enconding   string
	LogFileName string
	MaxBackups  int
	MaxAge      int
	MaxSize     int
	Compress    bool
}
