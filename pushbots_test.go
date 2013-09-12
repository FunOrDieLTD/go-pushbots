// Copyright 2012 David Pallinder, Fun or die ltd. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pushbots

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const (
	token             = "token"
	secret            = "secret"
	appId             = "appId"
	lat               = "59.3333333"
	lng               = "18.05"
	alias             = "alias"
	tag1              = "tag1"
	tag2              = "tag2"
	notificationType1 = "notificationtype1"
	notificationType2 = "notificationtype2"
	badge             = "0"
	sound             = "sound"
	msg               = "msg"
)

// Converts a string slice to a slice filled with interfaces
// needed because unmarshaling json in go arrays get treated like
// interfaces
func stringSliceToInterfaceSlice(original []string) []interface{} {
	interfaceSlice := make([]interface{}, len(original))

	for i := range original {
		interfaceSlice[i] = original[i]
	}

	return interfaceSlice
}

func isHeaderOK(header http.Header) bool {
	if header["X-Pushbots-Appid"][0] != appId {
		return false
	}

	if header["X-Pushbots-Secret"][0] != secret {
		return false
	}

	if header["Content-Type"][0] != "application/json" {
		return false
	}

	return true
}

func testHandler(t *testing.T, shouldEqual map[string]interface{}) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		bodyUnmarshaled := make(map[string]interface{})

		if isHeaderOK(r.Header) == false {
			t.Fatal("Malformed header")
		}

		defer r.Body.Close()

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			t.Fatal(err)
		}

		err = json.Unmarshal(body, &bodyUnmarshaled)

		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(shouldEqual, bodyUnmarshaled) {
			fmt.Println("Should equal:", shouldEqual)
			fmt.Println("bodyUnmarshaled", bodyUnmarshaled)
			t.Fatal("Body was not correct")
		}
	}
}

func TestRegisterDevice(t *testing.T) {
	tags := []string{tag1, tag2}
	notificationTypes := []string{notificationType1, notificationType2}

	shouldEqual := map[string]interface{}{
		"token":    token,
		"lat":      lat,
		"lng":      lng,
		"alias":    alias,
		"tags":     stringSliceToInterfaceSlice(tags),
		"active":   stringSliceToInterfaceSlice(notificationTypes),
		"platform": PlatformIos,
	}

	testServer := httptest.NewServer(http.HandlerFunc(testHandler(t, shouldEqual)))
	defer testServer.Close()

	pushBots := NewPushBots(appId, secret, false)

	pushBots.ApplyEndpointOverride(testServer.URL + "/")

	err := pushBots.RegisterDevice(token, PlatformIos, lat, lng, notificationTypes, tags, alias)

	if err != nil {
		t.Fatal(err)
	}
}

func TestUnregisteringDevice(t *testing.T) {
	shouldEqual := map[string]interface{}{
		"token":    token,
		"platform": PlatformIos,
	}

	testServer := httptest.NewServer(http.HandlerFunc(testHandler(t, shouldEqual)))
	defer testServer.Close()

	pushBots := NewPushBots(appId, secret, false)

	pushBots.ApplyEndpointOverride(testServer.URL + "/")

	err := pushBots.UnregisterDevice(token, PlatformIos)

	if err != nil {
		t.Fatal(err)
	}
}

func TestTagDevice(t *testing.T) {
	shouldEqual := map[string]interface{}{
		"token":    token,
		"platform": PlatformIos,
		"alias":    alias,
		"tag":      tag1,
	}

	testServer := httptest.NewServer(http.HandlerFunc(testHandler(t, shouldEqual)))
	defer testServer.Close()

	pushBots := NewPushBots(appId, secret, false)

	pushBots.ApplyEndpointOverride(testServer.URL + "/")

	err := pushBots.TagDevice(token, PlatformIos, alias, tag1)

	if err != nil {
		t.Fatal(err)
	}
}

func TestUntagDevice(t *testing.T) {
	shouldEqual := map[string]interface{}{
		"token":    token,
		"platform": PlatformIos,
		"alias":    alias,
		"tag":      tag1,
	}

	testServer := httptest.NewServer(http.HandlerFunc(testHandler(t, shouldEqual)))
	defer testServer.Close()

	pushBots := NewPushBots(appId, secret, false)

	pushBots.ApplyEndpointOverride(testServer.URL + "/")

	err := pushBots.UnTagDevice(token, PlatformIos, alias, tag1)

	if err != nil {
		t.Fatal(err)
	}
}

