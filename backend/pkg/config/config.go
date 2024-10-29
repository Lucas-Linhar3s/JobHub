package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// NewViper is a function that returns a new viper instance
func NewViper() *viper.Viper {
	p := flag.String("conf", "../../config/prod.yml", "config path, eg: -conf ../../config/prod.yml")
	flag.Parse()
	envConf := os.Getenv("APP_CONF")
	if envConf == "" {
		envConf = *p
	}
	fmt.Println("load conf file:", envConf)
	return getViper(envConf)
}

func getViper(dir string) *viper.Viper {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(currentDir, dir)
	conf := viper.New()
	conf.SetConfigFile(path)
	err = conf.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return conf
}

func LoadAttributes(viper *viper.Viper) *Config {
	return &Config{
		Env: viper.GetString("env"),
		Http: &Http{
			Host: viper.GetString("http.host"),
			Port: viper.GetString("http.port"),
		},
		Security: loadAttributesEnv(viper).Security,
		Data: &Data{
			DB: &Db{
				User: &User{
					Driver:            viper.GetString("data.db.user.driver"),
					Nick:              viper.GetString("data.db.user.nick"),
					Name:              viper.GetString("data.db.user.name"),
					Username:          viper.GetString("data.db.user.username"),
					Password:          viper.GetString("data.db.user.password"),
					HostName:          viper.GetString("data.db.user.hostname"),
					Port:              viper.GetString("data.db.user.port"),
					MaxConn:           viper.GetInt("data.db.user.max_conn"),
					MaxIdle:           viper.GetInt("data.db.user.max_idle"),
					TransationTimeout: viper.GetInt("data.db.user.transation_timeout"),
					Dsn:               viper.GetString("data.db.user.dsn"),
				},
			},
		},
		Log: &Log{
			LogLevel:    viper.GetString("log.log_level"),
			Enconding:   viper.GetString("log.enconding"),
			LogFileName: viper.GetString("log.log_file_name"),
			MaxBackups:  viper.GetInt("log.max_backups"),
			MaxAge:      viper.GetInt("log.max_age"),
			MaxSize:     viper.GetInt("log.max_size"),
			Compress:    viper.GetBool("log.compress"),
		},
	}
}

func loadAttributesEnv(viper *viper.Viper) *Config {
	envViper := viper.GetString("env")
	if envViper == "prod" {
		return &Config{
			Security: &Security{
				ApiSign: &ApiSign{
					AppKey:      os.Getenv(viper.GetString("security.api_sign.app_key")),
					AppSecurity: os.Getenv(viper.GetString("security.api_sign.app_security")),
				},
				Jwt: &Jwt{
					ExpiresAt: viper.GetInt("security.jwt.expire_at"),
					Key:       os.Getenv(viper.GetString("security.jwt.key")),
				},
				Oauth2: &Oauth2{
					Google: &Google{
						ClientId:     os.Getenv(viper.GetString("security.oauth2.google.client_id")),
						ClientSecret: os.Getenv(viper.GetString("security.oauth2.google.client_secret")),
						RedirectUrl:  viper.GetString("security.oauth2.google.redirect_url"),
						Scopes:       viper.GetStringSlice("security.oauth2.google.scopes"),
					},
					Github: &Github{
						RedirectUrl:  viper.GetString("security.oauth2.github.redirect_url"),
						ClientId:     os.Getenv(viper.GetString("security.oauth2.github.client_id")),
						ClientSecret: os.Getenv(viper.GetString("security.oauth2.github.client_secret")),
						Scopes:       viper.GetStringSlice("security.oauth2.github.scopes"),
					},
				},
			},
		}
	} else {
		return &Config{
			Security: &Security{
				ApiSign: &ApiSign{
					AppKey:      viper.GetString("security.api_sign.app_key"),
					AppSecurity: viper.GetString("security.api_sign.app_security"),
				},
				Jwt: &Jwt{
					ExpiresAt: viper.GetInt("security.jwt.expire_at"),
					Key:       viper.GetString("security.jwt.key"),
				},
				Oauth2: &Oauth2{
					Google: &Google{
						ClientId:     viper.GetString("security.oauth2.google.client_id"),
						ClientSecret: viper.GetString("security.oauth2.google.client_secret"),
						RedirectUrl:  viper.GetString("security.oauth2.google.redirect_url"),
						Scopes:       viper.GetStringSlice("security.oauth2.google.scopes"),
					},
					Github: &Github{
						RedirectUrl:  viper.GetString("security.oauth2.github.redirect_url"),
						ClientId:     viper.GetString("security.oauth2.github.client_id"),
						ClientSecret: viper.GetString("security.oauth2.github.client_secret"),
						Scopes:       viper.GetStringSlice("security.oauth2.github.scopes"),
					},
				},
			},
		}
	}
}
