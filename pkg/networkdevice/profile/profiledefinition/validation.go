// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package profiledefinition

import (
	"fmt"
	"regexp"
)

var validMetadataResources = map[string]map[string]bool{
	"device": {
		"name":          true,
		"description":   true,
		"sys_object_id": true,
		"location":      true,
		"serial_number": true,
		"vendor":        true,
		"version":       true,
		"product_name":  true,
		"model":         true,
		"os_name":       true,
		"os_version":    true,
		"os_hostname":   true,
		"type":          true,
	},
	"interface": {
		"name":         true,
		"alias":        true,
		"description":  true,
		"mac_address":  true,
		"admin_status": true,
		"oper_status":  true,
	},
}

// SymbolContext represent the context in which the symbol is used
type SymbolContext int64

// ScalarSymbol enums
const (
	ScalarSymbol SymbolContext = iota
	ColumnSymbol
	MetricTagSymbol
	MetadataSymbol
)

// ValidateEnrichProfile validates a profile and normalizes it.
func ValidateEnrichProfile(profile *ProfileDefinition) []string {
	NormalizeMetrics(profile.Metrics)
	if profile.Device.Vendor != "" {
		dev, ok := profile.Metadata["device"]
		if !ok {
			profile.Metadata["device"] = MetadataResourceConfig{
				Fields: make(map[string]MetadataField),
			}
			dev = profile.Metadata["device"]
		}
		_, ok = dev.Fields["vendor"]
		if !ok {
			dev.Fields["vendor"] = MetadataField{
				Value: profile.Device.Vendor,
			}
		}
	}
	profile.Device.Vendor = ""
	errors := ValidateEnrichMetadata(profile.Metadata)
	errors = append(errors, ValidateEnrichMetrics(profile.Metrics)...)
	errors = append(errors, ValidateEnrichMetricTags(profile.MetricTags)...)
	return errors
}

// ValidateEnrichMetricTags validates and normalizes metric tags
func ValidateEnrichMetricTags(metricTags []MetricTagConfig) []string {
	var errors []string
	for i := range metricTags {
		errors = append(errors, validateEnrichMetricTag(&metricTags[i])...)
	}
	return errors
}

// ValidateEnrichMetrics will validate MetricsConfig and enrich it.
// Example of enrichment:
// - storage of compiled regex pattern
func ValidateEnrichMetrics(metrics []MetricsConfig) []string {
	var errors []string
	for i := range metrics {
		metricConfig := &metrics[i]
		if !metricConfig.IsScalar() && !metricConfig.IsColumn() {
			errors = append(errors, fmt.Sprintf("either a table of symbols or a scalar symbol must be provided: %#v", metricConfig))
		}
		if metricConfig.IsScalar() && metricConfig.IsColumn() {
			errors = append(errors, fmt.Sprintf("table symbols and scalar symbol cannot be both provided: %#v", metricConfig))
		}
		// If the entry has a metric_type or the obsolete forced_type, migrate it into the symbols
		metricType := metricConfig.MetricType
		if metricType == "" && metricConfig.ForcedType != "" {
			metricType = metricConfig.ForcedType
		}
		metricConfig.MetricType = ""
		metricConfig.ForcedType = ""
		if metricType != "" {
			if metricConfig.IsScalar() {
				symbol := &metricConfig.Symbol
				if symbol.MetricType == ProfileMetricTypeUnset {
					symbol.MetricType = metricType
				} else if symbol.MetricType != metricType {
					errors = append(errors, fmt.Sprintf("deprecated metric_config.metric_type %s conflicts with symbol type %s at symbol %s", metricType, symbol.MetricType, symbol.Name))
				}
			} else {
				for i, symbol := range metricConfig.Symbols {
					if symbol.MetricType == ProfileMetricTypeUnset {
						metricConfig.Symbols[i].MetricType = metricType
					} else if symbol.MetricType != metricType {
						errors = append(errors, fmt.Sprintf("deprecated metric_config.metric_type %s conflicts with symbol type %s at symbol %s", metricType, symbol.MetricType, symbol.Name))
					}
				}
			}
		}
		if metricConfig.IsScalar() {
			errors = append(errors, validateEnrichSymbol(&metricConfig.Symbol, ScalarSymbol)...)
		}
		if metricConfig.IsColumn() {
			for j := range metricConfig.Symbols {
				errors = append(errors, validateEnrichSymbol(&metricConfig.Symbols[j], ColumnSymbol)...)
			}
			if len(metricConfig.MetricTags) == 0 {
				errors = append(errors, fmt.Sprintf("column symbols doesn't have a 'metric_tags' section (%+v), all its metrics will use the same tags; "+
					"if the table has multiple rows, only one row will be submitted; "+
					"please add at least one discriminating metric tag (such as a row index) "+
					"to ensure metrics of all rows are submitted", metricConfig.Symbols))
			}
			for i := range metricConfig.MetricTags {
				metricTag := &metricConfig.MetricTags[i]
				errors = append(errors, validateEnrichMetricTag(metricTag)...)
			}
		}
		// These are not exposed to JSON or the UI right now
		if len(metricConfig.StaticTags) != 0 {
			errors = append(errors, fmt.Sprintf("static tags are not supported: %#v", metricConfig))
		}
		if metricConfig.Options.Placement != 0 || metricConfig.Options.MetricSuffix != "" {
			errors = append(errors, fmt.Sprintf("metricConfig.Options is not supported: %#v", metricConfig))
		}
	}
	return errors
}

