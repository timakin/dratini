package dratini

import (
	"runtime"
	"testing"

	_ "github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	ConfDratiniPath = "../conf/dratini.toml"
)

type ConfigTestSuite struct {
	suite.Suite
	ConfDratiniDefault ConfToml
	ConfDratini        ConfToml
}

func (suite *ConfigTestSuite) SetupTest() {
	suite.ConfDratiniDefault = BuildDefaultConf()
	suite.ConfDratini = BuildDefaultConf()
	var err error
	suite.ConfDratini, err = LoadConf(suite.ConfDratini, ConfDratiniPath)
	if err != nil {
		panic("failed to load " + ConfDratiniPath)
	}
}

func (suite *ConfigTestSuite) TestValidateConfDefault() {
	// Core
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Core.Port, "1056")
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Core.WorkerNum, int64(runtime.NumCPU()))
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Core.QueueNum, int64(8192))
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Core.NotificationMax, int64(100))
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Core.PusherMax, int64(0))
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Core.Pid, "")
	// Android
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Android.Enabled, true)
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Android.ApiKey, "")
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Android.Timeout, 5)
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Android.KeepAliveTimeout, 90)
	assert.Equal(suite.T(), int64(suite.ConfDratiniDefault.Android.KeepAliveConns), suite.ConfDratiniDefault.Core.WorkerNum)
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Android.RetryMax, 1)
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Android.UseFCM, false)
	// Ios
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Ios.Enabled, true)
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Ios.PemCertPath, "")
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Ios.PemKeyPath, "")
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Ios.Sandbox, true)
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Ios.RetryMax, 1)
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Ios.Timeout, 5)
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Ios.KeepAliveTimeout, 90)
	assert.Equal(suite.T(), int64(suite.ConfDratiniDefault.Ios.KeepAliveConns), suite.ConfDratiniDefault.Core.WorkerNum)
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Ios.Topic, "")
	// Log
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Log.AccessLog, "stdout")
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Log.ErrorLog, "stderr")
	assert.Equal(suite.T(), suite.ConfDratiniDefault.Log.Level, "error")
}

func (suite *ConfigTestSuite) TestValidateConf() {
	// Core
	assert.Equal(suite.T(), suite.ConfDratini.Core.Port, "1056")
	assert.Equal(suite.T(), suite.ConfDratini.Core.WorkerNum, int64(8))
	assert.Equal(suite.T(), suite.ConfDratini.Core.QueueNum, int64(8192))
	assert.Equal(suite.T(), suite.ConfDratini.Core.NotificationMax, int64(100))
	assert.Equal(suite.T(), suite.ConfDratini.Core.PusherMax, int64(0))
	assert.Equal(suite.T(), suite.ConfDratini.Core.Pid, "")
	// Android
	assert.Equal(suite.T(), suite.ConfDratini.Android.Enabled, true)
	assert.Equal(suite.T(), suite.ConfDratini.Android.ApiKey, "apikey for GCM")
	assert.Equal(suite.T(), suite.ConfDratini.Android.Timeout, 5)
	assert.Equal(suite.T(), suite.ConfDratini.Android.KeepAliveTimeout, 30)
	assert.Equal(suite.T(), suite.ConfDratini.Android.KeepAliveConns, 4)
	assert.Equal(suite.T(), suite.ConfDratini.Android.RetryMax, 1)
	assert.Equal(suite.T(), suite.ConfDratini.Android.UseFCM, false)
	// Ios
	assert.Equal(suite.T(), suite.ConfDratini.Ios.Enabled, true)
	assert.Equal(suite.T(), suite.ConfDratini.Ios.PemCertPath, "cert.pem")
	assert.Equal(suite.T(), suite.ConfDratini.Ios.PemKeyPath, "key.pem")
	assert.Equal(suite.T(), suite.ConfDratini.Ios.Sandbox, true)
	assert.Equal(suite.T(), suite.ConfDratini.Ios.RetryMax, 1)
	assert.Equal(suite.T(), suite.ConfDratini.Ios.Timeout, 5)
	assert.Equal(suite.T(), suite.ConfDratini.Ios.KeepAliveTimeout, 30)
	// Log
	assert.Equal(suite.T(), suite.ConfDratini.Ios.KeepAliveConns, 6)
	assert.Equal(suite.T(), suite.ConfDratini.Log.AccessLog, "stdout")
	assert.Equal(suite.T(), suite.ConfDratini.Log.ErrorLog, "stderr")
	assert.Equal(suite.T(), suite.ConfDratini.Log.Level, "error")
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
