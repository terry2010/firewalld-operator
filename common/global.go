package Common

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
)

var Config = viper.New()

var Json = jsoniter.ConfigCompatibleWithStandardLibrary
