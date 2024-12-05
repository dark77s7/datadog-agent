// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build test

// //nolint:revive // TODO(PLINT) Fix revive linter
package wlan

import (
	"testing"

	"github.com/DataDog/datadog-agent/comp/core/autodiscovery/integration"
	"github.com/DataDog/datadog-agent/pkg/aggregator/mocksender"
	"github.com/stretchr/testify/mock"
)

func TestWLANOK(t *testing.T) {
	// setup mocks
	getWifiInfo = func() (WiFiInfo, error) {
		return WiFiInfo{
			Rssi:         10,
			Ssid:         "test-ssid",
			Bssid:        "test-bssid",
			Channel:      1,
			Noise:        20,
			TransmitRate: 4.0,
			SecurityType: "WPA/WPA2 Personal",
		}, nil
	}
	setupLocationAccess = func() {
	}
	defer func() {
		getWifiInfo = GetWiFiInfo
		setupLocationAccess = SetupLocationAccess
	}()

	wlanCheck := new(WLANCheck)

	senderManager := mocksender.CreateDefaultDemultiplexer()
	wlanCheck.Configure(senderManager, integration.FakeConfigHash, nil, nil, "test")

	mockSender := mocksender.NewMockSenderWithSenderManager(wlanCheck.ID(), senderManager)

	expectedTags := []string{"ssid:test-ssid", "bssid:test-bssid", "security_type:wpa/wpa2_personal"}

	mockSender.On("Gauge", "wlan.rssi", 10.0, mock.Anything, expectedTags).Return().Times(1)
	mockSender.On("Gauge", "wlan.noise", 20.0, mock.Anything, expectedTags).Return().Times(1)
	mockSender.On("Gauge", "wlan.transmit_rate", 4.0, mock.Anything, expectedTags).Return().Times(1)
	mockSender.On("Count", "wlan.channel_swap_events", 0.0, mock.Anything, expectedTags).Return().Times(1)

	mockSender.On("Commit").Return().Times(1)
	wlanCheck.Run()

	mockSender.AssertExpectations(t)
	mockSender.AssertNumberOfCalls(t, "Gauge", 3)
	mockSender.AssertNumberOfCalls(t, "Count", 1)
	mockSender.AssertNumberOfCalls(t, "Commit", 1)
}

func TestWLANEmptySSIDandBSSID(t *testing.T) {
	// setup mocks
	getWifiInfo = func() (WiFiInfo, error) {
		return WiFiInfo{
			Rssi:         10,
			Ssid:         "",
			Bssid:        "",
			Channel:      1,
			Noise:        20,
			TransmitRate: 4.0,
			SecurityType: "WPA/WPA2 Personal",
		}, nil
	}
	setupLocationAccess = func() {
	}
	defer func() {
		getWifiInfo = GetWiFiInfo
		setupLocationAccess = SetupLocationAccess
	}()

	wlanCheck := new(WLANCheck)

	senderManager := mocksender.CreateDefaultDemultiplexer()
	wlanCheck.Configure(senderManager, integration.FakeConfigHash, nil, nil, "test")

	mockSender := mocksender.NewMockSenderWithSenderManager(wlanCheck.ID(), senderManager)

	expectedTags := []string{"ssid:unknown", "bssid:unknown", "security_type:wpa/wpa2_personal"}

	mockSender.On("Gauge", "wlan.rssi", 10.0, mock.Anything, expectedTags).Return().Times(1)
	mockSender.On("Gauge", "wlan.noise", 20.0, mock.Anything, expectedTags).Return().Times(1)
	mockSender.On("Gauge", "wlan.transmit_rate", 4.0, mock.Anything, expectedTags).Return().Times(1)
	mockSender.On("Count", "wlan.channel_swap_events", 0.0, mock.Anything, expectedTags).Return().Times(1)

	mockSender.On("Commit").Return().Times(1)
	wlanCheck.Run()

	mockSender.AssertExpectations(t)
	mockSender.AssertNumberOfCalls(t, "Gauge", 3)
	mockSender.AssertNumberOfCalls(t, "Count", 1)
	mockSender.AssertNumberOfCalls(t, "Commit", 1)
}

