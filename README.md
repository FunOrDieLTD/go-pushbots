[![Build Status](https://travis-ci.org/FunOrDieLTD/go-pushbots.png)](https://travis-ci.org/FunOrDieLTD/go-pushbots)


## A pushbots API wrapper for go

A unofficial implementation around Pushbots service for sending APN in golang. 

### Setup
1. Sign up with pushbots at https://pushbots.com

### Installation
`go get github.com/FunOrDieLTD/go-pushbots`

### Examples

#### Registering a IOS device

```go
package main

import (
	"github.com/FunOrDieLTD/go-pushbots"
	"log"
)

func main() {
	appId := "Your app id"
	secret := "Your secret"
	deviceToken := "Your device token"
	pushBots := pushbots.NewPushBots(appId, secret, true) // true enables debugging
	deviceToken = "Your device token"

	err := pushBots.RegisterDevice(deviceToken, pushbots.PlatformIos, "", "", []string{}, []string{}, "")
	if err != nil {
		log.Fatal(err)
	}

}

```
#### Sending a push to a single device
```go
	appId := "Your app id"
	secret := "Your secret"
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
	appId := "Your app id"
	secret := "Your secret"
	msg := "Your message"
	sound := "your_sound"
	badge := "0"

	payload := map[string]interface{}{
		"your": "custom data here",
	}

	platforms := []string{PlatformIos}

	shouldEqual := map[string]interface{}{
		"platform": stringSliceToInterfaceSlice(platforms),
		"msg":      msg,
		"badge":    badge,
		"sound":    sound,
		"payload":  payload,
	}

	testServer := httptest.NewServer(http.HandlerFunc(testHandler(t, shouldEqual)))
	defer testServer.Close()

	pushBots := NewPushBots(appId, secret, false)

	pushBots.ApplyEndpointOverride(testServer.URL + "/")

	err := pushBots.Broadcast(platforms, msg, sound, badge, payload)

	if err != nil {
		t.Fatal(err)
	}
```

