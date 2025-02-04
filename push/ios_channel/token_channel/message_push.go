package token_channel

import (
	"context"
	"gitee.com/ling-bin/go-push-sdk/push/common/convert"
	"gitee.com/ling-bin/go-push-sdk/push/common/json"
	"gitee.com/ling-bin/go-push-sdk/push/common/message"
	"gitee.com/ling-bin/go-push-sdk/push/errcode"
	"gitee.com/ling-bin/go-push-sdk/push/ios_channel"
	"gitee.com/ling-bin/go-push-sdk/push/setting"
	"strings"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"
)

const (
	deviceTokenMax = 100
	deviceTokenMin = 1
)

type PushClient struct {
	conf      *setting.ConfigIosToken
	client    *apns2.Client //正式环境
	clientBox *apns2.Client //沙盒环境
}

func NewPushClient(conf *setting.ConfigIosToken) (setting.PushClientInterface, error) {
	errCheck := checkConf(conf)
	if errCheck != nil {
		return nil, errCheck
	}
	authKey, err := token.AuthKeyFromFile(conf.SecretFile)
	if err != nil {
		return nil, err
	}
	tokenClient := &token.Token{
		AuthKey: authKey,
		KeyID:   conf.KeyId,
		TeamID:  conf.TeamId,
	}
	client := &PushClient{
		conf:   conf,
		client: apns2.NewTokenClient(tokenClient).Production(),
	}

	if len(conf.SecretFileBox) != 0 {
		authKeyBox, err := token.AuthKeyFromFile(conf.SecretFile)
		if err != nil {
			return nil, err
		}
		tokenClientBox := &token.Token{
			AuthKey: authKeyBox,
			KeyID:   conf.KeyIdBox,
			TeamID:  conf.TeamIdBox,
		}
		client.clientBox = apns2.NewTokenClient(tokenClientBox).Development()
	}

	return client, nil
}

func checkConf(conf *setting.ConfigIosToken) error {
	if conf.TeamId == "" {
		return errcode.ErrIosTeamIdEmpty
	}
	if conf.KeyId == "" {
		return errcode.ErrIosKeyIdEmpty
	}
	if conf.SecretFile == "" {
		return errcode.ErrIosSecretFileEmpty
	}
	if conf.BundleId == "" {
		return errcode.ErrIosBundleIdEmpty
	}
	return nil
}

func (p *PushClient) PushNotice(ctx context.Context, pushRequest *setting.PushMessageRequest) (interface{}, error) {
	errCheck := p.checkParam(pushRequest)
	if errCheck != nil {
		return nil, errCheck
	}
	return p.pushNotice(ctx, pushRequest)
}

func (p *PushClient) checkParam(pushRequest *setting.PushMessageRequest) error {

	if err := message.CheckMessageParam(pushRequest, deviceTokenMin, deviceTokenMax, false); err != nil {
		return err
	}
	if pushRequest.Message.BusinessId == "" {
		return errcode.ErrBusinessIdEmpty
	}
	// 其余参数检查

	return nil
}

func (p *PushClient) pushNotice(ctx context.Context, pushRequest *setting.PushMessageRequest) (*ios_channel.PushMessageResponse, error) {

	return p.buildRequest(ctx, pushRequest)
}

func (p *PushClient) buildRequest(ctx context.Context, pushRequest *setting.PushMessageRequest) (*ios_channel.PushMessageResponse, error) {
	//bodyJson, _ := json.MarshalToString(pushRequest.Message.Extra)
	//payloadStr := fmt.Sprintf(ios_channel.PayloadTemplate, pushRequest.Message.Title, pushRequest.Message.SubTitle, pushRequest.Message.Content,
	//	pushRequest.Message.Badge, pushRequest.Message.Sound, bodyJson)
	pushPayload := &ios_channel.PushPayload{
		Aps: ios_channel.ApsData{
			ContentAvailable: 1,
			Alert: ios_channel.AlertData{
				Title:    pushRequest.Message.Title,
				Subtitle: pushRequest.Message.SubTitle,
				Body:     pushRequest.Message.Content,
			},
			Badge: pushRequest.Message.Badge,
			Sound: pushRequest.Message.Sound,
		},
		Body: json.MarshalToStringNoError(pushRequest.Message.Extra),
	}
	notification := &apns2.Notification{
		CollapseID:  pushRequest.Message.BusinessId,
		ApnsID:      pushRequest.Message.BusinessId,
		DeviceToken: strings.Join(pushRequest.DeviceTokens, ","),
		Topic:       p.conf.BundleId,
		Payload:     convert.Str2Byte(json.MarshalToStringNoError(pushPayload)),
	}
	var (
		client *apns2.Client
	)
	if pushRequest.IsSandBox {
		if p.clientBox == nil {
			return nil, errcode.ErrIosBoxEmpty
		}
		client = p.clientBox
	} else {
		client = p.client
	}
	res, err := client.PushWithContext(ctx, notification)
	if err != nil {
		return nil, err
	}
	return &ios_channel.PushMessageResponse{
		StatusCode: res.StatusCode,
		APNsId:     res.ApnsID,
		Reason:     res.Reason,
	}, nil
}

func (p *PushClient) GetAccessToken(ctx context.Context) (interface{}, error) {

	return nil, nil
}
