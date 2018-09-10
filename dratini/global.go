package dratini

import (
	"net/http"

	"github.com/timakin/dratini/gcm"

	"go.uber.org/zap"
)

var (
	// Toml configuration for Dratini
	ConfDratini ConfToml
	// push notification Queue
	QueueNotification chan RequestDratiniNotification
	// TLS certificate and key for APNs
	CertificatePemIos CertificatePem
	// Stat for Dratini
	StatDratini StatApp
	// http client for APNs and GCM/FCM
	APNSClient *http.Client
	GCMClient  *gcm.Client
	// access and error logger
	LogAccess *zap.Logger
	LogError  *zap.Logger
	// sequence ID for numbering push
	SeqID uint64
)
