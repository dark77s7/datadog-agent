// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build !windows && kubeapiserver

// Package diagnose implements 'cluster-agent diagnose'.
package diagnose

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/DataDog/datadog-agent/cmd/cluster-agent/command"
	"github.com/DataDog/datadog-agent/comp/core"
	"github.com/DataDog/datadog-agent/comp/core/config"
	log "github.com/DataDog/datadog-agent/comp/core/log/def"
	"github.com/DataDog/datadog-agent/comp/core/secrets"
	compressionfx "github.com/DataDog/datadog-agent/comp/serializer/compression/fx"
	"github.com/DataDog/datadog-agent/pkg/diagnose"
	"github.com/DataDog/datadog-agent/pkg/diagnose/diagnosis"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
)

// Commands returns a slice of subcommands for the 'cluster-agent' command.
func Commands(globalParams *command.GlobalParams) []*cobra.Command {
	cmd := &cobra.Command{
		Use:   "diagnose",
		Short: "Execute some connectivity diagnosis on your system",
		Long:  ``,
		RunE: func(_ *cobra.Command, _ []string) error {
			return fxutil.OneShot(run,
				fx.Supply(core.BundleParams{
					ConfigParams: config.NewClusterAgentParams(globalParams.ConfFilePath),
					SecretParams: secrets.NewEnabledParams(),
					LogParams:    log.ForOneShot(command.LoggerName, "off", true), // no need to show regular logs
				}),
				core.Bundle(),
				compressionfx.Module(),
			)
		},
	}

	return []*cobra.Command{cmd}
}

func run(_ config.Component) error {
	// Verbose:  true - to show details like if was done a while ago
	// RunLocal: true - do not attept to run in actual running agent but
	//                  may need to implement it in future
	// Include: connectivity-datadog-autodiscovery - limit to a single
	//                  diagnose suite as it was done in this agent for
	//                  a while. Most likely need to relax or add more
	//                  diagnose suites in the future
	diagCfg := diagnosis.Config{Verbose: true, RunLocal: true}
	diagnoses, err := diagnose.RunLocalCheck(diagCfg, diagnose.RegisterConnectivityAutodiscovery)
	if err != nil {
		return err
	}
	return diagnose.RunDiagnoseStdOut(color.Output, diagCfg, diagnoses)
}
