package services

import ()

//	A ConfigService manages configuratio
//
// needed for services and environment vars
type ConfigServiceApi interface {
	GetConfigValueInt64(k string) int64
	GetConfigValueString(k string) string
}
