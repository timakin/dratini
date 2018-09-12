package dratini

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/timakin/dratini/dratini"
)

func main() {
	confPath := flag.String("c", "", "configuration file path for dratini")
	targetPath := flag.String("t", "", "configuration file path for dratini push target")
	workerNum := flag.Int64("w", 0, "number of workers for push notification")
	queueNum := flag.Int64("q", 0, "size of internal queue for push notification")
	flag.Parse()

	// set default parameters
	dratini.ConfDratini = dratini.BuildDefaultConf()

	// load configuration
	conf, err := dratini.LoadConf(dratini.ConfDratini, *confPath)
	if err != nil {
		dratini.LogSetupFatal(err)
	}
	dratini.ConfDratini = conf

	// exit if push target is not specified by flags
	if *targetPath == "" {
		dratini.LogSetupFatal(errors.New("targetPath is not specified"))
	}
	dratini.TargetFilePath = *targetPath

	// overwrite if workerNum is specified by flags
	if *workerNum > 0 {
		dratini.ConfDratini.Core.WorkerNum = *workerNum
	}

	// overwrite if queueNum is specified by flags
	if *queueNum > 0 {
		dratini.ConfDratini.Core.QueueNum = *queueNum
	}

	// set logger
	accessLogger, accessLogReopener, err := dratini.InitLog(dratini.ConfDratini.Log.AccessLog, "info")
	if err != nil {
		dratini.LogSetupFatal(err)
	}
	errorLogger, errorLogReopener, err := dratini.InitLog(dratini.ConfDratini.Log.ErrorLog, dratini.ConfDratini.Log.Level)
	if err != nil {
		dratini.LogSetupFatal(err)
	}

	dratini.LogAccess = accessLogger
	dratini.LogError = errorLogger

	if !dratini.ConfDratini.Ios.Enabled && !dratini.ConfDratini.Android.Enabled {
		dratini.LogSetupFatal(fmt.Errorf("What do you want to do?"))
	}

	if dratini.ConfDratini.Ios.Enabled {
		dratini.CertificatePemIos.Cert, err = ioutil.ReadFile(dratini.ConfDratini.Ios.PemCertPath)
		if err != nil {
			dratini.LogSetupFatal(fmt.Errorf("A certification file for iOS is not found."))
		}

		dratini.CertificatePemIos.Key, err = ioutil.ReadFile(dratini.ConfDratini.Ios.PemKeyPath)
		if err != nil {
			dratini.LogSetupFatal(fmt.Errorf("A key file for iOS is not found."))
		}

	}

	if dratini.ConfDratini.Android.Enabled {
		if dratini.ConfDratini.Android.ApiKey == "" {
			dratini.LogSetupFatal(fmt.Errorf("APIKey for Android is empty."))
		}
	}

	sigHUPChan := make(chan os.Signal, 1)
	signal.Notify(sigHUPChan, syscall.SIGHUP)

	sighupHandler := func() {
		if err := accessLogReopener.Reopen(); err != nil {
			dratini.LogError.Warn(fmt.Sprintf("failed to reopen access log: %v", err))
		}
		if err := errorLogReopener.Reopen(); err != nil {
			dratini.LogError.Warn(fmt.Sprintf("failed to reopen error log: %v", err))
		}
	}

	go signalHandler(sigHUPChan, sighupHandler)

	if dratini.ConfDratini.Android.Enabled {
		if err := dratini.InitGCMClient(); err != nil {
			dratini.LogSetupFatal(fmt.Errorf("failed to init gcm/fcm client: %v", err))
		}
	}

	if dratini.ConfDratini.Ios.Enabled {
		if err := dratini.InitAPNSClient(); err != nil {
			dratini.LogSetupFatal(fmt.Errorf("failed to init http client for APNs: %v", err))
		}
	}
	dratini.InitStat()
	dratini.StartPushWorkers(dratini.ConfDratini.Core.WorkerNum, dratini.ConfDratini.Core.QueueNum)
	dratini.StartBatch()

	// Block until all job is kicked
	dratini.BootWg.Wait()
	// Block until all pusher worker job is done.
	dratini.PusherWg.Wait()

	dratini.LogError.Info("successfully shutdown")
}

func signalHandler(ch <-chan os.Signal, sighupFn func()) {
	for {
		select {
		case sig := <-ch:
			switch sig {
			case syscall.SIGHUP:
				sighupFn()
			}
		}
	}
}
