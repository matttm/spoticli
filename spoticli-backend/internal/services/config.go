package services

/// TODO: use go:generate to make singleton

import (
	"context"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

//	A ConfigService manages configuratio
//
// needed for services and environment vars
type ConfigService struct {
	CloudConfig aws.Config
	Config      map[string]string
}

var configService *ConfigService

func GetConfigService() *ConfigService {
	if configService == nil {
		configService = &ConfigService{}
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			panic(err)
		}
		configService.CloudConfig = cfg
		configService.Config = map[string]string{}
		configService.Config["STREAM_SEGMENT_SIZE"] = os.Getenv("STREAM_SEGMENT_SIZE")
		configService.Config["FRAME_CLUSTER_SIZE"] = os.Getenv("FRAME_CLUSTER_SIZE")
		println("ConfigService Instantiated")
	}
	return configService
}

// GetConfigValue gets an config var
func (cs *ConfigService) GetConfigValueString(k string) string {
	return cs.Config[k]
}
func (cs *ConfigService) GetConfigValueInt64(k string) int64 {
	intString := cs.Config[k]
	_int, _ := strconv.Atoi(intString)
	return int64(_int)
}
