package services

import ()

//	A ConfigService manages configuratio
//
// needed for services and environment vars
type ConfigServiceApi interface {
	GetConfigService() *ConfigService
	GetConfigValue(k string) any
}
