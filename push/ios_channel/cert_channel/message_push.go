package cert_channel

import (
	"context"
	"crypto/tls"
	"gitee.com/ling-bin/go-push-sdk/push/common/convert"
	"gitee.com/ling-bin/go-push-sdk/push/common/json"
	"gitee.com/ling-bin/go-push-sdk/push/common/message"
	"gitee.com/ling-bin/go-push-sdk/push/errcode"
	"gitee.com/ling-bin/go-push-sdk/push/ios_channel"
	"gitee.com/ling-bin/go-push-sdk/push/setting"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"strings"
	"time"
)

/*
链接：http://events.jianshu.io/p/f4f016848a4e

苹果APNS push推送错误定位/错误码
Status code	Error string	Description
400	BadCollapseId	The collapse identifier exceeds the maximum allowed size
400	BadDeviceToken	The specified device token was bad. Verify that the request contains a valid token and that the token matches the environment.
400	BadExpirationDate	The apns-expiration value is bad.
400	BadMessageId	The apns-id value is bad.
400	BadPriority	The apns-priority value is bad.
400	BadTopic	The apns-topic was invalid.
400	DeviceTokenNotForTopic	The device token does not match the specified topic.
400	DuplicateHeaders	One or more headers were repeated.
400	IdleTimeout	Idle time out.
400	MissingDeviceToken	The device token is not specified in the request :path. Verify that the :path header contains the device token.
400	MissingTopic	The apns-topic header of the request was not specified and was required. The apns-topic header is mandatory when the client is connected using a certificate that supports multiple topics.
400	PayloadEmpty	The message payload was empty.
400	TopicDisallowed	Pushing to this topic is not allowed.
403	BadCertificate	The certificate was bad.
403	BadCertificateEnvironment	The client certificate was for the wrong environment.
403	ExpiredProviderToken	The provider token is stale and a new token should be generated.
403	Forbidden	The specified action is not allowed.
403	InvalidProviderToken	The provider token is not valid or the token signature could not be verified.
403	MissingProviderToken	No provider certificate was used to connect to APNs and Authorization header was missing or no provider token was specified.
404	BadPath	The request contained a bad :path value.
405	MethodNotAllowed	The specified :method was not POST.
410	Unregistered	The device token is inactive for the specified topic. Expected HTTP/2 status code is 410; see Table 8-4.
413	PayloadTooLarge	The message payload was too large. See Creating the Remote Notification Payload for details on maximum payload size.
429	TooManyProviderTokenUpdates	The provider token is being updated too often.
429	TooManyRequests	Too many requests were made consecutively to the same device token.
500	InternalServerError	An internal server error occurred.
503	ServiceUnavailable	The service is unavailable.
503	Shutdown	The server is shutting down.
*/

const (
	deviceTokenMax = 100
	deviceTokenMin = 1
)

type PushClient struct {
	conf      *setting.ConfigIosCert
	client    *apns2.Client //正式环境
	topic     string        //正式环境里的证书topic
	clientBox *apns2.Client //沙盒环境
	topicBox  string        //沙盒环境里的证书topic
}

func NewPushClient(conf *setting.ConfigIosCert) (setting.PushClientInterface, error) {
	errCheck := checkConf(conf)
	if errCheck != nil {
		return nil, errCheck
	}
	cert, err := certificate.FromP12File(conf.CertPath, conf.Password)
	if err != nil {
		return nil, err
	}
	client := &PushClient{
		conf:   conf,
		client: apns2.NewClient(cert).Production(),
	}
	client.topic = client.TopicFromCert(cert)

	//测试地址
	if len(conf.CertPathBox) != 0 {

		certBox, err := certificate.FromP12File(conf.CertPathBox, conf.PasswordBox)
		if err != nil {
			return nil, err
		}
		client.clientBox = apns2.NewClient(certBox).Development()
		client.topicBox = client.TopicFromCert(certBox)
	}
	return client, nil
}

// TopicFromCert extracts topic from a certificate's common name.
func (p *PushClient) TopicFromCert(cert tls.Certificate) string {
	commonName := cert.Leaf.Subject.CommonName
	var topic string
	n := strings.Index(commonName, ":")
	if n != -1 {
		topic = strings.TrimSpace(commonName[n+1:])
	}
	return topic
}

func (p *PushClient) buildRequest(ctx context.Context, pushRequest *setting.PushMessageRequest) (*ios_channel.PushMessageResponse, error) {
	//payloadStr := fmt.Sprintf(ios_channel.PayloadTemplate, pushRequest.Message.Title, pushRequest.Message.SubTitle, pushRequest.Message.Content,
	//	pushRequest.Message.Badge, pushRequest.Message.Sound, "")
	var (
		client *apns2.Client
	)
	if pushRequest.IsSandBox {
		if p.clientBox == nil {
			return nil, errcode.ErrIosBoxEmpty
		}
		client = p.clientBox
		//证书过期判断，如果是过期证书则不推送，2023-03-04
		if time.Now().UTC().Sub(client.Certificate.Leaf.NotAfter).Seconds() >= 0 {
			return nil,errcode.ErrIosBoxNotAfter
		}
	} else {
		client = p.client
		//证书过期判断，如果是过期证书则不推送，2023-03-04
		if time.Now().UTC().Sub(client.Certificate.Leaf.NotAfter).Seconds() >= 0 {
			return nil, errcode.ErrIosNotAfter
		}
	}
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
		DeviceToken: strings.Join(pushRequest.DeviceTokens, ","),
		ApnsID:      pushRequest.Message.BusinessId,
		CollapseID:  pushRequest.Message.BusinessId,
		Payload:     convert.Str2Byte(json.MarshalToStringNoError(pushPayload)),
	}
	if pushRequest.IsSandBox {
		notification.Topic = p.topicBox
	} else {
		notification.Topic = p.topic
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

func checkConf(conf *setting.ConfigIosCert) error {
	if conf.CertPath == "" {
		return errcode.ErrIosCertPathEmpty
	}
	if conf.Password == "" {
		return errcode.ErrIosPasswordEmpty
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

func (p *PushClient) pushNotice(ctx context.Context, pushRequest *setting.PushMessageRequest) (*ios_channel.PushMessageResponse, error) {

	return p.buildRequest(ctx, pushRequest)
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

func (p *PushClient) GetAccessToken(ctx context.Context) (interface{}, error) {

	return nil, nil
}
