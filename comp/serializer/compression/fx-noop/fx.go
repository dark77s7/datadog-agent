// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build !zlib && !zstd

// Package fx provides the fx module for the serializer/compression component
package fx

import (
	compressionnoop "github.com/DataDog/datadog-agent/comp/serializer/compression/impl-noop"
	"go.uber.org/fx"

	"github.com/DataDog/datadog-agent/comp/core/config"
	compression "github.com/DataDog/datadog-agent/comp/serializer/compression/def"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
)

// Module defines the fx options for the component.
func Module() fxutil.Module {
	return fxutil.Component(
		fx.Provide(NewCompressor),
	)
}

// NewCompressor returns a new Compressor based on serializer_compressor_kind
// This function is called only when the zlib build tag is included
func NewCompressor(_ config.Component) compression.Component {
	return compressionnoop.NewNoopStrategy()
}
