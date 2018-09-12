# Spec of push request

You must specify the JSON file what is used to let API to see where to send bulk-pushes.

The JSON below is the example.

```json
{
    "notifications" : [
        {
            "token" : ["xxx"],
            "platform" : 1,
            "message" : "Hello, iOS!",
            "title": "Greeting",
            "subtitle": "greeting",
            "badge" : 1,
            "category": "category1",
            "sound" : "default",
            "content_available" : false,
            "mutable_content" : false,
            "expiry" : 10,
            "extend" : [{ "key": "url", "val": "..." }, { "key": "intent", "val": "..." }]
        },
        {
            "token" : ["yyy"],
            "platform" : 2,
            "message" : "Hello, Android!",
            "collapse_key" : "update",
            "delay_while_idle" : true,
            "time_to_live" : 10
        }
    ]
}
```

The entity must has the `notifications` array. There is the parameter table for each notification below.

|name             |type        |description                              |required|default|note                                      |
|-----------------|------------|-----------------------------------------|--------|-------|------------------------------------------|
|token            |string array|device tokens                            |o       |       |                                          |
|platform         |int         |platform(iOS,Android)                    |o       |       |1=iOS, 2=Android                          |
|message          |string      |message for notification                 |o       |       |                                          |
|title            |string      |title for notification                   |-       |       |only iOS                                  |
|subtitle         |string      |subtitle for notification                |-       |       |only iOS                                  |
|badge            |int         |badge count                              |-       |0      |only iOS                                  |
|category         |string      |unnotification category                  |-       |       |only iOS                                  |
|sound            |string      |sound type                               |-       |       |only iOS                                  |
|expiry           |int         |expiration for notification              |-       |0      |only iOS.                                 |
|content_available|bool        |indicate that new content is available   |-       |false  |only iOS.                                 |
|mutable_content  |bool        |enable Notification Service app extension|-       |false  |only iOS(10.0+).                          |
|collapse_key     |string      |the key for collapsing notifications     |-       |       |only Android                              |
|delay_while_idle |bool        |the flag for device idling               |-       |false  |only Android                              |
|time_to_live     |int         |expiration of message kept on GCM storage|-       |0      |only Android                              |
|extend           |string array|extensible partition                     |-       |       |                                          |
|identifier       |string      |notification identifier                  |-       |       |an optional value to identify notification|
