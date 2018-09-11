package dratini

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"time"
	"os"
	"github.com/timakin/dratini/gcm"
)

type DratiniPush struct {
	Notifications []DratiniPushNotification `json:"notifications"`
}

type DratiniPushNotification struct {
	// Common
	Tokens     []string `json:"token"`
	Platform   int      `json:"platform"`
	Message    string   `json:"message"`
	Identifier string   `json:"identifier,omitempty"`
	// Android
	CollapseKey    string `json:"collapse_key,omitempty"`
	DelayWhileIdle bool   `json:"delay_while_idle,omitempty"`
	TimeToLive     int    `json:"time_to_live,omitempty"`
	// iOS
	Title            string       `json:"title,omitempty"`
	Subtitle         string       `json:"subtitle,omitempty"`
	Badge            int          `json:"badge,omitempty"`
	Category         string       `json:"category,omitempty"`
	Sound            string       `json:"sound,omitempty"`
	ContentAvailable bool         `json:"content_available,omitempty"`
	MutableContent   bool         `json:"mutable_content,omitempty"`
	Expiry           int          `json:"expiry,omitempty"`
	Retry            int          `json:"retry,omitempty"`
	Extend           []ExtendJSON `json:"extend,omitempty"`
	// meta
	ID uint64 `json:"seq_id,omitempty"`
}

type ExtendJSON struct {
	Key   string `json:"key"`
	Value string `json:"val"`
}

type ResponseDratini struct {
	Message string `json:"message"`
}

type CertificatePem struct {
	Cert []byte
	Key  []byte
}

func enqueueNotifications(notifications []DratiniPushNotification) {
	BootWg.Add(1)
	for _, notification := range notifications {
		err := validateNotification(&notification)
		if err != nil {
			LogError.Error(err.Error())
			continue
		}
		var enabledPush bool
		switch notification.Platform {
		case PlatFormIos:
			enabledPush = ConfDratini.Ios.Enabled
		case PlatFormAndroid:
			enabledPush = ConfDratini.Android.Enabled
		}
		// Enqueue notification per token
		for _, token := range notification.Tokens {
			notification2 := notification
			notification2.Tokens = []string{token}
			notification2.ID = numberingPush()
			if enabledPush {
				LogPush(notification2.ID, StatusAcceptedPush, token, 0, notification2, nil)
				QueueNotification <- notification2
			} else {
				LogPush(notification2.ID, StatusDisabledPush, token, 0, notification2, nil)
			}
		}
	}
	BootWg.Done()
}

func pushNotificationIos(req DratiniPushNotification) error {
	LogError.Debug("START push notification for iOS")

	service := NewApnsServiceHttp2(APNSClient)

	token := req.Tokens[0]

	headers := NewApnsHeadersHttp2(&req)
	payload := NewApnsPayloadHttp2(&req)

	stime := time.Now()
	err := ApnsPushHttp2(token, service, headers, payload)

	etime := time.Now()
	ptime := etime.Sub(stime).Seconds()

	if err != nil {
		atomic.AddInt64(&StatDratini.Ios.PushError, 1)
		LogPush(req.ID, StatusFailedPush, token, ptime, req, err)
		return err
	}

	atomic.AddInt64(&StatDratini.Ios.PushSuccess, 1)
	LogPush(req.ID, StatusSucceededPush, token, ptime, req, nil)

	LogError.Debug("END push notification for iOS")

	return nil
}

func pushNotificationAndroid(req DratiniPushNotification) error {
	LogError.Debug("START push notification for Android")

	data := map[string]interface{}{"message": req.Message}
	if len(req.Extend) > 0 {
		for _, extend := range req.Extend {
			data[extend.Key] = extend.Value
		}
	}

	token := req.Tokens[0]

	msg := gcm.NewMessage(data, token)
	msg.CollapseKey = req.CollapseKey
	msg.DelayWhileIdle = req.DelayWhileIdle
	msg.TimeToLive = req.TimeToLive

	stime := time.Now()
	resp, err := GCMClient.SendNoRetry(msg)
	etime := time.Now()
	ptime := etime.Sub(stime).Seconds()
	if err != nil {
		atomic.AddInt64(&StatDratini.Android.PushError, 1)
		LogPush(req.ID, StatusFailedPush, token, ptime, req, err)
		return err
	}

	if resp.Failure > 0 {
		atomic.AddInt64(&StatDratini.Android.PushSuccess, int64(resp.Success))
		atomic.AddInt64(&StatDratini.Android.PushError, int64(resp.Failure))
		LogPush(req.ID, StatusFailedPush, token, ptime, req, errors.New(resp.Results[0].Error))
		return errors.New(resp.Results[0].Error)
	}

	LogPush(req.ID, StatusSucceededPush, token, ptime, req, nil)

	atomic.AddInt64(&StatDratini.Android.PushSuccess, int64(len(req.Tokens)))
	LogError.Debug("END push notification for Android")

	return nil
}

func validateNotification(notification *DratiniPushNotification) error {
	for _, token := range notification.Tokens {
		if len(token) == 0 {
			return errors.New("empty token")
		}
	}

	if notification.Platform < 1 || notification.Platform > 2 {
		return errors.New("invalid platform")
	}

	if len(notification.Message) == 0 {
		return errors.New("empty message")
	}

	return nil
}


func StartBatch() {
	LogError.Debug("push-batch started")

	// file load
	var dratiniPush DratiniPush
	targetFile, err := os.Open(TargetFilePath)
	if err != nil {
		LogError.Error(err.Error())
		return
	}

	err = json.NewDecoder(targetFile).Decode(&dratiniPush)

	if err != nil {
		LogError.Error(err.Error())
		return
	}

	if len(dratiniPush.Notifications) == 0 {
		LogError.Error("empty notification")
		return
	} else if int64(len(dratiniPush.Notifications)) > ConfDratini.Core.NotificationMax {
		msg := fmt.Sprintf("number of notifications(%d) over limit(%d)", len(dratiniPush.Notifications), ConfDratini.Core.NotificationMax)
		LogError.Error(msg)
		return
	}

	LogError.Debug("enqueue notification")
	enqueueNotifications(dratiniPush.Notifications)

	LogError.Debug("push-batch finished")
}
