// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

// Package fx provides the fx module for the haagent component
package fx

import (
	haagent "github.com/DataDog/datadog-agent/comp/haagent/def"
	haagentimpl "github.com/DataDog/datadog-agent/comp/haagent/impl"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
)

// Module defines the fx options for this component
func Module() fxutil.Module {
	return fxutil.Component(
		fxutil.ProvideComponentConstructor(
			haagentimpl.NewComponent,
		),
		fxutil.ProvideOptional[haagent.Component](),
	)
}
