// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build !linux

// Package process implements the local process collector for Workloadmeta.
package process

import (
	"github.com/DataDog/datadog-agent/comp/core/config"
	wmcatalog "github.com/DataDog/datadog-agent/comp/core/wmcatalog/def"
)

// NewCollector is a no-op constructor
func NewCollector(_ config.Component) (wmcatalog.Collector, error) {
	return nil, nil
}
