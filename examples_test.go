package pushbots_test

import (
	"github.com/FunOrDieLTD/go-pushbots"
	"log"
)

func ExamplePushBots_TagDevice() {
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

}

func ExamplePushBots_RegisterDevice() {
	var appId string = "your app id"
	var secret string = "your secret"

	deviceToken := "Your device token"
	pushBots := pushbots.NewPushBots(appId, secret, true) // true enables debugging
	deviceToken = "Your device token"

	err := pushBots.RegisterDevice(deviceToken, pushbots.PlatformIos, "", "", []string{}, []string{}, "") // Registers an ios device using only it's token
	if err != nil {
		log.Fatal(err)
	}
}

func ExamplePushBots_SendPushToDevice() {
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
}

func ExamplePushBots_Broadcast() {
	var appId string = "your app id"
	var secret string = "your secret"

	msg := "Your message"
	sound := "your_sound"
	badge := "0"

	payload := map[string]interface{}{
		"your": "custom data here",
	}

	pushBots := pushbots.NewPushBots(appId, secret, true)

	err := pushBots.Broadcast(pushbots.PlatformIos, msg, sound, badge, payload)

	if err != nil {
		log.Fatal(err)
	}
}
