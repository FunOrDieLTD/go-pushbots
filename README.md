[![Build Status](https://travis-ci.org/FunOrDieLTD/go-pushbots.png)](https://travis-ci.org/FunOrDieLTD/go-pushbots)


## A pushbots API wrapper for go

A unofficial implementation around Pushbots service for sending APN in golang. 

### Documentation
http://godoc.org/github.com/FunOrDieLTD/go-pushbots

### Setup
1. Sign up with pushbots at https://pushbots.com

### Installation
`go get github.com/FunOrDieLTD/go-pushbots`

### Examples

#### Registering a IOS device

```go
	var appId string = "your app id"
	var secret string = "your secret"

	deviceToken := "Your device token"
	pushBots := pushbots.NewPushBots(appId, secret, true) // true enables debugging
	deviceToken = "Your device token"

	err := pushBots.RegisterDevice(deviceToken, pushbots.PlatformIos, "", "", []string{}, []string{}, "") // Registers an ios device using only it's token
	if err != nil {
		log.Fatal(err)
	}

```

#### Tagging a device
```go
	var appId string = "your app id"
	var secret string = "your secret"

	deviceToken := "Your device token"
	pushBots := pushbots.NewPushBots(appId, secret, true) // true enables debugging
	deviceToken = "Your device token"
	tag := "your_new_tag"

	err := pushBots.TagDevice(deviceToken, pushbots.PlatformIos, "", tag)

	if err != nil {
		log.Fatal(err)
	}

```

#### Sending a push to a single device
```go
	var appId string = "your app id"
	var secret string = "your secret"

	deviceToken := "Your device token"
	msg := "Your message"
	sound := "your_sound"
	badge := "0"
	payload := map[string]interface{}{
		"your": "custom data here",
	}

	pushBots := pushbots.NewPushBots(appId, secret, true) // true enables debugging

	err := pushBots.SendPushToDevice(pushbots.PlatformIos, deviceToken, msg, sound, badge, payload)

	if err != nil {
		log.Fatal(err)
	}

```
#### Broadcasting a push to all devices
```go
	var appId string = "your app id"
	var secret string = "your secret"

	msg := "Your message"
	sound := "your_sound"
	badge := "0"

	payload := map[string]interface{}{
		"your": "custom data here",
	}

	pushBots := pushbots.NewPushBots(appId, secret, true)

	err := pushBots.Broadcast(pushbots.PlatformAll, msg, sound, badge, payload)

	if err != nil {
		log.Fatal(err)
	}

```