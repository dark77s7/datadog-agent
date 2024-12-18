// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

// Package fx provides the fx module for the authtoken component
package fx

import (
	authtoken "github.com/DataDog/datadog-agent/comp/api/authtoken/def"
	createandfetchimpl "github.com/DataDog/datadog-agent/comp/api/authtoken/impl-createandfetch"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
)

// Module defines the fx options for this component
func Module() fxutil.Module {
	return fxutil.Component(
		fxutil.ProvideComponentConstructor(
			createandfetchimpl.NewComponent,
		),
		fxutil.ProvideOptional[authtoken.Component](),
	)
}
