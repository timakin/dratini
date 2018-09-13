package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/timakin/dratini/pusher"
)

func main() {
	confPath := flag.String("c", "", "configuration file path for pusher")
	targetPath := flag.String("t", "", "configuration file path for pusher push target")
	workerNum := flag.Int64("w", 0, "number of workers for push notification")
	queueNum := flag.Int64("q", 0, "size of internal queue for push notification")
	flag.Parse()

	// set default parameters
	pusher.ConfDratini = pusher.BuildDefaultConf()

	// load configuration
	conf, err := pusher.LoadConf(pusher.ConfDratini, *confPath)
	if err != nil {
		pusher.LogSetupFatal(err)
	}
	pusher.ConfDratini = conf

	// exit if push target is not specified by flags
	if *targetPath == "" {
		pusher.LogSetupFatal(errors.New("targetPath is not specified"))
	}
	pusher.TargetFilePath = *targetPath

	// overwrite if workerNum is specified by flags
	if *workerNum > 0 {
		pusher.ConfDratini.Core.WorkerNum = *workerNum
	}

	// overwrite if queueNum is specified by flags
	if *queueNum > 0 {
		pusher.ConfDratini.Core.QueueNum = *queueNum
	}

	// set logger
	accessLogger, accessLogReopener, err := pusher.InitLog(pusher.ConfDratini.Log.AccessLog, "info")
	if err != nil {
		pusher.LogSetupFatal(err)
	}
	errorLogger, errorLogReopener, err := pusher.InitLog(pusher.ConfDratini.Log.ErrorLog, pusher.ConfDratini.Log.Level)
	if err != nil {
		pusher.LogSetupFatal(err)
	}

	pusher.LogAccess = accessLogger
	pusher.LogError = errorLogger

	if !pusher.ConfDratini.Ios.Enabled && !pusher.ConfDratini.Android.Enabled {
		pusher.LogSetupFatal(fmt.Errorf("What do you want to do?"))
	}

	if pusher.ConfDratini.Ios.Enabled {
		pusher.CertificatePemIos.Cert, err = ioutil.ReadFile(pusher.ConfDratini.Ios.PemCertPath)
		if err != nil {
			pusher.LogSetupFatal(fmt.Errorf("A certification file for iOS is not found."))
		}

		pusher.CertificatePemIos.Key, err = ioutil.ReadFile(pusher.ConfDratini.Ios.PemKeyPath)
		if err != nil {
			pusher.LogSetupFatal(fmt.Errorf("A key file for iOS is not found."))
		}

	}

	if pusher.ConfDratini.Android.Enabled {
		if pusher.ConfDratini.Android.ApiKey == "" {
			pusher.LogSetupFatal(fmt.Errorf("APIKey for Android is empty."))
		}
	}

	sigHUPChan := make(chan os.Signal, 1)
	signal.Notify(sigHUPChan, syscall.SIGHUP)

	sighupHandler := func() {
		if err := accessLogReopener.Reopen(); err != nil {
			pusher.LogError.Warn(fmt.Sprintf("failed to reopen access log: %v", err))
		}
		if err := errorLogReopener.Reopen(); err != nil {
			pusher.LogError.Warn(fmt.Sprintf("failed to reopen error log: %v", err))
		}
	}

	go signalHandler(sigHUPChan, sighupHandler)

	if pusher.ConfDratini.Android.Enabled {
		if err := pusher.InitGCMClient(); err != nil {
			pusher.LogSetupFatal(fmt.Errorf("failed to init gcm/fcm client: %v", err))
		}
	}

	if pusher.ConfDratini.Ios.Enabled {
		if err := pusher.InitAPNSClient(); err != nil {
			pusher.LogSetupFatal(fmt.Errorf("failed to init http client for APNs: %v", err))
		}
	}
	pusher.InitStat()
	pusher.StartPushWorkers(pusher.ConfDratini.Core.WorkerNum, pusher.ConfDratini.Core.QueueNum)
	pusher.StartBatch()

	// Block until all job is kicked
	pusher.BootWg.Wait()
	// Block until all pusher worker job is done.
	pusher.PusherWg.Wait()

	pusher.LogError.Info("successfully shutdown")
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
