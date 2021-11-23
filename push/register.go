package push

import (
	"gitee.com/ling-bin/go-push-sdk/push/common/convert"
	"gitee.com/ling-bin/go-push-sdk/push/common/file"
	"gitee.com/ling-bin/go-push-sdk/push/common/json"
	"gitee.com/ling-bin/go-push-sdk/push/errcode"
	"gitee.com/ling-bin/go-push-sdk/push/huawei_channel"
	"gitee.com/ling-bin/go-push-sdk/push/ios_channel/cert_channel"
	"gitee.com/ling-bin/go-push-sdk/push/ios_channel/token_channel"
	"gitee.com/ling-bin/go-push-sdk/push/meizu_channel"
	"gitee.com/ling-bin/go-push-sdk/push/oppo_channel"
	"gitee.com/ling-bin/go-push-sdk/push/setting"
	"gitee.com/ling-bin/go-push-sdk/push/vivo_channel"
	"gitee.com/ling-bin/go-push-sdk/push/xiaomi_channel"
	"sync"
)

const (
	//DefaultConfFile 推荐的配置文件存放路径
	DefaultConfFile = "/usr/local/etc/go-push-sdk/setting.json"
)

//RegisterClient  厂商客户端连接
type RegisterClient struct {
	cfg    interface{} //配置
	client sync.Map    //连接客户端
}

func NewRegisterClient(configFilePath string) (*RegisterClient, error) {
	if configFilePath == "" {
		return nil, errcode.ErrCfgFileEmpty
	}
	fileRead := file.NewFileRead()
	jsonByte, err := fileRead.Read(configFilePath)
	if err != nil {
		return nil, errcode.ErrParseConfigFile
	}
	return NewRegisterClientWithConf(convert.Byte2Str(jsonByte), "")
}

func NewRegisterClientMap(configMap map[string]map[string]string) (*RegisterClient, error) {
	jsonByte,err:= json.Marshal(configMap)
	if  err != nil{
		return nil, err
	}
	return NewRegisterClientWithConf(convert.Byte2Str(jsonByte), "")
}

func newRegisterClient(cfgJson string, obj interface{}) (*RegisterClient, error) {
	if cfgJson == "" {
		return nil, errcode.ErrConfigEmpty
	}
	err := json.Unmarshal(cfgJson, obj)
	if err != nil {
		return nil, errcode.ErrParseConfigFile
	}
	return &RegisterClient{
		cfg: obj,
	}, nil
}

func NewRegisterClientWithConf(cfgJson string, platformType setting.PlatformType) (*RegisterClient, error) {
	var obj interface{}
	switch platformType {
	case setting.HuaweiPlatform:
		obj = &setting.ConfigHuawei{}
	case setting.MeizuPlatform:
		obj = &setting.ConfigMeizu{}
	case setting.OppoPlatform:
		obj = &setting.ConfigOppo{}
	case setting.VivoPlatform:
		obj = &setting.ConfigVivo{}
	case setting.XiaomiPlatform:
		obj = &setting.ConfigXiaomi{}
	case setting.IosCertPlatform:
		obj = &setting.ConfigIosCert{}
	case setting.IosTokenPlatform:
		obj = &setting.ConfigIosToken{}
	default:
		obj = &setting.PushConfig{}
	}
	return newRegisterClient(cfgJson, obj)
}

func (r *RegisterClient) GetPlatformClient(platform setting.PlatformType) (setting.PushClientInterface, error) {
	switch platform {
	case setting.HuaweiPlatform:
		return r.GetHUAWEIClient()
	case setting.MeizuPlatform:
		return r.GetMEIZUClient()
	case setting.OppoPlatform:
		return r.GetOPPOClient()
	case setting.VivoPlatform:
		return r.GetVIVOClient()
	case setting.XiaomiPlatform:
		return r.GetXIAOMIClient()
	case setting.IosCertPlatform:
		return r.GetIosCertClient()
	case setting.IosTokenPlatform:
		return r.GetIosTokenClient()
	default:
		return nil, errcode.ErrUnknownPlatform
	}
}

