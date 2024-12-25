// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package agenttelemetry implements a component to generate Agent telemetry
package agenttelemetry

// team: agent-shared-components

// Component is the component type
type Component interface {
	// GetAsJSON returns the payload as a JSON string. Useful to be displayed in the CLI or added to a flare.
	GetAsJSON() ([]byte, error)

	// Sends telemetry payload.
	//    payloadType - should be registered in datadog-agent\comp\core\agenttelemetry\impl\config.go
	//    message - top level log message accompanying the payload
	//    payload - de-serializable into JSON payload
	Send(payloadType string, message string, payload []byte) error
}
