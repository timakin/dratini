package dratini

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/timakin/dratini/dratini"
)

const (
	DefaultPidPermission = 0644
)

func main() {
	versionPrinted := flag.Bool("v", false, "dratini version")
	confPath := flag.String("c", "", "configuration file path for dratini")
	listenPort := flag.String("p", "", "port number or unix socket path")
	workerNum := flag.Int64("w", 0, "number of workers for push notification")
	queueNum := flag.Int64("q", 0, "size of internal queue for push notification")
	flag.Parse()

	if *versionPrinted {
		dratini.PrintVersion()
		return
	}

	// set default parameters
	dratini.ConfDratini = dratini.BuildDefaultConf()

	// load configuration
	conf, err := dratini.LoadConf(dratini.ConfDratini, *confPath)
	if err != nil {
		dratini.LogSetupFatal(err)
	}
	dratini.ConfDratini = conf

	// overwrite if port is specified by flags
	if *listenPort != "" {
		dratini.ConfDratini.Core.Port = *listenPort
	}

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

	if len(conf.Core.Pid) > 0 {
		if _, err := os.Stat(filepath.Dir(conf.Core.Pid)); os.IsNotExist(err) {
			dratini.LogSetupFatal(fmt.Errorf("directory for pid file is not exist: %v", err))
		} else if err := ioutil.WriteFile(conf.Core.Pid, []byte(strconv.Itoa(os.Getpid())), DefaultPidPermission); err != nil {
			dratini.LogSetupFatal(fmt.Errorf("failed to create a pid file: %v", err))
		}
	}

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

	// Start PushJob
	// <- job started

	//mux := http.NewServeMux()
	//dratini.RegisterHandlers(mux)

	//server := &http.Server{
	//	Handler: mux,
	//}
	//go func() {
	//	dratini.LogError.Info("start server")
	//	if err := dratini.RunServer(server, &dratini.ConfDratini); err != nil {
	//		dratini.LogError.Info(fmt.Sprintf("failed to serve: %s", err))
	//	}
	//}()

	// Graceful shutdown (kicked by SIGTERM).
	//
	// First, it shutdowns server and stops accepting new requests.
	// Then wait until all remaining queues in buffer are flushed.
	//sigTERMChan := make(chan os.Signal, 1)
	//signal.Notify(sigTERMChan, syscall.SIGTERM)
	//
	//<-sigTERMChan
	//dratini.LogError.Info("shutdown server")
	//timeout := time.Duration(conf.Core.ShutdownTimeout) * time.Second
	//ctx, cancel := context.WithTimeout(context.Background(), timeout)
	//defer cancel()
	//if err := server.Shutdown(ctx); err != nil {
	//	dratini.LogError.Error(fmt.Sprintf("failed to shutdown server: %v", err))
	//}

	// Start a goroutine to log number of job queue.
	go func() {
		for {
			queue := len(dratini.QueueNotification)
			if queue == 0 {
				break
			}

			dratini.LogError.Info(fmt.Sprintf("wait until queue is empty. Current queue len: %d", queue))
			time.Sleep(1 * time.Second)
		}
	}()

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
