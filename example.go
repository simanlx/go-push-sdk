package main

import (
	"context"
	"fmt"
	"gitee.com/ling-bin/go-push-sdk/push"
	"gitee.com/ling-bin/go-push-sdk/push/common/json"
	"gitee.com/ling-bin/go-push-sdk/push/huawei_channel"
	"gitee.com/ling-bin/go-push-sdk/push/setting"
	"github.com/google/uuid"
)

func main() {
	ios()
}

func xiaomi()  {
	register, err := push.NewRegisterClient("E:\\GoProject\\go-push-sdk\\setting.json")
	if err != nil {
		fmt.Printf("NewRegisterClient err: %v", err)
		return
	}
	client, err := register.GetPlatformClient("xiaomi")
	if err != nil {
		fmt.Printf("GetIosClient err: %v", err)
		return
	}
	var deviceTokens = []string{
		"2L2MTwfXvRr5UllZXG2eFOLmdsSRKJJOiH+pEO0M85vhZWxR3pxWIDeV8FEmb8Xr",
	}
	msg := &setting.PushMessageRequest{
		AccessToken:  "",
		DeviceTokens: deviceTokens,
		IsSandBox:    true,
		Message: &setting.Message{
			BusinessId: uuid.New().String(),
			Title:      "待办任务提醒-Title",
			SubTitle:   "您有待办哦-SubTitle",
			Content:    "早上好！一天开始了-Content",
			Sound:      "ios_alarm_sos.mp3",
			Badge:      1,
			Extra: map[string]string{
				"type":        "4",
				"link_type":   "1",
				"link_params": "",
			},
			CallBack:      "",
			CallbackParam: "",
		},
	}
	//bytes, _ := json.Marshal(msg)
	//fmt.Println(string(bytes))
	ctx := context.Background()
	respPush, err := client.PushNotice(ctx, msg)
	if err != nil {
		fmt.Printf("ios push err: %v", err)
		return
	}
	fmt.Printf("ios push resp: %+v", respPush)
}

//测试通过
func huawei()  {
	//data := []byte("88a53b2706dbc48cea69554bf3bae2bfc5fd582a6a9c97fd578626d621b63b99")
	//has := md5.Sum(data)
	//md5str := fmt.Sprintf("%x", has)
	//fmt.Println(md5str)
	register, err := push.NewRegisterClient("E:\\GoProject\\go-push-sdk\\setting.json")
	if err != nil {
		fmt.Printf("NewRegisterClient err: %v", err)
		return
	}
	client, err := register.GetPlatformClient("huawei")
	if err != nil {
		fmt.Printf("GetIosClient err: %v", err)
		return
	}
	var deviceTokens = []string{
		"ACIym4TyPki6h0k1F4NB_QUVO48mZ-D9HSvBVeKnPIIFNrSmdR_OjqxzUxP1WKCDygPtT5EUVL35d5iN8bzW2hc7WN-tOsshhEHV4KcGTpW5LRP3Trk0uRQrjYACqN5V6w",
	}
	token, err := client.GetAccessToken(context.Background())
	if err != nil{
		fmt.Printf("GetToken err: %v", err)
		return
	}
	accessToken := token.(*huawei_channel.AccessTokenResp)
	msg := &setting.PushMessageRequest{
		AccessToken:  accessToken.AccessToken,
		DeviceTokens: deviceTokens,
		IsSandBox:    true,
		Message: &setting.Message{
			BusinessId: uuid.New().String(),
			Title:      "待办任务提醒-Title",
			SubTitle:   "您有待办哦-SubTitle",
			Content:    "早上好！一天开始了-Content",
			Sound:      "ios_alarm_sos.mp3",
			Badge:      1,
			Extra: map[string]string{
				"type":        "4",
				"link_type":   "1",
				"link_params": "[]",
			},
			CallBack:      "",
			CallbackParam: "",
		},
	}
	//bytes, _ := json.Marshal(msg)
	//fmt.Println(string(bytes))
	ctx := context.Background()
	respPush, err := client.PushNotice(ctx, msg)
	if err != nil {
		fmt.Printf("ios push err: %v", err)
		return
	}
	fmt.Printf("ios push resp: %+v", respPush)
}

