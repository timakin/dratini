package dratini

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
)

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

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	var result StatApp
	result.QueueMax = cap(QueueNotification)
	result.QueueUsage = len(QueueNotification)
	result.PusherMax = ConfDratini.Core.PusherMax * ConfDratini.Core.WorkerNum
	result.PusherCount = atomic.LoadInt64(&PusherCountAll)
	result.Ios.PushSuccess = atomic.LoadInt64(&StatDratini.Ios.PushSuccess)
	result.Ios.PushError = atomic.LoadInt64(&StatDratini.Ios.PushError)
	result.Android.PushSuccess = atomic.LoadInt64(&StatDratini.Android.PushSuccess)
	result.Android.PushError = atomic.LoadInt64(&StatDratini.Android.PushError)

	respBody, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		msg := "Response-body could not be created"
		LogError.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Server", serverHeader())
	w.Write(respBody)
}
