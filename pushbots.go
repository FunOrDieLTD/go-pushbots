// pushbots contains a wrapper around pushbots REST api (https://pushbots.com/) to send Apple push notifications
// and android push notifications
// Please note that this is NOT an official library by pushbots
// All code Copyright David Pallinder and Fun or die LTD
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pushbots

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Constants for the different platforms supported
const (
	PlatformIos     = "0"
	PlatformAndroid = "1"
)

// Internal constants for keeping track of what endpoint to connect to
const (
	productionEndPoint = "https://api.pushbots.com/"
)

// Simply holds an endpoint and what http verb to use when connecting to that endpoint
type pushBotRequest struct {
	Endpoint string
	HttpVerb string
}

// Holds the appid and app secret for use in requests
type PushBots struct {
	AppId     string
	Secret    string
	Debug     bool
	endpoints map[string]pushBotRequest
}

type serverErrorResponse struct {
	Message interface{} `json:"message"`
}

// A struct to contain all arguments for a request
type apiRequest struct {
	Payload                 map[string]interface{} `json:"payload,omitempty"`
	Token                   string                 `json:"token,omitempty"`
	Platform                interface{}            `json:"platform,omitempty"` // Sometimes string sometimes []string
	Badge                   string                 `json:"badge,omitempty"`
	Sound                   string                 `json:"sound,omitempty"`
	Alias                   string                 `json:"alias,omitempty"`
	ExceptAlias             string                 `json:"except_alias,omitempty"`
	Tag                     string                 `json:"tag,omitempty"`
	Lat                     string                 `json:"lat,omitempty"`
	Lng                     string                 `json:"lng,omitempty"`
	Msg                     string                 `json:"msg,omitempty"`
	NotificationType        interface{}            `json:"active,omitempty"` // Sometimes string sometimes []string
	Stats                   string                 `json:"stats,omitempty"`
	Tags                    []string               `json:"tags,omitempty"`
	ExceptTags              []string               `json:"except_tags,omitempty"`
	ExceptNotificationTypes []string               `json:"except_active,omitempty"`
	BadgeCount              *int                   `json:"setbadgecount,omitempty"` //Hack to avoid badgecount being omitted if it's value is 0
}

// Create a new pushbots object
func NewPushBots(appId string, secret string, debug bool) PushBots {
	pushBots := PushBots{AppId: appId, Secret: secret, Debug: debug}
	return pushBots
}

// Override the base endpoint used to construct the API urls
func (pushBots *PushBots) ApplyEndpointOverride(endpointOverride string) {
	pushBots.initializeEndpoints(endpointOverride)
}

// Initializes all known endpoints
func (pushBots *PushBots) initializeEndpoints(endpointOverride string) {
	endpointBase := productionEndPoint

	if endpointOverride != "" {
		endpointBase = endpointOverride
	}

	pushBots.endpoints = map[string]pushBotRequest{
		"registerdevice":         pushBotRequest{Endpoint: endpointBase + "deviceToken", HttpVerb: "PUT"},
		"unregisterdevice":       pushBotRequest{Endpoint: endpointBase + "deviceToken/del", HttpVerb: "PUT"},
		"alias":                  pushBotRequest{Endpoint: endpointBase + "alias", HttpVerb: "PUT"},
		"tagdevice":              pushBotRequest{Endpoint: endpointBase + "tag", HttpVerb: "PUT"},
		"untagdevice":            pushBotRequest{Endpoint: endpointBase + "tag/del", HttpVerb: "PUT"},
		"geos":                   pushBotRequest{Endpoint: endpointBase + "geo", HttpVerb: "PUT"},
		"addnotificationtype":    pushBotRequest{Endpoint: endpointBase + "activate", HttpVerb: "PUT"},
		"removenotificationtype": pushBotRequest{Endpoint: endpointBase + "deactivate", HttpVerb: "PUT"},
		"broadcast":              pushBotRequest{Endpoint: endpointBase + "push/all", HttpVerb: "POST"},
		"pushone":                pushBotRequest{Endpoint: endpointBase + "push/one", HttpVerb: "POST"},
		"batch":                  pushBotRequest{Endpoint: endpointBase + "push/all", HttpVerb: "POST"},
		"badge":                  pushBotRequest{Endpoint: endpointBase + "badge", HttpVerb: "PUT"},
		"recordanalytics":        pushBotRequest{Endpoint: endpointBase + "stats", HttpVerb: "PUT"},
	}
}

// Register a device with PushBots
func (pushbots *PushBots) RegisterDevice(token, platform, lat, lng string, notificationTypes, tags []string, alias string) error {
	if err := checkForArgErrors(token, platform); err != nil {
		return err
	}

	args := apiRequest{
		Token:    token,
		Platform: platform,
		Lat:      lat,
		Lng:      lng,
		Tags:     tags,
		Alias:    alias,
	}

	if notificationTypes != nil && len(notificationTypes) > 0 {
		args.NotificationType = notificationTypes
	}

	return checkAndReturn(pushbots.sendToEndpoint("registerdevice", args))
}

