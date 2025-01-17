// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

// Package configsyncimpl implements synchronizing the configuration using the core agent config API
package configsyncimpl

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go.uber.org/fx"

	"github.com/DataDog/datadog-agent/comp/api/authtoken"
	"github.com/DataDog/datadog-agent/comp/core/config"
	"github.com/DataDog/datadog-agent/comp/core/configsync"
	log "github.com/DataDog/datadog-agent/comp/core/log/def"
	apiutil "github.com/DataDog/datadog-agent/pkg/api/util"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
)

type dependencies struct {
	fx.In
	Lc fx.Lifecycle

	Config     config.Component
	Log        log.Component
	Authtoken  authtoken.Component
	SyncParams Params
}

// Module defines the fx options for this component.
func Module() fxutil.Module {
	return fxutil.Component(
		fx.Provide(newComponent),
		fx.Supply(Params{}),
	)
}

// ModuleWithParams defines the fx options for this component, but
// requires additionally specifying custom Params from the fx App, to be
// passed to the constructor.
func ModuleWithParams() fxutil.Module {
	return fxutil.Component(
		fx.Provide(newComponent),
	)
}

type configSync struct {
	Config    config.Component
	Log       log.Component
	Authtoken authtoken.Component

	url       *url.URL
	client    *http.Client
	connected bool
	ctx       context.Context
	enabled   bool
}

// newComponent checks if the component was enabled as per the config and return a enable/disabled configsync
func newComponent(deps dependencies) configsync.Component {
	agentIPCPort := deps.Config.GetInt("agent_ipc.port")
	configRefreshIntervalSec := deps.Config.GetInt("agent_ipc.config_refresh_interval")

	if agentIPCPort <= 0 || configRefreshIntervalSec <= 0 {
		deps.Log.Infof("configsync disabled (agent_ipc.port: %d | agent_ipc.config_refresh_interval: %d)", agentIPCPort, configRefreshIntervalSec)
		return configSync{}
	}

	deps.Log.Infof("configsync enabled (agent_ipc '%s:%d' | agent_ipc.config_refresh_interval: %d)", deps.Config.GetString("agent_ipc.host"), agentIPCPort, configRefreshIntervalSec)
	return newConfigSync(deps, agentIPCPort, configRefreshIntervalSec)
}

// newConfigSync creates a new configSync component.
// agentIPCPort and configRefreshIntervalSec must be strictly positive.
func newConfigSync(deps dependencies, agentIPCPort int, configRefreshIntervalSec int) configsync.Component {
	agentIPCHost := deps.Config.GetString("agent_ipc.host")

	url := &url.URL{
		Scheme: "https",
		Host:   net.JoinHostPort(agentIPCHost, strconv.Itoa(agentIPCPort)),
		Path:   "/config/v1",
	}

	ctx, cancel := context.WithCancel(context.Background())
	client := apiutil.GetClientWithTimeout(deps.SyncParams.Timeout, false)
	configRefreshInterval := time.Duration(configRefreshIntervalSec) * time.Second

	configSync := configSync{
		Config:    deps.Config,
		Log:       deps.Log,
		Authtoken: deps.Authtoken,
		url:       url,
		client:    client,
		ctx:       ctx,
		enabled:   true,
	}

	if deps.SyncParams.OnInit {
		if deps.SyncParams.Delay != 0 {
			select {
			case <-ctx.Done(): //context cancelled
				// TODO: this component should return an error
				cancel()
				return nil
			case <-time.After(deps.SyncParams.Delay):
			}
		}
		configSync.updater()
	}

	// start and stop the routine in fx hooks
	deps.Lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go configSync.runWithInterval(configRefreshInterval)
			return nil
		},
		OnStop: func(_ context.Context) error {
			cancel()
			return nil
		},
	})

	return configSync
}
