// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package profile

import (
	"github.com/DataDog/datadog-agent/pkg/networkdevice/profile/profiledefinition"
)

// ProfileConfigMap represent a map of ProfileConfig
type ProfileConfigMap map[string]ProfileConfig

func (m ProfileConfigMap) Clone() ProfileConfigMap {
	return profiledefinition.CloneMap(m)
}

// ProfileConfig represent a profile configuration
type ProfileConfig struct {
	DefinitionFile string                              `yaml:"definition_file"`
	Definition     profiledefinition.ProfileDefinition `yaml:"definition"`

	IsUserProfile bool `yaml:"-"`
}

func (p ProfileConfig) Clone() ProfileConfig {
	return ProfileConfig{
		DefinitionFile: p.DefinitionFile,
		Definition:     *p.Definition.Clone(),
		IsUserProfile:  p.IsUserProfile,
	}
}
