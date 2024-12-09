module github.com/DataDog/datadog-agent/comp/api/authtoken

go 1.22.0

replace (
	github.com/DataDog/datadog-agent/comp/api/api/def => ../../../comp/api/api/def
	github.com/DataDog/datadog-agent/comp/core/config => ../../../comp/core/config
	github.com/DataDog/datadog-agent/comp/core/flare/builder => ../../../comp/core/flare/builder
	github.com/DataDog/datadog-agent/comp/core/flare/types => ../../../comp/core/flare/types
	github.com/DataDog/datadog-agent/comp/core/log/def => ../../../comp/core/log/def
	github.com/DataDog/datadog-agent/comp/core/log/mock => ../../../comp/core/log/mock
	github.com/DataDog/datadog-agent/comp/core/secrets => ../../../comp/core/secrets
	github.com/DataDog/datadog-agent/comp/core/telemetry => ../../../comp/core/telemetry
	github.com/DataDog/datadog-agent/comp/def => ../../../comp/def
	github.com/DataDog/datadog-agent/pkg/api => ../../../pkg/api
	github.com/DataDog/datadog-agent/pkg/collector/check/defaults => ../../../pkg/collector/check/defaults
	github.com/DataDog/datadog-agent/pkg/config/env => ../../../pkg/config/env
	github.com/DataDog/datadog-agent/pkg/config/mock => ../../../pkg/config/mock
	github.com/DataDog/datadog-agent/pkg/config/model => ../../../pkg/config/model
	github.com/DataDog/datadog-agent/pkg/config/nodetreemodel => ../../../pkg/config/nodetreemodel
	github.com/DataDog/datadog-agent/pkg/config/setup => ../../../pkg/config/setup
	github.com/DataDog/datadog-agent/pkg/config/structure => ../../../pkg/config/structure
	github.com/DataDog/datadog-agent/pkg/config/teeconfig => ../../../pkg/config/teeconfig
	github.com/DataDog/datadog-agent/pkg/config/utils => ../../../pkg/config/utils
	github.com/DataDog/datadog-agent/pkg/util/defaultpaths => ../../../pkg/util/defaultpaths
	github.com/DataDog/datadog-agent/pkg/util/executable => ../../../pkg/util/executable
	github.com/DataDog/datadog-agent/pkg/util/filesystem => ../../../pkg/util/filesystem
	github.com/DataDog/datadog-agent/pkg/util/fxutil => ../../../pkg/util/fxutil
	github.com/DataDog/datadog-agent/pkg/util/hostname/validate => ../../../pkg/util/hostname/validate/
	github.com/DataDog/datadog-agent/pkg/util/log => ../../../pkg/util/log
	github.com/DataDog/datadog-agent/pkg/util/log/setup => ../../../pkg/util/log/setup
	github.com/DataDog/datadog-agent/pkg/util/optional => ../../../pkg/util/optional
	github.com/DataDog/datadog-agent/pkg/util/pointer => ../../../pkg/util/pointer
	github.com/DataDog/datadog-agent/pkg/util/scrubber => ../../../pkg/util/scrubber
	github.com/DataDog/datadog-agent/pkg/util/system => ../../../pkg/util/system
	github.com/DataDog/datadog-agent/pkg/util/system/socket => ../../../pkg/util/system/socket/
	github.com/DataDog/datadog-agent/pkg/util/testutil => ../../../pkg/util/testutil
	github.com/DataDog/datadog-agent/pkg/util/winutil => ../../../pkg/util/winutil
	github.com/DataDog/datadog-agent/pkg/version => ../../../pkg/version

)

require (
	github.com/DataDog/datadog-agent/comp/core/config v0.63.0
	github.com/DataDog/datadog-agent/comp/core/log/def v0.63.0-devel
	github.com/DataDog/datadog-agent/comp/core/log/mock v0.63.0-devel
	github.com/DataDog/datadog-agent/pkg/api v0.63.0
	github.com/DataDog/datadog-agent/pkg/util/fxutil v0.63.0
	github.com/DataDog/datadog-agent/pkg/util/optional v0.63.0
	github.com/stretchr/testify v1.10.0
	go.uber.org/fx v1.23.0
)

require (
	github.com/DataDog/datadog-agent/comp/core/flare/builder v0.63.0 // indirect
	github.com/DataDog/datadog-agent/comp/core/flare/types v0.63.0 // indirect
	github.com/DataDog/datadog-agent/comp/core/secrets v0.63.0 // indirect
	github.com/DataDog/datadog-agent/comp/def v0.63.0 // indirect
	github.com/DataDog/datadog-agent/pkg/collector/check/defaults v0.63.0 // indirect
	github.com/DataDog/datadog-agent/pkg/config/env v0.63.0 // indirect
	github.com/DataDog/datadog-agent/pkg/config/mock v0.63.0 // indirect
	github.com/DataDog/datadog-agent/pkg/config/model v0.63.0 // indirect
	github.com/DataDog/datadog-agent/pkg/config/nodetreemodel v0.63.0-devel // indirect
	github.com/DataDog/datadog-agent/pkg/config/setup v0.63.0 // indirect
	github.com/DataDog/datadog-agent/pkg/config/structure v0.63.0 // indirect
	github.com/DataDog/datadog-agent/pkg/config/teeconfig v0.63.0-devel // indirect
	github.com/DataDog/datadog-agent/pkg/config/utils v0.63.0 // indirect
	github.com/DataDog/datadog-agent/pkg/util/executable v0.63.0 // indirect
	github.com/DataDog/datadog-agent/pkg/util/filesystem v0.63.0 // indirect
	github.com/DataDog/datadog-agent/pkg/util/hostname/validate v0.63.0 // indirect
	github.com/DataDog/datadog-agent/pkg/util/log v0.63.0 // indirect
	github.com/DataDog/datadog-agent/pkg/util/log/setup v0.63.0-devel // indirect
	github.com/DataDog/datadog-agent/pkg/util/pointer v0.63.0 // indirect
	github.com/DataDog/datadog-agent/pkg/util/scrubber v0.63.0 // indirect
	github.com/DataDog/datadog-agent/pkg/util/system v0.63.0 // indirect
	github.com/DataDog/datadog-agent/pkg/util/system/socket v0.63.0 // indirect
	github.com/DataDog/datadog-agent/pkg/util/winutil v0.63.0 // indirect
	github.com/DataDog/datadog-agent/pkg/version v0.63.0 // indirect
	github.com/DataDog/viper v1.13.5 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/hashicorp/hcl v1.0.1-vault-5 // indirect
	github.com/hectane/go-acl v0.0.0-20190604041725-da78bae5fc95 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/lufia/plan9stats v0.0.0-20220913051719-115f729f3c8c // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.1-0.20231216201459-8508981c8b6c // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20220216144756-c35f1ee13d7c // indirect
	github.com/shirou/gopsutil/v3 v3.24.5 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.7.0 // indirect
	github.com/spf13/cobra v1.8.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tklauser/go-sysconf v0.3.14 // indirect
	github.com/tklauser/numcpus v0.8.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/dig v1.18.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/exp v0.0.0-20241108190413-2d47ceb2692f // indirect
	golang.org/x/sys v0.27.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
