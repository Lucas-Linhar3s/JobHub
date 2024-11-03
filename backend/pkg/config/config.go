package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

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

// Função principal para carregar as configurações
func LoadAttributes(viper *viper.Viper) *Config {
	env := &Config{}
	missingVars := findMissingVars(viper, reflect.ValueOf(env).Elem())

	err := fmt.Errorf("missing required environment variables:\n -- %v --", strings.Join(missingVars, "\n"))

	if len(missingVars) > 0 {
		panic(err)
	}
	return env
}

// Função recursiva para verificar variáveis obrigatórias e setar valores
func findMissingVars(viper *viper.Viper, v reflect.Value) []string {
	var missingVars []string
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("yaml")

		if tag == "" {
			continue
		}

		parts := splitTag(tag)
		name := parts[0]
		required := len(parts) > 1 && parts[1] == "required"
		isEnvironment := len(parts) > 2 && parts[2] == "environment"
		fieldValue := v.Field(i)

		// Se o campo é um struct, fazemos a chamada recursiva
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
			missingVars = append(missingVars, findMissingVars(viper, fieldValue.Elem())...)
			continue
		}

		// Setando valores de campos não-estruturados
		switch fieldValue.Kind() {
		case reflect.String:
			if isEnvironment {
				value := os.Getenv(viper.GetString(name))
				if required && value == "" {
					missingVars = append(missingVars, name)
				}
				fieldValue.SetString(value)
				break
			}
			value := viper.GetString(name)
			if required && value == "" {
				missingVars = append(missingVars, name)
			}
			fieldValue.SetString(value)
		case reflect.Int:
			value := viper.GetInt(name)
			if required && value == 0 {
				missingVars = append(missingVars, name)
			}
			fieldValue.SetInt(int64(value))
		case reflect.Slice, reflect.Array:
			value := viper.GetStringSlice(name)
			if required && len(value) == 0 {
				missingVars = append(missingVars, name)
			}
			fieldValue.Set(reflect.ValueOf(value))
		case reflect.Bool:
			value := viper.GetBool(name)
			if required && !value {
				missingVars = append(missingVars, name)
			}
			fieldValue.SetBool(value)
		}
	}
	return missingVars
}

// Função para dividir a tag por vírgulas
func splitTag(tag string) []string {
	var result []string
	var current string
	for _, char := range tag {
		if char == ',' {
			result = append(result, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	result = append(result, current)
	return result
}
