package vivo_channel

import (
	"context"
	"fmt"
	"gitee.com/ling-bin/go-push-sdk/push/common/http"
	"gitee.com/ling-bin/go-push-sdk/push/common/intent"
	"gitee.com/ling-bin/go-push-sdk/push/common/json"
	"gitee.com/ling-bin/go-push-sdk/push/common/message"
	"gitee.com/ling-bin/go-push-sdk/push/errcode"
	"gitee.com/ling-bin/go-push-sdk/push/setting"
	"strings"
)

const (
	timeout           = 5
	deviceTokenMax    = 1000
	deviceTokenMin    = 1
	urlBase           = "https://api-push.vivo.com.cn"
	actionSinglePush  = "message/send"
	actionSaveMessage = "message/saveListPayload"
	actionMultiPush   = "message/pushToList"
)

type PushClient struct {
	httpClient *http.Client
	conf       *setting.ConfigVivo
	authClient *AuthToken
}

func NewPushClient(conf *setting.ConfigVivo) (setting.PushClientInterface, error) {
	errCheck := checkConf(conf)
	if errCheck != nil {
		return nil, errCheck
	}
	return &PushClient{
		conf:       conf,
		httpClient: http.NewClient(timeout),
		authClient: NewAuthToken(),
	}, nil
}

func checkConf(conf *setting.ConfigVivo) error {
	if conf.AppPkgName == "" {
		return errcode.ErrVivoAppPkgNameEmpty
	}
	if conf.AppId == "" {
		return errcode.ErrVivoAppIdEmpty
	}
	if conf.AppKey == "" {
		return errcode.ErrVivoAppKeyEmpty
	}
	if conf.AppSecret == "" {
		return errcode.ErrVivoAppSecretEmpty
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

func (p *PushClient) parseBody(body []byte) (*PushMessageResponse, error) {
	resp := &PushMessageResponse{}
	err := json.UnmarshalByte(body, resp)
	if err != nil {
		return nil, errcode.ErrVivoParseBody
	}
	return resp, nil
}

func (p *PushClient) GetAccessToken(ctx context.Context) (interface{}, error) {

	authTokenReq := &AuthTokenReq{
		AppId:     p.conf.AppId,
		AppKey:    p.conf.AppKey,
		AppSecret: p.conf.AppSecret,
	}

	return p.authClient.Get(ctx, authTokenReq)
}

func (p *PushClient) checkParam(pushRequest *setting.PushMessageRequest) error {

	err := message.CheckMessageParam(pushRequest, deviceTokenMin, deviceTokenMax, true)
	if err != nil {
		return err
	}
	if pushRequest.Message.BusinessId == "" {
		return errcode.ErrBusinessIdEmpty
	}

	return nil
}

func (p *PushClient) pushNotice(ctx context.Context, pushRequest *setting.PushMessageRequest) (*PushMessageResponse, error) {

	body, err := p.pushGateWay(ctx, pushRequest)
	if err != nil {
		return nil, err
	}

	return p.parseBody(body)
}

func (p *PushClient) pushGateWay(ctx context.Context, pushRequest *setting.PushMessageRequest) ([]byte, error) {
	if len(pushRequest.DeviceTokens) > deviceTokenMin {
		return p.pushMultiNotify(ctx, pushRequest)
	} else {
		return p.pushSingleNotify(ctx, pushRequest)
	}
}

func (p *PushClient) pushMultiNotify(ctx context.Context, pushRequest *setting.PushMessageRequest) ([]byte, error) {

	saveMessageTaskId, err := p.saveMessageToCloud(ctx, pushRequest)
	if err != nil {
		return nil, err
	}
	pushMultiNotify := &PushMultiNotify{
		RegIds:    pushRequest.DeviceTokens,
		TaskId:    saveMessageTaskId,
		RequestId: pushRequest.Message.BusinessId,
	}
	url := p.buildMultiNotifyUrl()

	param := json.MarshalToStringNoError(pushMultiNotify)
	request, err := p.httpClient.BuildRequest(ctx, "POST", url, param)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("authToken", pushRequest.AccessToken)

	return p.httpClient.Do(ctx, request)
}

func (p *PushClient) saveMessageToCloud(ctx context.Context, pushRequest *setting.PushMessageRequest) (string, error) {

	saveMessageToCloud := &SaveMessageToCloud{
		Title:       pushRequest.Message.Title,
		Content:     pushRequest.Message.Content,
		SkipType:    4,
		SkipContent: intent.GenerateIntent(p.conf.AppPkgName, pushRequest.Message.Extra),
		RequestId:   pushRequest.Message.BusinessId,
		NotifyType:  1,
		Extra: &SaveMessageToCloudExtra{
			CallBack:      pushRequest.Message.CallBack,
			CallBackParam: pushRequest.Message.CallbackParam,
		},
	}

	uri := p.buildSaveMessageToCloudUrl()
	param := json.MarshalToStringNoError(saveMessageToCloud)
	request, err := p.httpClient.BuildRequest(ctx, "POST", uri, param)
	if err != nil {
		return "", err
	}
	request.Header.Set("authToken", pushRequest.AccessToken)
	request.Header.Set("Content-Type", "application/json")
	body, err := p.httpClient.Do(ctx, request)
	if err != nil {
		return "", err
	}
	saveResult := &SaveMessageToCloudResponse{}
	errParse := json.UnmarshalByte(body, saveResult)
	if errParse != nil {
		return "", errcode.ErrVivoParseBody
	}

	return saveResult.TaskId, nil
}

func (p *PushClient) pushSingleNotify(ctx context.Context, pushRequest *setting.PushMessageRequest) ([]byte, error) {

	singleNotify := &PushSingleNotify{
		RegId:       strings.Join(pushRequest.DeviceTokens, ","),
		Title:       pushRequest.Message.Title,
		Content:     pushRequest.Message.Content,
		SkipType:    4,
		SkipContent: intent.GenerateIntent(p.conf.AppPkgName, pushRequest.Message.Extra),
		RequestId:   pushRequest.Message.BusinessId,
		NotifyType:  1,
		Extra: &SingleNotifyExtra{
			CallBack:      pushRequest.Message.CallBack,
			CallBackParam: pushRequest.Message.CallbackParam,
		},
	}

	uri := p.buildSingleNotifyUrl()

	param := json.MarshalToStringNoError(singleNotify)
	request, err := p.httpClient.BuildRequest(ctx, "POST", uri, param)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("authToken", pushRequest.AccessToken)

	return p.httpClient.Do(ctx, request)
}

func (p *PushClient) buildSingleNotifyUrl() string {

	return fmt.Sprintf("%s/%s", urlBase, actionSinglePush)
}

func (p *PushClient) buildSaveMessageToCloudUrl() string {

	return fmt.Sprintf("%s/%s", urlBase, actionSaveMessage)
}

func (p *PushClient) buildMultiNotifyUrl() string {

	return fmt.Sprintf("%s/%s", urlBase, actionMultiPush)
}