func (r *RegisterClient) GetHUAWEIClient() (client setting.PushClientInterface, err error) {
	value, ok := r.client.Load(setting.HuaweiPlatform)
	if ok {
		return value.(*huawei_channel.PushClient), nil
	}
	if conf, ok := r.cfg.(*setting.ConfigHuawei); ok {
		client, err = huawei_channel.NewPushClient(conf)
	} else {
		client, err = huawei_channel.NewPushClient(&r.cfg.(*setting.PushConfig).ConfigHuawei)
	}
	if err == nil {
		r.client.Store(setting.HuaweiPlatform, client)
	}
	return client, err
}

func (r *RegisterClient) GetMEIZUClient() (client setting.PushClientInterface, err error) {
	value, ok := r.client.Load(setting.MeizuPlatform)
	if ok {
		return value.(*meizu_channel.PushClient), nil
	}
	if conf, ok := r.cfg.(*setting.ConfigMeizu); ok {
		client, err = meizu_channel.NewPushClient(conf)
	}else {
		client, err = meizu_channel.NewPushClient(&r.cfg.(*setting.PushConfig).ConfigMeizu)
	}
	if err == nil {
		r.client.Store(setting.MeizuPlatform, client)
	}
	return client, err
}

func (r *RegisterClient) GetXIAOMIClient() (client setting.PushClientInterface, err error) {
	value, ok := r.client.Load(setting.XiaomiPlatform)
	if ok {
		return value.(*xiaomi_channel.PushClient), nil
	}
	if conf, ok := r.cfg.(*setting.ConfigXiaomi); ok {
		client,err = xiaomi_channel.NewPushClient(conf)
	}else {
		client, err = xiaomi_channel.NewPushClient(&r.cfg.(*setting.PushConfig).ConfigXiaomi)
	}
	if err == nil {
		r.client.Store(setting.XiaomiPlatform, client)
	}
	return client,err
}

func (r *RegisterClient) GetOPPOClient() (client setting.PushClientInterface, err error) {
	value, ok := r.client.Load(setting.OppoPlatform)
	if ok {
		return value.(*oppo_channel.PushClient), nil
	}
	if conf, ok := r.cfg.(*setting.ConfigOppo); ok {
		client, err = oppo_channel.NewPushClient(conf)
	}else {
		client, err = oppo_channel.NewPushClient(&r.cfg.(*setting.PushConfig).ConfigOppo)
	}
	if err == nil {
		r.client.Store(setting.OppoPlatform, client)
	}
	return client,err
}

func (r *RegisterClient) GetVIVOClient() (client setting.PushClientInterface, err error) {
	value, ok := r.client.Load(setting.VivoPlatform)
	if ok {
		return value.(*vivo_channel.PushClient), nil
	}
	if conf, ok := r.cfg.(*setting.ConfigVivo); ok {
		client, err = vivo_channel.NewPushClient(conf)
	}else {
		client, err = vivo_channel.NewPushClient(&r.cfg.(*setting.PushConfig).ConfigVivo)
	}
	if err == nil {
		r.client.Store(setting.VivoPlatform, client)
	}
	return client,err
}

func (r *RegisterClient) GetIosCertClient() (client setting.PushClientInterface, err error)  {
	value, ok := r.client.Load(setting.IosCertPlatform)
	if ok {
		return value.(*cert_channel.PushClient), nil
	}
	if conf, ok := r.cfg.(*setting.ConfigIosCert); ok {
		client, err = cert_channel.NewPushClient(conf)
	}else {
		client, err = cert_channel.NewPushClient(&r.cfg.(*setting.PushConfig).ConfigIosCert)
	}
	if err == nil {
		r.client.Store(setting.IosCertPlatform, client)
	}
	return client,err
}

func (r *RegisterClient) GetIosTokenClient() (client setting.PushClientInterface, err error) {
	value, ok := r.client.Load(setting.IosTokenPlatform)
	if ok {
		return value.(*token_channel.PushClient), nil
	}
	if conf, ok := r.cfg.(*setting.ConfigIosToken); ok {
		client, err = token_channel.NewPushClient(conf)
	}else {
		client, err = token_channel.NewPushClient(&r.cfg.(*setting.PushConfig).ConfigIosToken)
	}
	if err == nil {
		r.client.Store(setting.IosTokenPlatform, client)
	}
	return client,err
}
