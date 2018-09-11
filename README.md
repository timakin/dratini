Dratini
====

<img src="https://github.com/timakin/dratini/blob/master/dratini.png" alt="logo" align="right"/>

Dratini is a push notification handler works on a spot instance. Normally, push notification server is resident API, but, like daily notification job, most of time it stands by and costs meaninglessly.

You can reduce the cost if the handler works only at the moment. Dratini cannot serve request like normal push notification handler. However, it will send bulk push notifications in parallel with background workers based on goroutine.

# Dependency