//测试通过
func ios()  {
	//data := []byte("88a53b2706dbc48cea69554bf3bae2bfc5fd582a6a9c97fd578626d621b63b99")
	//has := md5.Sum(data)
	//md5str := fmt.Sprintf("%x", has)
	//fmt.Println(md5str)
	register, err := push.NewRegisterClient("E:\\GoProject\\go-push-sdk\\setting.json")
	if err != nil {
		fmt.Printf("NewRegisterClient err: %v", err)
		return
	}
	iosClient, err := register.GetPlatformClient("ios")
	if err != nil {
		fmt.Printf("GetIosClient err: %v", err)
		return
	}
	_,_ =register.GetIosCertClient()
	var deviceTokens = []string{
		"88a53b2706dbc48cea69554bf3bae2bfc5fd582a6a9c97fd578626d621b63b99",
	}
	msg := &setting.PushMessageRequest{
		AccessToken:  "",
		DeviceTokens: deviceTokens,
		IsSandBox:    true,
		Message: &setting.Message{
			BusinessId: uuid.New().String(),
			Title:      "待办任务提醒-Title",
			SubTitle:   "您有待办哦-SubTitle",
			Content:    "早上好！一天开始了-Content",
			Sound:      "ios_alarm_sos.mp3",
			Badge:      1,
			Extra: map[string]string{
				"type":        "4",
				"link_type":   "1",
				"link_params": "[]",
			},
			CallBack:      "",
			CallbackParam: "",
		},
	}
	bytes, _ := json.Marshal(msg)
	fmt.Println(string(bytes))
	ctx := context.Background()
	respPush, err := iosClient.PushNotice(ctx, msg)
	if err != nil {
		fmt.Printf("ios push err: %v", err)
		return
	}
	fmt.Printf("ios push resp: %+v", respPush)
}
/*
	register, err := push.NewRegisterClientMap(map[string]map[string]string{
			"huawei": {
				"appPkgName": "应用包名",
				"clientId": "用户在联盟申请的APPID",
				"clientSecret": "应用ID对应的秘钥",
			},
			"meizu": {
				"appPkgName": "应用包名",
				"appId": "应用ID",
				"appSecret": "应用秘钥",
			},
			"xiaomi": {
				"appPkgName": "应用包名",
				"appSecret": "应用秘钥",
			},
			"oppo": {
				"appPkgName": "应用包名",
				"appKey": "应用key",
				"masterSecret": "主秘钥",
			},
			"vivo": {
				"appPkgName": "应用包名",
				"appId": "应用ID",
				"appKey": "应用key",
				"appSecret": "应用秘钥",
			},
			"ios": {
				"certPath": "E:\\GoProject\\dev_push.p12",
				"certPathComment": ".p12格式推送证书绝对路径",
				"password": "123456",
				"passwordComment": "推送证书密码",
				"certPathBox": "E:\\GoProject\\dev_push.p12",
				"certPathCommentBox": ".p12格式推送证书绝对路径[沙盒环境]",
				"passwordBox": "123456",
				"passwordCommentBox": "推送证书密码[沙盒环境]",
			},
			"ios-token": {
				"teamId": "xxxxx",
				"teamIdComment": "开发者帐号teamId",
				"keyId": "xxxxx",
				"keyIdComment": "token认证keyId",
				"secretFile": "xxx.p8",
				"secretFileComment": "token认证密钥文件本地绝对路径",
				"bundleId": "com.xxx.xxx",
				"bundleIdComment": "应用ID",
				"teamIdBox": "xxxxx",
				"teamIdCommentBox": "开发者帐号teamId[沙盒环境]",
				"keyIdBox": "xxxxx",
				"keyIdCommentBox": "token认证keyId[沙盒环境]",
				"secretFileBox": "xxx.p8",
				"secretFileCommentBox": "token认证密钥文件本地绝对路径[沙盒环境]",
				"bundleIdBox": "com.xxx.xxx",
				"bundleIdCommentBox": "应用ID[沙盒环境]",
			},
	})
*/