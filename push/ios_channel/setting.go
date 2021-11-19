package ios_channel

const (
	PayloadTemplate = `{"aps" : {"content-available": 1,"alert" : {"title" : "%s","subtitle" : "%s","body" : "%s"},"badge" : %d,"sound" : "%s"},"body" : "%s"}`
)

type PushMessageResponse struct {
	StatusCode int    `json:"status_code"`
	APNsId     string `json:"ap_ns_id"`
	Reason     string `json:"reason"`
}
