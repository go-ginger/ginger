package ginger

import (
	"github.com/kulichak/models"
	"reflect"
)

type Config struct {
	models.IConfig

	CorsEnabled          bool
	CorsAllowOrigins     string
	CorsAllowCredentials string
	CorsAllowHeaders     string
	CorsAllowMethods     string
}

var config Config

func InitializeConfig(input interface{}) {
	v := reflect.Indirect(reflect.ValueOf(input))
	CorsEnabled := v.FieldByName("CorsEnabled")
	CorsAllowOrigins := v.FieldByName("CorsAllowOrigins")
	CorsAllowCredentials := v.FieldByName("CorsAllowCredentials")
	CorsAllowHeaders := v.FieldByName("CorsAllowHeaders")
	CorsAllowMethods := v.FieldByName("CorsAllowMethods")

	config = Config{
		CorsEnabled:          CorsEnabled.String() == "true",
		CorsAllowOrigins:     CorsAllowOrigins.String(),
		CorsAllowCredentials: CorsAllowCredentials.String(),
		CorsAllowHeaders:     CorsAllowHeaders.String(),
		CorsAllowMethods:     CorsAllowMethods.String(),
	}
}
