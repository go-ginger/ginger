package ginger

import (
	"github.com/BurntSushi/toml"
	"github.com/go-ginger/models"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"reflect"
)

type Config struct {
	models.IConfig

	LanguageBundle       *i18n.Bundle
	CorsEnabled          bool
	CorsAllowOrigins     string
	CorsAllowCredentials string
	CorsAllowHeaders     string
	CorsAllowMethods     string
}

var config Config

func (c *Config) InitializeConfig(input interface{}) {
}

func InitializeConfig(conf models.IConfig, input interface{}) {
	v := reflect.Indirect(reflect.ValueOf(input))
	CorsEnabled := v.FieldByName("CorsEnabled")
	CorsAllowOrigins := v.FieldByName("CorsAllowOrigins")
	CorsAllowCredentials := v.FieldByName("CorsAllowCredentials")
	CorsAllowHeaders := v.FieldByName("CorsAllowHeaders")
	CorsAllowMethods := v.FieldByName("CorsAllowMethods")
	config = Config{
		IConfig:              conf,
		CorsEnabled:          CorsEnabled.Interface() == true,
		CorsAllowOrigins:     CorsAllowOrigins.String(),
		CorsAllowCredentials: CorsAllowCredentials.String(),
		CorsAllowHeaders:     CorsAllowHeaders.String(),
		CorsAllowMethods:     CorsAllowMethods.String(),
		LanguageBundle:       i18n.NewBundle(language.English),
	}
	config.LanguageBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	config.IConfig.InitializeConfig(&config)
}
