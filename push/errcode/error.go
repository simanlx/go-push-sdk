package errcode

import (
	"errors"
)


var (
	ErrCfgFileEmpty            = errors.New("conf file empty")
	ErrHuaweiClientIdEmpty     = errors.New( "huawei clientId Empty")
	ErrHuaweiClientSecretEmpty = errors.New( "huawei clientSecret Empty")
	ErrXiaomiAppSecretEmpty    = errors.New( "xiaomi appSecret Empty")
	ErrVivoAppSecretEmpty      = errors.New( "vivo appSecret Empty")
	ErrMeizuAppSecretEmpty     = errors.New( "meizu appSecret Empty")
	ErrXiaomiAppPkgNameEmpty   = errors.New( "xiaomi appPkgName Empty")
	ErrVivoAppPkgNameEmpty     = errors.New( "vivo appPkgName Empty")
	ErrOppoAppPkgNameEmpty     = errors.New( "oppo appPkgName Empty")
	ErrMeizuAppPkgNameEmpty    = errors.New( "meizu appPkgName Empty")
	ErrHuaweiAppPkgNameEmpty   = errors.New( "huawei appPkgName Empty")
	ErrAccessTokenEmpty        = errors.New( "accessToken Empty")
	ErrMessageTitleEmpty       = errors.New( "message title empty")
	ErrMessageContentEmpty     = errors.New( "message content empty")
	ErrDeviceTokenMax          = errors.New( "device token max limited")
	ErrDeviceTokenMin          = errors.New( "device token min limited")
	ErrVivoAppIdEmpty          = errors.New( "vivo appId Empty")
	ErrMeizuAppIdEmpty         = errors.New( "meizu appId Empty")
	ErrVivoAppKeyEmpty         = errors.New( "vivo appKey Empty")
	ErrOppoAppKeyEmpty         = errors.New( "oppo appKey Empty")
	ErrOppoMasterSecretEmpty   = errors.New( "oppo masterSecret Empty")
	ErrOppoSaveMessageToCloud  = errors.New( "oppo save Message to Cloud error")
	ErrBusinessIdEmpty         = errors.New( "businessId Empty")
	ErrXiaomiParseBody         = errors.New( "xiaomi parse response body error")
	ErrVivoParseBody           = errors.New( "vivo parse response body error")
	ErrOppoParseBody           = errors.New( "oppo parse response body error")
	ErrMeizuParseBody          = errors.New( "meizu parse response body error")
	ErrHuaweiParseBody         = errors.New( "huawei parse response body error")
	ErrIosCertPathEmpty        = errors.New( "ios certPath Empty")
	ErrIosPasswordEmpty        = errors.New( "ios passWord Empty")
	ErrIosBoxEmpty             = errors.New( "ios Box Empty")
	ErrIosTeamIdEmpty          = errors.New( "ios teamId Empty")
	ErrIosKeyIdEmpty           = errors.New( "ios keyId Empty")
	ErrIosSecretFileEmpty      = errors.New( "ios secretFile Empty")
	ErrIosBundleIdEmpty        = errors.New( "ios bundleId Empty")
	ErrIosBoxNotAfter          = errors.New( "ios Box tls file out Expirytime")
	ErrIosNotAfter             = errors.New( "ios tls file out Expirytime")
	ErrUnknownPlatform         = errors.New( "unknown platform ")
	ErrParseConfigFile         = errors.New( "parse config file err")
	ErrConfigEmpty             = errors.New( "config file null")
)
