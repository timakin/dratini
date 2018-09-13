package pusher

import (
	"net"
	"net/http"
	"time"

	"github.com/timakin/dratini/gcm"
)

func keepAliveInterval(keepAliveTimeout int) int {
	const minInterval = 30
	const maxInterval = 90
	if keepAliveTimeout <= minInterval {
		return keepAliveTimeout
	}
	result := keepAliveTimeout / 3
	if result < minInterval {
		return minInterval
	}
	if result > maxInterval {
		return maxInterval
	}
	return result
}

// InitGCMClient initializes GCMClient which is globally declared.
func InitGCMClient() error {
	// By default, use GCM endpoint. If UseFCM is explicitly enabled via configuration,
	// use FCM endpoint.
	url := gcm.GCMSendEndpoint
	if ConfDratini.Android.UseFCM {
		url = gcm.FCMSendEndpoint
	}

	var err error
	GCMClient, err = gcm.NewClient(url, ConfDratini.Android.ApiKey)
	if err != nil {
		return err
	}

	transport := &http.Transport{
		MaxIdleConnsPerHost: ConfDratini.Android.KeepAliveConns,
		Dial: (&net.Dialer{
			Timeout:   time.Duration(ConfDratini.Android.Timeout) * time.Second,
			KeepAlive: time.Duration(keepAliveInterval(ConfDratini.Android.KeepAliveTimeout)) * time.Second,
		}).Dial,
		IdleConnTimeout: time.Duration(ConfDratini.Android.KeepAliveTimeout) * time.Second,
	}

	GCMClient.Http = &http.Client{
		Transport: transport,
		Timeout:   time.Duration(ConfDratini.Android.Timeout) * time.Second,
	}

	return nil
}

func InitAPNSClient() error {
	var err error
	APNSClient, err = NewApnsClientHttp2(
		ConfDratini.Ios.PemCertPath,
		ConfDratini.Ios.PemKeyPath,
		ConfDratini.Ios.PemKeyPassphrase,
	)
	if err != nil {
		return err
	}
	return nil
}
