// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package profile

import (
	"expvar"
	"fmt"
	"reflect"
	"strings"

	"github.com/DataDog/datadog-agent/pkg/util/log"
	"github.com/mohae/deepcopy"

	"github.com/DataDog/datadog-agent/pkg/networkdevice/profile/profiledefinition"
	"github.com/DataDog/datadog-agent/pkg/networkdevice/utils"

	"github.com/DataDog/datadog-agent/pkg/collector/corechecks/snmp/internal/configvalidation"
)

var (
	profileExpVar = expvar.NewMap("snmpProfileErrors")
)

func resolveProfiles(userProfiles, defaultProfiles ProfileConfigMap) (ProfileConfigMap, error) {
	rawProfiles := mergeProfiles(defaultProfiles, userProfiles)
	userExpandedProfiles, err := loadResolveProfiles(rawProfiles, defaultProfiles)
	if err != nil {
		return nil, fmt.Errorf("failed to load profiles: %w", err)
	}
	return userExpandedProfiles, nil
}

func loadResolveProfiles(pConfig ProfileConfigMap, defaultProfiles ProfileConfigMap) (ProfileConfigMap, error) {
	profiles := make(ProfileConfigMap, len(pConfig))

	for name := range pConfig {
		// No need to resolve abstract profile
		if strings.HasPrefix(name, "_") {
			continue
		}

		newProfileConfig := deepcopy.Copy(pConfig[name]).(ProfileConfig)
		err := recursivelyExpandBaseProfiles(name, &newProfileConfig.Definition, newProfileConfig.Definition.Extends, []string{}, pConfig, defaultProfiles)
		if err != nil {
			log.Warnf("failed to expand profile %q: %v", name, err)
			continue
		}
		profiledefinition.NormalizeMetrics(newProfileConfig.Definition.Metrics)
		errors := configvalidation.ValidateEnrichMetadata(newProfileConfig.Definition.Metadata)
		errors = append(errors, configvalidation.ValidateEnrichMetrics(newProfileConfig.Definition.Metrics)...)
		errors = append(errors, configvalidation.ValidateEnrichMetricTags(newProfileConfig.Definition.MetricTags)...)
		if len(errors) > 0 {
			log.Warnf("validation errors in profile %q: %s", name, strings.Join(errors, "\n"))
			profileExpVar.Set(name, expvar.Func(func() interface{} {
				return strings.Join(errors, "\n")
			}))
			continue
		}
		profiles[name] = newProfileConfig
	}

	return profiles, nil
}

func recursivelyExpandBaseProfiles(parentExtend string, definition *profiledefinition.ProfileDefinition, extends []string, extendsHistory []string, profiles ProfileConfigMap, defaultProfiles ProfileConfigMap) error {
	for _, extendEntry := range extends {
		extendEntry = strings.TrimSuffix(extendEntry, ".yaml")

		var baseDefinition *profiledefinition.ProfileDefinition
		// User profile can extend default profile by extending the default profile.
		// If the extend entry has the same name as the profile name, we assume the extend entry is referring to a default profile.
		if extendEntry == parentExtend {
			profile, ok := defaultProfiles[extendEntry]
			if !ok {
				return fmt.Errorf("extend does not exist: `%s`", extendEntry)
			}
			baseDefinition = &profile.Definition
		} else {
			profile, ok := profiles[extendEntry]
			if !ok {
				profile, ok = defaultProfiles[extendEntry]
				if !ok {
					return fmt.Errorf("extend does not exist: `%s`", extendEntry)
				}
			}
			baseDefinition = &profile.Definition
		}
		for _, extend := range extendsHistory {
			if extend == extendEntry {
				return fmt.Errorf("cyclic profile extend detected, `%s` has already been extended, extendsHistory=`%v`", extendEntry, extendsHistory)
			}
		}

		mergeProfileDefinition(definition, baseDefinition)

		newExtendsHistory := append(utils.CopyStrings(extendsHistory), extendEntry)
		err := recursivelyExpandBaseProfiles(extendEntry, definition, baseDefinition.Extends, newExtendsHistory, profiles, defaultProfiles)
		if err != nil {
			return err
		}
	}
	return nil
}

func mergeProfileDefinition(targetDefinition *profiledefinition.ProfileDefinition, baseDefinition *profiledefinition.ProfileDefinition) {
	targetDefinition.Metrics = append(targetDefinition.Metrics, baseDefinition.Metrics...)
	targetDefinition.MetricTags = append(targetDefinition.MetricTags, baseDefinition.MetricTags...)
	targetDefinition.StaticTags = append(targetDefinition.StaticTags, baseDefinition.StaticTags...)
	// we only have two types of resources: Device and Interface
	// merge Device Metadata Fields
	baseProfileDeviceMetadataFields := reflect.ValueOf(baseDefinition.Metadata.Device.Fields).Elem()
	targetProfileDeviceMetadataFields := reflect.ValueOf(targetDefinition.Metadata.Device.Fields).Elem()
	for i := 0; i < targetProfileDeviceMetadataFields.NumField(); i++ {
		targetField := targetProfileDeviceMetadataFields.Field(i)
		baseField := baseProfileDeviceMetadataFields.Field(i)

		if isEmpty(targetField) {
			targetField.Set(baseField)
		}
	}
	
	// merge Interface Metadata Fields
	baseProfileInterfaceMetadataFields := reflect.ValueOf(baseDefinition.Metadata.Interface.Fields).Elem()
	targetProfileInterfaceMetadataFields := reflect.ValueOf(targetDefinition.Metadata.Interface.Fields).Elem()
	for i := 0; i < targetProfileInterfaceMetadataFields.NumField(); i++ {
		targetField := targetProfileInterfaceMetadataFields.Field(i)
		baseField := baseProfileInterfaceMetadataFields.Field(i)

		if isEmpty(targetField) {
			targetField.Set(baseField)
		}
	}
	// merge Interface Metadata IDTags
	for _, tagConfig := range baseDefinition.Metadata.Interface.IDTags {
		targetDefinition.Metadata.Interface.IDTags = append(targetDefinition.Metadata.Interface.IDTags, tagConfig)
	}
}
