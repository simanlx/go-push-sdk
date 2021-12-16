package ios_channel

const (
	PayloadTemplate = `{"aps" : {"content-available": 1,"alert" : {"title" : "%s","subtitle" : "%s","body" : "%s"},"badge" : %d,"sound" : "%s"},"body" : "%s"}`
)

type PushPayload struct {
	Aps  ApsData `json:"aps"`
	Body string  `json:"body"`
}

type ApsData struct {
	ContentAvailable int       `json:"content_available"`
	Alert            AlertData `json:"alert"`
	Badge            int       `json:"badge"`
	Sound            string    `json:"sound"`
}

type AlertData struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Body     string `json:"body"`
}

type PushMessageResponse struct {
	StatusCode int    `json:"status_code"`
	APNsId     string `json:"ap_ns_id"`
	Reason     string `json:"reason"`
}
