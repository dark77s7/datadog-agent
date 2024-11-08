// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Code generated by "go.opentelemetry.io/collector/cmd/builder". DO NOT EDIT.

package collectorcontribimpl

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/receiver"
	spanmetricsconnector "github.com/open-telemetry/opentelemetry-collector-contrib/connector/spanmetricsconnector"
	debugexporter "go.opentelemetry.io/collector/exporter/debugexporter"
	nopexporter "go.opentelemetry.io/collector/exporter/nopexporter"
	otlpexporter "go.opentelemetry.io/collector/exporter/otlpexporter"
	otlphttpexporter "go.opentelemetry.io/collector/exporter/otlphttpexporter"
	sapmexporter "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/sapmexporter"
	zpagesextension "go.opentelemetry.io/collector/extension/zpagesextension"
	healthcheckextension "github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension"
	pprofextension "github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension"
	dockerobserver "github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/dockerobserver"
	ecsobserver "github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/ecsobserver"
	ecstaskobserver "github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/ecstaskobserver"
	hostobserver "github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/hostobserver"
	k8sobserver "github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/k8sobserver"
	batchprocessor "go.opentelemetry.io/collector/processor/batchprocessor"
	memorylimiterprocessor "go.opentelemetry.io/collector/processor/memorylimiterprocessor"
	attributesprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor"
	cumulativetodeltaprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/cumulativetodeltaprocessor"
	filterprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor"
	groupbyattrsprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor"
	k8sattributesprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor"
	probabilisticsamplerprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor"
	resourcedetectionprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor"
	resourceprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor"
	routingprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/routingprocessor"
	tailsamplingprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/tailsamplingprocessor"
	transformprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
	nopreceiver "go.opentelemetry.io/collector/receiver/nopreceiver"
	otlpreceiver "go.opentelemetry.io/collector/receiver/otlpreceiver"
	filelogreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver"
	fluentforwardreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/fluentforwardreceiver"
	hostmetricsreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver"
	jaegerreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver"
	prometheusreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver"
	receivercreator "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/receivercreator"
	zipkinreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinreceiver"
)