// Unregister a device
func (pushbots *PushBots) UnregisterDevice(token, platform string) error {

	if err := checkForArgErrors(token, platform); err != nil {
		return err
	}

	args := apiRequest{
		Token:    token,
		Platform: platform,
	}

	return checkAndReturn(pushbots.sendToEndpoint("unregisterdevice", args))
}

// Add a tag to a device
func (pushbots *PushBots) TagDevice(token, platform, alias, tag string) error {

	if err := checkForArgErrorsWithAlias(token, platform, alias); err != nil {
		fmt.Println(err)
		return err
	}

	args := apiRequest{
		Token:    token,
		Alias:    alias,
		Platform: platform,
		Tag:      tag,
	}

	return checkAndReturn(pushbots.sendToEndpoint("tagdevice", args))
}

// Remove a tag from a device
func (pushbots *PushBots) UnTagDevice(token, platform, alias, tag string) error {
	if err := checkForArgErrorsWithAlias(token, platform, alias); err != nil {
		return err
	}

	args := apiRequest{
		Token:    token,
		Alias:    alias,
		Platform: platform,
		Tag:      tag,
	}

	return checkAndReturn(pushbots.sendToEndpoint("untagdevice", args))
}

// Add geo information to a device
func (pushbots *PushBots) Geo(token, platform, lat, lng string) error {
	if err := checkForArgErrors(token, platform); err != nil {
		return err
	}

	if lat == "" || lng == "" {
		return errors.New("Latitude/Longitude not specified")
	}

	args := apiRequest{
		Token:    token,
		Platform: platform,
		Lat:      lat,
		Lng:      lng,
	}

	return checkAndReturn(pushbots.sendToEndpoint("geos", args))
}

// Adds a notification type to a device
func (pushbots *PushBots) AddNotificationType(token, platform, alias, notificationType string) error {
	if err := checkForArgErrorsWithAlias(token, platform, alias); err != nil {
		return err
	}

	if notificationType == "" {
		return errors.New("No notification type specified")
	}

	args := apiRequest{
		Token:            token,
		Alias:            alias,
		Platform:         platform,
		NotificationType: notificationType,
	}

	return checkAndReturn(pushbots.sendToEndpoint("addnotificationtype", args))
}

// Removes a notification type from a device
func (pushbots *PushBots) RemoveNotificationType(token, platform, alias, notificationType string) error {
	if err := checkForArgErrorsWithAlias(token, platform, alias); err != nil {
		return err
	}

	if notificationType == "" {
		return errors.New("No notification type specified")
	}

	args := apiRequest{
		Token:            token,
		Alias:            alias,
		Platform:         platform,
		NotificationType: notificationType,
	}
	return checkAndReturn(pushbots.sendToEndpoint("removenotificationtype", args))
}

// Send a broadcast to multiple devices
func (pushbots *PushBots) Broadcast(platforms []string, msg, sound, badge string, payload map[string]interface{}) error {
	var supportsIos, supportsAndroid bool

	for _, platform := range platforms {
		if platform != PlatformIos && platform != PlatformAndroid {
			return errors.New("Platform neither IOS nor android")
		} else if platform == PlatformIos {
			supportsIos = true
		} else if platform == PlatformAndroid {
			supportsAndroid = true
		}
	}

	if supportsIos == false && supportsAndroid == false {
		return errors.New("Either android or ios must be specified as platforms")
	}

	if msg == "" {
		return errors.New("Message not specified")
	}

	if badge == "" {
		badge = "0"
	}

	if sound == "" {
		if supportsIos == true && supportsAndroid == false {
			sound = "default"
		} else {
			return errors.New("No sound specified")
		}
	}

	args := apiRequest{
		Platform: platforms,
		Msg:      msg,
		Badge:    badge,
		Sound:    sound,
		Payload:  payload,
	}

	return checkAndReturn(pushbots.sendToEndpoint("broadcast", args))
}

// Send a push to one device
func (pushbots *PushBots) SendPushToDevice(platform, token, msg, sound, badge string, payload map[string]interface{}) error {
	if err := checkForArgErrors(token, platform); err != nil {
		return err
	}

	if sound == "" {
		if platform == PlatformIos {
			sound = "default"
		} else {
			return errors.New("No sound specified")
		}
	}

	if msg == "" {
		return errors.New("No message specified")
	}

	if badge == "" {
		badge = "0"
	}

	args := apiRequest{
		Platform: platform,
		Token:    token,
		Msg:      msg,
		Sound:    sound,
		Badge:    badge,
		Payload:  payload,
	}
	return checkAndReturn(pushbots.sendToEndpoint("broadcast", args))
}