// ValidateEnrichMetadata will validate MetadataConfig and enrich it.
func ValidateEnrichMetadata(metadata MetadataConfig) []string {
	var errors []string
	for resName := range metadata {
		_, isValidRes := validMetadataResources[resName]
		if !isValidRes {
			errors = append(errors, fmt.Sprintf("invalid resource: %s", resName))
		} else {
			res := metadata[resName]
			for fieldName := range res.Fields {
				_, isValidField := validMetadataResources[resName][fieldName]
				if !isValidField {
					errors = append(errors, fmt.Sprintf("invalid resource (%s) field: %s", resName, fieldName))
					continue
				}
				field := res.Fields[fieldName]
				for i := range field.Symbols {
					errors = append(errors, validateEnrichSymbol(&field.Symbols[i], MetadataSymbol)...)
				}
				if field.Symbol.OID != "" {
					errors = append(errors, validateEnrichSymbol(&field.Symbol, MetadataSymbol)...)
				}
				res.Fields[fieldName] = field
			}
			metadata[resName] = res
		}
		if resName == "device" && len(metadata[resName].IDTags) > 0 {
			errors = append(errors, "device resource does not support custom id_tags")
		}
		for i := range metadata[resName].IDTags {
			metricTag := &metadata[resName].IDTags[i]
			errors = append(errors, validateEnrichMetricTag(metricTag)...)
		}
	}
	return errors
}

