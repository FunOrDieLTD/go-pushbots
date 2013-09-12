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

	if err := pushBots.RegisterDevice(deviceToken, pushbots.PlatformIos, "", "", []string{}, []string{}, ""); err != nil {
		log.Fatal(err)
	}

}

```
#### Sending a push to a single device


#### Broadcasting a push to all devices