// Batch push notifications to matching devices
func (pushbots *PushBots) Batch(platform, msg, sound, badge string, tags, exceptTags, notificationTypes, exceptNotificationTypes []string,
	alias, exceptAlias string, payload map[string]interface{}) error {

	if platform != PlatformIos && platform != PlatformAndroid {
		return errors.New("Platform must be either PlatformIos or PlatformAndroid")
	}

	if msg == "" {
		return errors.New("No message specified")
	}

	if sound == "" && platform != PlatformIos {
		return errors.New("No sound specified")
	} else if sound == "" && platform == PlatformIos {
		sound = "default"
	}

	if badge == "" {
		badge = "0"
	}

	args := apiRequest{
		Payload:                 payload,
		Alias:                   alias,
		ExceptAlias:             exceptAlias,
		Platform:                platform,
		Msg:                     msg,
		Sound:                   sound,
		Badge:                   badge,
		Tags:                    tags,
		ExceptTags:              exceptTags,
		NotificationType:        notificationTypes,
		ExceptNotificationTypes: exceptNotificationTypes,
	}

	return checkAndReturn(pushbots.sendToEndpoint("broadcast", args))

}

// Set the badgecount for a device
func (pushbots *PushBots) Badge(token, platform string, badgeCount int) error {
	if err := checkForArgErrors(token, platform); err != nil {
		return err
	}

	args := apiRequest{
		Token:      token,
		Platform:   platform,
		BadgeCount: &badgeCount,
	}
	return checkAndReturn(pushbots.sendToEndpoint("badge", args))
}

// Record analytics for a device
func (pushbots *PushBots) RecordAnalytics(token, platform, stats string) error {
	if err := checkForArgErrors(token, platform); err != nil {
		return err
	}

	args := apiRequest{
		Token:    token,
		Platform: platform,
		Stats:    stats,
	}
	return checkAndReturn(pushbots.sendToEndpoint("recordanalytics", args))
}

// Prepare and send the request to the endpoint
func (pushbots *PushBots) sendToEndpoint(endpoint string, args apiRequest) ([]byte, error) {

	if pushbots.endpoints == nil {
		pushbots.initializeEndpoints("")
	}

	pushbotEndpoint, available := pushbots.endpoints[endpoint]

	if available == false {
		return []byte{}, errors.New("Could not find endpoint")
	}

	jsonPayload, err := json.Marshal(args)

	if err != nil {
		return []byte{}, err
	}

	if pushbots.Debug == true {

		fmt.Println("Sending JSON:", string(jsonPayload))
	}

	req, err := http.NewRequest(pushbotEndpoint.HttpVerb, pushbotEndpoint.Endpoint, strings.NewReader(string(jsonPayload)))

	if err != nil {
		return []byte{}, err
	}

	if pushbots.AppId == "" || pushbots.Secret == "" {
		return []byte{}, errors.New("Appid and/or secret key not set")
	}

	req.Header.Set("x-pushbots-appid", pushbots.AppId)
	req.Header.Set("x-pushbots-secret", pushbots.Secret)
	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)

	resp, err := client.Do(req)

	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if pushbots.Debug == true {
		fmt.Println("Response object:", resp)
		fmt.Println("Response content", string(body))
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return body, errors.New("Error response from server")
	}

	return body, err
}

// Checks for errors within arguments
func checkForArgErrors(token string, platform string) error {
	if token == "" {
		return errors.New("Token needs to be a device token")
	} else if platform != PlatformIos && platform != PlatformAndroid {
		return errors.New("Platform must be either PlatformIos or PlatformAndroid")
	}
	return nil
}

// Checks for errors when either a token or an alias is required
func checkForArgErrorsWithAlias(token string, platform, alias string) error {
	if token == "" && alias == "" {
		return errors.New("Either token or alias need to be set")
	} else if platform != PlatformIos && platform != PlatformAndroid {
		return errors.New("Platform must be either PlatformIos or PlatformAndroid")
	}
	return nil
}

func checkAndReturn(content []byte, err error) error {
	if err != nil {
		return err
	}

	if len(content) > 0 {
		serverErr := new(serverErrorResponse)
		err := json.Unmarshal(content, serverErr)

		if err != nil {
			return err
		}

		if serverErr.Message == "" {
			return errors.New("Got response from server but failed to parse")
		} else {
			switch serverErr.Message.(type) {
			case string:
				return errors.New(serverErr.Message.(string))
			case map[string]interface{}:
				return fmt.Errorf("A server error occurred: %s", string(content))
			default:
				return fmt.Errorf("Could not parse server message: %s", string(content))
			}
		}
	}
	return nil
}
