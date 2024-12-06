// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package newstylefx

import (
	newstyleimpl "github.com/DataDog/datadog-agent/comp/newstyle/impl"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
)

func Module() fxutil.Module {
	return fxutil.Component(
		fxutil.ProvideComponentConstructor(
			newstyleimpl.NewComponent,
		),
	)
}
