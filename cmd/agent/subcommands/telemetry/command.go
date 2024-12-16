// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package telemetry implements 'agent telemetry'.
package telemetry

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/DataDog/datadog-agent/cmd/agent/command"
	"github.com/DataDog/datadog-agent/pkg/flare"
)

func Commands(globalParams *command.GlobalParams) []*cobra.Command {
	cmd := &cobra.Command{
		Use:   "telemetry",
		Short: "Print the telemetry metrics exposed by the agent",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			payload, err := flare.QueryAgentTelemetry()
			if err != nil {
				return err
			}
			fmt.Print(string(payload))
			return nil
		},
	}

	return []*cobra.Command{cmd}
}