func TestGeo(t *testing.T) {
	shouldEqual := map[string]interface{}{
		"token":    token,
		"platform": PlatformIos,
		"lat":      lat,
		"lng":      lng,
	}

	testServer := httptest.NewServer(http.HandlerFunc(testHandler(t, shouldEqual)))
	defer testServer.Close()

	pushBots := NewPushBots(appId, secret, false)

	pushBots.ApplyEndpointOverride(testServer.URL + "/")

	err := pushBots.Geo(token, PlatformIos, lat, lng)

	if err != nil {
		t.Fatal(err)
	}
}

func TestAddNotificationType(t *testing.T) {
	shouldEqual := map[string]interface{}{
		"token":    token,
		"platform": PlatformIos,
		"active":   notificationType1,
		"alias":    alias,
	}

	testServer := httptest.NewServer(http.HandlerFunc(testHandler(t, shouldEqual)))
	defer testServer.Close()

	pushBots := NewPushBots(appId, secret, false)

	pushBots.ApplyEndpointOverride(testServer.URL + "/")

	err := pushBots.AddNotificationType(token, PlatformIos, alias, notificationType1)

	if err != nil {
		t.Fatal(err)
	}
}

func TestRemoveNotificationType(t *testing.T) {
	shouldEqual := map[string]interface{}{
		"token":    token,
		"platform": PlatformIos,
		"active":   notificationType1,
		"alias":    alias,
	}

	testServer := httptest.NewServer(http.HandlerFunc(testHandler(t, shouldEqual)))
	defer testServer.Close()

	pushBots := NewPushBots(appId, secret, false)

	pushBots.ApplyEndpointOverride(testServer.URL + "/")

	err := pushBots.RemoveNotificationType(token, PlatformIos, alias, notificationType1)

	if err != nil {
		t.Fatal(err)
	}
}

func TestBroadcast(t *testing.T) {
	payload := map[string]interface{}{"a": "b"}
	platforms := []string{PlatformIos, PlatformAndroid}

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

	err := pushBots.Broadcast(PlatformAll, msg, sound, badge, payload)

	if err != nil {
		t.Fatal(err)
	}
}

func TestSendPushToDevice(t *testing.T) {
	payload := map[string]interface{}{"a": "b"}

	shouldEqual := map[string]interface{}{
		"platform": PlatformIos,
		"token":    token,
		"msg":      msg,
		"badge":    badge,
		"sound":    sound,
		"payload":  payload,
	}

	testServer := httptest.NewServer(http.HandlerFunc(testHandler(t, shouldEqual)))
	defer testServer.Close()

	pushBots := NewPushBots(appId, secret, false)

	pushBots.ApplyEndpointOverride(testServer.URL + "/")

	err := pushBots.SendPushToDevice(PlatformIos, token, msg, sound, badge, payload)

	if err != nil {
		t.Fatal(err)
	}
}

func TestBatch(t *testing.T) {
	payload := map[string]interface{}{"a": "b"}
	tags := []string{tag1}
	exceptTags := []string{tag2}
	notificationTypes := []string{notificationType1}
	exceptNotificationTypes := []string{notificationType2}
	exceptAlias := "exceptAlias"

	shouldEqual := map[string]interface{}{
		"platform":      PlatformIos,
		"msg":           msg,
		"badge":         badge,
		"sound":         sound,
		"payload":       payload,
		"tags":          stringSliceToInterfaceSlice(tags),
		"except_tags":   stringSliceToInterfaceSlice(exceptTags),
		"active":        stringSliceToInterfaceSlice(notificationTypes),
		"except_active": stringSliceToInterfaceSlice(exceptNotificationTypes),
		"alias":         alias,
		"except_alias":  exceptAlias,
	}

	testServer := httptest.NewServer(http.HandlerFunc(testHandler(t, shouldEqual)))
	defer testServer.Close()

	pushBots := NewPushBots(appId, secret, false)

	pushBots.ApplyEndpointOverride(testServer.URL + "/")

	err := pushBots.Batch(PlatformIos, msg, sound, badge, tags, exceptTags,
		notificationTypes, exceptNotificationTypes, alias, exceptAlias, payload)

	if err != nil {
		t.Fatal(err)
	}
}

