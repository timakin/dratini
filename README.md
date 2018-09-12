Dratini [![CircleCI](https://circleci.com/gh/timakin/dratini.svg?style=svg)](https://circleci.com/gh/timakin/dratini)
====

<img src="https://github.com/timakin/dratini/blob/master/dratini.png" alt="logo" align="right"/>

Dratini is a push notification handler works on a spot instance. Normally, push notification server is resident API, but, like daily notification job, most of time it stands by and costs meaninglessly.

You can reduce the cost if the handler works only at the moment. Dratini cannot serve request like normal push notification handler. However, it will send bulk push notifications in parallel with background workers based on goroutine.

# Installation

```
go get -u github.com/timakin/dratini
```

This project uses [Go Modules](https://github.com/golang/go/wiki/Modules), so if you've already installed Go (>= 1.11) and enables the flag `GO111MODULE=on`, you can automatically install the dependencies.

# How to use

To run dratini, you must provide configuration path via -c option (See [CONFIGURATION.md](/CONFIGURATION.md) about details), and path (or url) written notification request JSON via -t option (See [PUSH_REQUEST.md](/PUSH_REQUEST.md) about detail parameters).

```
$ cd $GOPATH/src/github.com/timakin/dratini
$ go build -o app
$ app -c /path/to/config.toml -t /path/or/url/to/notifications.json
```

## License

Licensed under the MIT License.


