// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

//go:build test

// Package fx provides the fx module for the serializer/compression component
package fx

import (
	compressionnoop "github.com/DataDog/datadog-agent/comp/serializer/compression/impl-noop"
	"go.uber.org/fx"

	compression "github.com/DataDog/datadog-agent/comp/serializer/compression/def"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
)

// MockModule defines the fx options for the mock component.
func MockModule() fxutil.Module {
	return fxutil.Component(
		fx.Provide(NewMockCompressor),
	)
}

// NewMockCompressor returns a new Mock
func NewMockCompressor() compression.Component {
	return compressionnoop.NewNoopStrategy()
}
