// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

//go:build linux && arm64

// Package modules is all the module definitions for system-probe
package modules

import (
	"time"

	"github.com/DataDog/datadog-agent/cmd/system-probe/api/module"
)

// All System Probe modules should register their factories here
var All = []module.Factory{
	EBPFProbe,
	NetworkTracer,
	TCPQueueLength,
	OOMKillProbe,
	// there is a dependency from EventMonitor -> NetworkTracer
	// so EventMonitor has to follow NetworkTracer
	EventMonitor,
	Process,
	DynamicInstrumentation,
	LanguageDetectionModule,
	ComplianceModule,
	Pinger,
	Traceroute,
	DiscoveryModule,
	GPUMonitoring, // GPU monitoring needs to be initialized afer EventMonitor, so that we have the event consumer ready
}

func inactivityEventLog(_ time.Duration) {

}
