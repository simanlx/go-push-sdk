package main

import (
	"context"
	"fmt"
	"gitee.com/cristiane/go-push-sdk/push"
	"gitee.com/cristiane/go-push-sdk/push/setting"
	"github.com/google/uuid"
)

func main() {
	//data := []byte("88a53b2706dbc48cea69554bf3bae2bfc5fd582a6a9c97fd578626d621b63b99")
	//has := md5.Sum(data)
	//md5str := fmt.Sprintf("%x", has)
	//fmt.Println(md5str)
	register, err := push.NewRegisterClient("E:\\GoProject\\go-push-sdk\\setting.json")
	if err != nil {
		fmt.Printf("NewRegisterClient err: %v", err)
		return
	}
	iosClient, err := register.GetIosCertClient()
	if err != nil {
		fmt.Printf("GetIosClient err: %v", err)
		return
	}
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
			SubTitle:   "您有待办任务哦-SubTitle",
			Content:    "早上好！新的一天开始了，目前您还有任务需要马上处理~-Content",
			Sound:      "ios_alarm_sos.mp3",
			Badge:      1,
			Extra: map[string]string{
				"type":        "TodoRemind",
				"link_type":   "TaskList",
				"link_params": "[]",
			},
			CallBack:      "",
			CallbackParam: "",
		},
	}
	//bytes, _ := json.Marshal(msg)
	//fmt.Println(string(bytes))
	ctx := context.Background()
	respPush, err := iosClient.PushNotice(ctx, msg)
	if err != nil {
		fmt.Printf("ios push err: %v", err)
		return
	}
	fmt.Printf("ios push resp: %+v", respPush)
}
