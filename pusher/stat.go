package pusher

type StatApp struct {
	QueueMax    int         `json:"queue_max"`
	QueueUsage  int         `json:"queue_usage"`
	PusherMax   int64       `json:"pusher_max"`
	PusherCount int64       `json:"pusher_count"`
	Ios         StatIos     `json:"ios"`
	Android     StatAndroid `json:"android"`
}

type StatAndroid struct {
	PushSuccess int64 `json:"push_success"`
	PushError   int64 `json:"push_error"`
}

type StatIos struct {
	PushSuccess int64 `json:"push_success"`
	PushError   int64 `json:"push_error"`
}

func InitStat() {
	StatDratini.QueueUsage = 0
	StatDratini.PusherCount = 0
	StatDratini.Ios.PushSuccess = 0
	StatDratini.Ios.PushError = 0
	StatDratini.Android.PushSuccess = 0
	StatDratini.Android.PushError = 0
}
