module github.com/DataDog/datadog-agent/comp/core/config

go 1.22.0

replace (
	github.com/DataDog/datadog-agent/comp/api/api/def => ../../api/api/def
	github.com/DataDog/datadog-agent/comp/core/flare/builder => ../flare/builder
	github.com/DataDog/datadog-agent/comp/core/flare/types => ../flare/types
	github.com/DataDog/datadog-agent/comp/core/secrets => ../secrets/
	github.com/DataDog/datadog-agent/comp/core/telemetry => ../telemetry/
	github.com/DataDog/datadog-agent/comp/def => ../../def
	github.com/DataDog/datadog-agent/pkg/collector/check/defaults => ../../../pkg/collector/check/defaults
	github.com/DataDog/datadog-agent/pkg/config/env => ../../../pkg/config/env
	github.com/DataDog/datadog-agent/pkg/config/mock => ../../../pkg/config/mock
	github.com/DataDog/datadog-agent/pkg/config/model => ../../../pkg/config/model/
	github.com/DataDog/datadog-agent/pkg/config/nodetreemodel => ../../../pkg/config/nodetreemodel
	github.com/DataDog/datadog-agent/pkg/config/setup => ../../../pkg/config/setup
	github.com/DataDog/datadog-agent/pkg/config/teeconfig => ../../../pkg/config/teeconfig
	github.com/DataDog/datadog-agent/pkg/telemetry => ../../../pkg/telemetry
	github.com/DataDog/datadog-agent/pkg/util/defaultpaths => ../../../pkg/util/defaultpaths
	github.com/DataDog/datadog-agent/pkg/util/executable => ../../../pkg/util/executable
	github.com/DataDog/datadog-agent/pkg/util/filesystem => ../../../pkg/util/filesystem
	github.com/DataDog/datadog-agent/pkg/util/fxutil => ../../../pkg/util/fxutil
	github.com/DataDog/datadog-agent/pkg/util/hostname/validate => ../../../pkg/util/hostname/validate
	github.com/DataDog/datadog-agent/pkg/util/log => ../../../pkg/util/log
	github.com/DataDog/datadog-agent/pkg/util/optional => ../../../pkg/util/optional/
	github.com/DataDog/datadog-agent/pkg/util/pointer => ../../../pkg/util/pointer
	github.com/DataDog/datadog-agent/pkg/util/scrubber => ../../../pkg/util/scrubber/
	github.com/DataDog/datadog-agent/pkg/util/system => ../../../pkg/util/system
	github.com/DataDog/datadog-agent/pkg/util/system/socket => ../../../pkg/util/system/socket
	github.com/DataDog/datadog-agent/pkg/util/testutil => ../../../pkg/util/testutil
	github.com/DataDog/datadog-agent/pkg/util/winutil => ../../../pkg/util/winutil
)

require (
	github.com/DataDog/datadog-agent/comp/core/flare/types v0.61.0-rc.1
	github.com/DataDog/datadog-agent/comp/core/secrets v0.61.0-rc.1
	github.com/DataDog/datadog-agent/comp/core/telemetry v0.61.0-rc.1
	github.com/DataDog/datadog-agent/pkg/config/mock v0.61.0-rc.1
	github.com/DataDog/datadog-agent/pkg/config/model v0.61.0-rc.1
	github.com/DataDog/datadog-agent/pkg/config/setup v0.61.0-rc.1
	github.com/DataDog/datadog-agent/pkg/util/defaultpaths v0.61.0-rc.1
	github.com/DataDog/datadog-agent/pkg/util/fxutil v0.61.0-rc.1
	github.com/DataDog/datadog-agent/pkg/util/optional v0.61.0-rc.1
	github.com/DataDog/datadog-agent/pkg/util/winutil v0.61.0-rc.1
	github.com/DataDog/viper v1.13.5
	github.com/stretchr/testify v1.9.0
	go.uber.org/fx v1.22.2
)

require (
	github.com/DataDog/datadog-agent/comp/api/api/def v0.61.0-rc.1 // indirect
	github.com/DataDog/datadog-agent/comp/core/flare/builder v0.61.0-rc.1 // indirect
	github.com/DataDog/datadog-agent/comp/def v0.61.0-rc.1 // indirect
	github.com/DataDog/datadog-agent/pkg/collector/check/defaults v0.61.0-rc.1 // indirect
	github.com/DataDog/datadog-agent/pkg/config/env v0.61.0-rc.1 // indirect
	github.com/DataDog/datadog-agent/pkg/config/nodetreemodel v0.61.0-rc.1 // indirect
	github.com/DataDog/datadog-agent/pkg/config/structure v0.61.0-rc.1 // indirect
	github.com/DataDog/datadog-agent/pkg/config/teeconfig v0.61.0-rc.1 // indirect
	github.com/DataDog/datadog-agent/pkg/util/executable v0.61.0-rc.1 // indirect
	github.com/DataDog/datadog-agent/pkg/util/filesystem v0.61.0-rc.1 // indirect
	github.com/DataDog/datadog-agent/pkg/util/hostname/validate v0.61.0-rc.1 // indirect
	github.com/DataDog/datadog-agent/pkg/util/log v0.61.0-rc.1 // indirect
	github.com/DataDog/datadog-agent/pkg/util/pointer v0.61.0-rc.1 // indirect
	github.com/DataDog/datadog-agent/pkg/util/scrubber v0.61.0-rc.1 // indirect
	github.com/DataDog/datadog-agent/pkg/util/system v0.61.0-rc.1 // indirect
	github.com/DataDog/datadog-agent/pkg/util/system/socket v0.61.0-rc.1 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hectane/go-acl v0.0.0-20190604041725-da78bae5fc95 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/pelletier/go-toml v1.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/prometheus/client_golang v1.17.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.1 // indirect
	github.com/shirou/gopsutil/v3 v3.23.12 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.7.0 // indirect
	github.com/spf13/cobra v1.7.0 // indirect
	github.com/spf13/jwalterweatherman v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	go.opentelemetry.io/otel v1.27.0 // indirect
	go.opentelemetry.io/otel/metric v1.27.0 // indirect
	go.opentelemetry.io/otel/sdk v1.27.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.27.0 // indirect
	go.opentelemetry.io/otel/trace v1.27.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/dig v1.18.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/exp v0.0.0-20241108190413-2d47ceb2692f // indirect
	golang.org/x/mod v0.22.0 // indirect
	golang.org/x/sync v0.9.0 // indirect
	golang.org/x/sys v0.27.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	golang.org/x/tools v0.27.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/DataDog/datadog-agent/pkg/config/structure => ../../../pkg/config/structure