func validateEnrichSymbol(symbol *SymbolConfig, symbolContext SymbolContext) []string {
	var errors []string
	if symbol.Name == "" {
		errors = append(errors, fmt.Sprintf("symbol name missing: name=`%s` oid=`%s`", symbol.Name, symbol.OID))
	}
	// Percent is deprecated in favor of rate and scale factor
	if symbol.MetricType == ProfileMetricTypePercent {
		symbol.MetricType = ProfileMetricTypeRate
		if symbol.ScaleFactor == 0 {
			symbol.ScaleFactor = 1
		}
		symbol.ScaleFactor *= 100
	}
	// Counter is deprecated in favor of rate
	if symbol.MetricType == ProfileMetricTypeCounter {
		symbol.MetricType = ProfileMetricTypeRate
	}
	// Flag stream isn't supported in the frontend (yet)
	if symbol.MetricType == ProfileMetricTypeFlagStream {
		errors = append(errors, fmt.Sprintf("metric type %s is not supported (name=%q, oid=%q)", symbol.MetricType, symbol.Name, symbol.OID))
	}
	if symbol.OID == "" {
		if symbolContext == ColumnSymbol && !symbol.ConstantValueOne {
			errors = append(errors, fmt.Sprintf("symbol oid or constant_value_one missing: name=`%s` oid=`%s`", symbol.Name, symbol.OID))
		} else if symbolContext != ColumnSymbol {
			errors = append(errors, fmt.Sprintf("symbol oid missing: name=`%s` oid=`%s`", symbol.Name, symbol.OID))
		}
	}
	if symbol.ExtractValue != "" {
		pattern, err := regexp.Compile(symbol.ExtractValue)
		if err != nil {
			errors = append(errors, fmt.Sprintf("cannot compile `extract_value` (%s): %s", symbol.ExtractValue, err.Error()))
		} else {
			symbol.ExtractValueCompiled = pattern
		}
	}
	if symbol.MatchPattern != "" {
		pattern, err := regexp.Compile(symbol.MatchPattern)
		if err != nil {
			errors = append(errors, fmt.Sprintf("cannot compile `extract_value` (%s): %s", symbol.ExtractValue, err.Error()))
		} else {
			symbol.MatchPatternCompiled = pattern
		}
	}
	if symbolContext != ColumnSymbol && symbol.ConstantValueOne {
		errors = append(errors, "`constant_value_one` cannot be used outside of tables")
	}
	if (symbolContext != ColumnSymbol && symbolContext != ScalarSymbol) && symbol.MetricType != "" {
		errors = append(errors, "`metric_type` cannot be used outside scalar/table metric symbols and metrics root")
	}
	return errors
}
func validateEnrichMetricTag(metricTag *MetricTagConfig) []string {
	var errors []string
	if (metricTag.Column.OID != "" || metricTag.Column.Name != "") && (metricTag.Symbol.OID != "" || metricTag.Symbol.Name != "") {
		errors = append(errors, fmt.Sprintf("metric tag symbol and column cannot be both declared: symbol=%v, column=%v", metricTag.Symbol, metricTag.Column))
	}

	// Move deprecated metricTag.Column to metricTag.Symbol
	if metricTag.Column.OID != "" || metricTag.Column.Name != "" {
		metricTag.Symbol = SymbolConfigCompat(metricTag.Column)
		metricTag.Column = SymbolConfig{}
	}

	// OID/Name to Symbol harmonization:
	// When users declare metric tag like:
	//   metric_tags:
	//     - OID: 1.2.3
	//       symbol: aSymbol
	// this will lead to OID stored as MetricTagConfig.OID  and name stored as MetricTagConfig.Symbol.Name
	// When this happens, we harmonize by moving MetricTagConfig.OID to MetricTagConfig.Symbol.OID.
	if metricTag.OID != "" && metricTag.Symbol.OID != "" {
		errors = append(errors, fmt.Sprintf("metric tag OID and symbol.OID cannot be both declared: OID=%s, symbol.OID=%s", metricTag.OID, metricTag.Symbol.OID))
	}
	if metricTag.OID != "" && metricTag.Symbol.OID == "" {
		metricTag.Symbol.OID = metricTag.OID
		metricTag.OID = ""
	}
	if metricTag.Symbol.OID != "" || metricTag.Symbol.Name != "" {
		symbol := SymbolConfig(metricTag.Symbol)
		errors = append(errors, validateEnrichSymbol(&symbol, MetricTagSymbol)...)
		metricTag.Symbol = SymbolConfigCompat(symbol)
	}
	if metricTag.Match != "" {
		errors = append(errors, "MetricTag.Match not supported.")
		pattern, err := regexp.Compile(metricTag.Match)
		if err != nil {
			errors = append(errors, fmt.Sprintf("cannot compile `match` (`%s`): %s", metricTag.Match, err.Error()))
		} else {
			metricTag.Pattern = pattern
		}
		if len(metricTag.Tags) == 0 {
			errors = append(errors, fmt.Sprintf("`tags` mapping must be provided if `match` (`%s`) is defined", metricTag.Match))
		}
	}
	if len(metricTag.Tags) > 0 {
		errors = append(errors, "MetricTag.Tags not supported.")
	}
	if len(metricTag.Mapping) > 0 && metricTag.Tag == "" {
		errors = append(errors, fmt.Sprintf("``tag` must be provided if `mapping` (`%s`) is defined", metricTag.Mapping))
	}
	for _, transform := range metricTag.IndexTransform {
		if transform.Start > transform.End {
			errors = append(errors, fmt.Sprintf("transform rule end should be greater than start. Invalid rule: %#v", transform))
		}
	}
	return errors
}
