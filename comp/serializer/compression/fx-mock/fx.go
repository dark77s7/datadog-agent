// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

//go:build test

// Package fx provides the fx module for the serializer/compression component
package fx

import (
	compression "github.com/DataDog/datadog-agent/comp/serializer/compression/def"
	compressionnoop "github.com/DataDog/datadog-agent/comp/serializer/compression/impl-noop"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
)

// MockModule defines the fx options for the mock component.
func MockModule() fxutil.Module {
	return fxutil.Component(
		fxutil.ProvideComponentConstructor(
			NewMockCompressor,
		),
	)
}

// MockModuleFactory defines the fx options for the mock component factory.
func MockModuleFactory() fxutil.Module {
	return fxutil.Component(
		fxutil.ProvideComponentConstructor(
			NewMockCompressorFactory,
		),
	)
}

type mockFactory struct{}

func (*mockFactory) NewNoopCompressor() compression.Component {
	return compressionnoop.NewComponent().Comp
}

func (*mockFactory) NewCompressor(_ string, _ int, _ string, _ []string) compression.Component {
	return compressionnoop.NewComponent().Comp
}

// NewMockCompressorFactory returns a factory that always return a Noop Compressor
func NewMockCompressorFactory() compression.Factory {
	return &mockFactory{}
}

// NewMockCompressor returns a new Mock
func NewMockCompressor() compression.Component {
	return compressionnoop.NewComponent().Comp
}