func TestBadge(t *testing.T) {
	badgeCount := 1
	shouldEqual := map[string]interface{}{
		"token":         token,
		"platform":      PlatformIos,
		"setbadgecount": float64(badgeCount),
	}

	testServer := httptest.NewServer(http.HandlerFunc(testHandler(t, shouldEqual)))
	defer testServer.Close()
	pushBots := NewPushBots(appId, secret, false)

	pushBots.ApplyEndpointOverride(testServer.URL + "/")

	err := pushBots.Badge(token, PlatformIos, badgeCount)

	if err != nil {
		t.Fatal(err)
	}
}

func TestRecordAnalytics(t *testing.T) {
	stats := "o"
	shouldEqual := map[string]interface{}{
		"token":    token,
		"platform": PlatformIos,
		"stats":    stats,
	}

	testServer := httptest.NewServer(http.HandlerFunc(testHandler(t, shouldEqual)))
	defer testServer.Close()
	pushBots := NewPushBots(appId, secret, false)

	pushBots.ApplyEndpointOverride(testServer.URL + "/")

	err := pushBots.RecordAnalytics(token, PlatformIos, stats)

	if err != nil {
		t.Fatal(err)
	}
}

func TestErrorHandling(t *testing.T) {
	stats := "o"

	testServer := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(resp, `{"message":"An error"}`)
	}))
	defer testServer.Close()
	pushBots := NewPushBots(appId, secret, false)

	pushBots.ApplyEndpointOverride(testServer.URL + "/")

	err := pushBots.RecordAnalytics(token, PlatformIos, stats)

	if err == nil {
		t.Fatal("No error was returned")
	}
}

func TestCheckForArgErrors(t *testing.T) {
	t.Parallel()
	if err := checkForArgErrors("", PlatformIos); err == nil {
		t.Fatal("Should have been false when not specifying token")
	}

	if err := checkForArgErrors(token, "8"); err == nil {
		t.Fatal("Should have been false when specifying wrong platform")
	}

	if err := checkForArgErrors(token, PlatformAndroid); err != nil {
		t.Fatal("Should have been true when passing correct arguments")
	}
}

func TestCheckForArgErrorsWithAlias(t *testing.T) {
	t.Parallel()
	if err := checkForArgErrorsWithAlias("", PlatformIos, ""); err == nil {
		t.Fatal("Should have been false when not specifying token nor alias")
	}

	if err := checkForArgErrorsWithAlias(token, "8", alias); err == nil {
		t.Fatal("Should have been false when specifying wrong platform")
	}

	if err := checkForArgErrorsWithAlias(token, PlatformAndroid, alias); err != nil {
		t.Fatal("Should have been true when passing correct arguments")
	}
}

func TestGeneratePlatform(t *testing.T) {
	t.Parallel()
	platforms, err := generatePlatform(PlatformAndroid, false)

	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(platforms).String() != "string" {
		t.Fatal("Failed to return string")
	}

	if platforms.(string) != "1" {
		t.Fatal("Failed to return correct platform")
	}

	platforms, err = generatePlatform(PlatformIos, false)

	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(platforms).String() != "string" {
		t.Fatal("Failed to return string")
	}

	if platforms.(string) != "0" {
		t.Fatal("Failed to return string")
	}

	platforms, err = generatePlatform(PlatformIos, true)

	if reflect.TypeOf(platforms).String() != "[]string" {
		t.Fatal("Failed to return array")
	}

	if len(platforms.([]string)) == 0 {
		t.Fatal("Array wasnt properly populated")
	}

	platform := platforms.([]string)[0]

	if platform != PlatformIos {
		t.Fatal("Failed to return correct platform")

	}

	platforms, err = generatePlatform(PlatformAll, false)

	if err == nil {
		t.Fatal("Failed to generate array when issuing platformall and no array")
	}

	platforms, err = generatePlatform(PlatformAll, true)

	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(platforms).String() != "[]string" {
		t.Fatal("Array not returned")
	}

	if len(platforms.([]string)) != 2 {
		t.Fatal("Array wasnt properly populated")
	}

	for _, v := range platforms.([]string) {
		if v != PlatformIos && v != PlatformAndroid {
			t.Fatal("Wrong type in array")
		}
	}
}