func components() (otelcol.Factories, error) {
	var err error
	factories := otelcol.Factories{}

	factories.Extensions, err = extension.MakeFactoryMap(
		zpagesextension.NewFactory(),
		healthcheckextension.NewFactory(),
		pprofextension.NewFactory(),
		dockerobserver.NewFactory(),
		ecsobserver.NewFactory(),
		ecstaskobserver.NewFactory(),
		hostobserver.NewFactory(),
		k8sobserver.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}
	factories.ExtensionModules = make(map[component.Type]string, len(factories.Extensions))
	factories.ExtensionModules[zpagesextension.NewFactory().Type()] = "go.opentelemetry.io/collector/extension/zpagesextension v0.111.0"
	factories.ExtensionModules[healthcheckextension.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension v0.111.0"
	factories.ExtensionModules[pprofextension.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension v0.111.0"
	factories.ExtensionModules[dockerobserver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/dockerobserver v0.111.0"
	factories.ExtensionModules[ecsobserver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/ecsobserver v0.111.0"
	factories.ExtensionModules[ecstaskobserver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/ecstaskobserver v0.111.0"
	factories.ExtensionModules[hostobserver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/hostobserver v0.111.0"
	factories.ExtensionModules[k8sobserver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/k8sobserver v0.111.0"

	factories.Receivers, err = receiver.MakeFactoryMap(
		nopreceiver.NewFactory(),
		otlpreceiver.NewFactory(),
		filelogreceiver.NewFactory(),
		fluentforwardreceiver.NewFactory(),
		hostmetricsreceiver.NewFactory(),
		jaegerreceiver.NewFactory(),
		prometheusreceiver.NewFactory(),
		receivercreator.NewFactory(),
		zipkinreceiver.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}
	factories.ReceiverModules = make(map[component.Type]string, len(factories.Receivers))
	factories.ReceiverModules[nopreceiver.NewFactory().Type()] = "go.opentelemetry.io/collector/receiver/nopreceiver v0.111.0"
	factories.ReceiverModules[otlpreceiver.NewFactory().Type()] = "go.opentelemetry.io/collector/receiver/otlpreceiver v0.111.0"
	factories.ReceiverModules[filelogreceiver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver v0.111.0"
	factories.ReceiverModules[fluentforwardreceiver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/fluentforwardreceiver v0.111.0"
	factories.ReceiverModules[hostmetricsreceiver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver v0.111.0"
	factories.ReceiverModules[jaegerreceiver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver v0.111.0"
	factories.ReceiverModules[prometheusreceiver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver v0.111.0"
	factories.ReceiverModules[receivercreator.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/receivercreator v0.111.0"
	factories.ReceiverModules[zipkinreceiver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinreceiver v0.111.0"

	factories.Exporters, err = exporter.MakeFactoryMap(
		debugexporter.NewFactory(),
		nopexporter.NewFactory(),
		otlpexporter.NewFactory(),
		otlphttpexporter.NewFactory(),
		sapmexporter.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}
	factories.ExporterModules = make(map[component.Type]string, len(factories.Exporters))
	factories.ExporterModules[debugexporter.NewFactory().Type()] = "go.opentelemetry.io/collector/exporter/debugexporter v0.111.0"
	factories.ExporterModules[nopexporter.NewFactory().Type()] = "go.opentelemetry.io/collector/exporter/nopexporter v0.111.0"
	factories.ExporterModules[otlpexporter.NewFactory().Type()] = "go.opentelemetry.io/collector/exporter/otlpexporter v0.111.0"
	factories.ExporterModules[otlphttpexporter.NewFactory().Type()] = "go.opentelemetry.io/collector/exporter/otlphttpexporter v0.111.0"
	factories.ExporterModules[sapmexporter.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/sapmexporter v0.111.0"

	factories.Processors, err = processor.MakeFactoryMap(
		batchprocessor.NewFactory(),
		memorylimiterprocessor.NewFactory(),
		attributesprocessor.NewFactory(),
		cumulativetodeltaprocessor.NewFactory(),
		filterprocessor.NewFactory(),
		groupbyattrsprocessor.NewFactory(),
		k8sattributesprocessor.NewFactory(),
		probabilisticsamplerprocessor.NewFactory(),
		resourcedetectionprocessor.NewFactory(),
		resourceprocessor.NewFactory(),
		routingprocessor.NewFactory(),
		tailsamplingprocessor.NewFactory(),
		transformprocessor.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}
	factories.ProcessorModules = make(map[component.Type]string, len(factories.Processors))
	factories.ProcessorModules[batchprocessor.NewFactory().Type()] = "go.opentelemetry.io/collector/processor/batchprocessor v0.111.0"
	factories.ProcessorModules[memorylimiterprocessor.NewFactory().Type()] = "go.opentelemetry.io/collector/processor/memorylimiterprocessor v0.111.0"
	factories.ProcessorModules[attributesprocessor.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor v0.111.0"
	factories.ProcessorModules[cumulativetodeltaprocessor.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/cumulativetodeltaprocessor v0.111.0"
	factories.ProcessorModules[filterprocessor.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor v0.111.0"
	factories.ProcessorModules[groupbyattrsprocessor.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor v0.111.0"
	factories.ProcessorModules[k8sattributesprocessor.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor v0.111.0"
	factories.ProcessorModules[probabilisticsamplerprocessor.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor v0.111.0"
	factories.ProcessorModules[resourcedetectionprocessor.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor v0.111.0"
	factories.ProcessorModules[resourceprocessor.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor v0.111.0"
	factories.ProcessorModules[routingprocessor.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/routingprocessor v0.111.0"
	factories.ProcessorModules[tailsamplingprocessor.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/tailsamplingprocessor v0.111.0"
	factories.ProcessorModules[transformprocessor.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor v0.111.0"

	factories.Connectors, err = connector.MakeFactoryMap(
		spanmetricsconnector.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}
	factories.ConnectorModules = make(map[component.Type]string, len(factories.Connectors))
	factories.ConnectorModules[spanmetricsconnector.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/connector/spanmetricsconnector v0.111.0"

	return factories, nil
}