func TestWLANChannelSwapEvents(t *testing.T) {
	// setup mocks
	getWifiInfo = func() (WiFiInfo, error) {
		return WiFiInfo{
			Rssi:         10,
			Ssid:         "",
			Bssid:        "",
			Channel:      1,
			Noise:        20,
			TransmitRate: 4.0,
			SecurityType: "WPA/WPA2 Personal",
		}, nil
	}
	setupLocationAccess = func() {
	}
	defer func() {
		getWifiInfo = GetWiFiInfo
		setupLocationAccess = SetupLocationAccess
	}()

	wlanCheck := new(WLANCheck)

	senderManager := mocksender.CreateDefaultDemultiplexer()
	wlanCheck.Configure(senderManager, integration.FakeConfigHash, nil, nil, "test")

	mockSender := mocksender.NewMockSenderWithSenderManager(wlanCheck.ID(), senderManager)

	expectedTags := []string{"ssid:unknown", "bssid:unknown", "security_type:wpa/wpa2_personal"}

	mockSender.On("Gauge", "wlan.rssi", 10.0, mock.Anything, expectedTags).Return().Times(1)
	mockSender.On("Gauge", "wlan.noise", 20.0, mock.Anything, expectedTags).Return().Times(1)
	mockSender.On("Gauge", "wlan.transmit_rate", 4.0, mock.Anything, expectedTags).Return().Times(1)
	mockSender.On("Count", "wlan.channel_swap_events", 0.0, mock.Anything, expectedTags).Return().Times(1)
	mockSender.On("Commit").Return().Times(1)

	wlanCheck.Run()

	mockSender.AssertExpectations(t)
	mockSender.AssertNumberOfCalls(t, "Gauge", 3)
	mockSender.AssertNumberOfCalls(t, "Count", 1)
	mockSender.AssertNumberOfCalls(t, "Commit", 1)

	// change channel number from 1 to 2
	getWifiInfo = func() (WiFiInfo, error) {
		return WiFiInfo{
			Rssi:         10,
			Ssid:         "",
			Bssid:        "",
			Channel:      2,
			Noise:        20,
			TransmitRate: 4.0,
			SecurityType: "WPA/WPA2 Personal",
		}, nil
	}

	mockSender = mocksender.NewMockSenderWithSenderManager(wlanCheck.ID(), senderManager)

	mockSender.On("Gauge", "wlan.rssi", 10.0, mock.Anything, expectedTags).Return().Times(1)
	mockSender.On("Gauge", "wlan.noise", 20.0, mock.Anything, expectedTags).Return().Times(1)
	mockSender.On("Gauge", "wlan.transmit_rate", 4.0, mock.Anything, expectedTags).Return().Times(1)
	mockSender.On("Count", "wlan.channel_swap_events", 1.0, mock.Anything, expectedTags).Return().Times(1)
	mockSender.On("Commit").Return().Times(1)

	wlanCheck.Run()

	mockSender.AssertExpectations(t)
	mockSender.AssertNumberOfCalls(t, "Gauge", 3)
	mockSender.AssertNumberOfCalls(t, "Count", 1)
	mockSender.AssertNumberOfCalls(t, "Commit", 1)

	// change channel number from 2 to 1
	getWifiInfo = func() (WiFiInfo, error) {
		return WiFiInfo{
			Rssi:         10,
			Ssid:         "",
			Bssid:        "",
			Channel:      1,
			Noise:        20,
			TransmitRate: 4.0,
			SecurityType: "WPA/WPA2 Personal",
		}, nil
	}

	mockSender = mocksender.NewMockSenderWithSenderManager(wlanCheck.ID(), senderManager)

	mockSender.On("Gauge", "wlan.rssi", 10.0, mock.Anything, expectedTags).Return().Times(1)
	mockSender.On("Gauge", "wlan.noise", 20.0, mock.Anything, expectedTags).Return().Times(1)
	mockSender.On("Gauge", "wlan.transmit_rate", 4.0, mock.Anything, expectedTags).Return().Times(1)
	mockSender.On("Count", "wlan.channel_swap_events", 1.0, mock.Anything, expectedTags).Return().Times(1)
	mockSender.On("Commit").Return().Times(1)

	wlanCheck.Run()

	mockSender.AssertExpectations(t)
	mockSender.AssertNumberOfCalls(t, "Gauge", 3)
	mockSender.AssertNumberOfCalls(t, "Count", 1)
	mockSender.AssertNumberOfCalls(t, "Commit", 1)
}
