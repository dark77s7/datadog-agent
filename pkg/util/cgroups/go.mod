module github.com/DataDog/datadog-agent/pkg/util/cgroups

go 1.22.0

replace (
	github.com/DataDog/datadog-agent/pkg/util/log => ../log
	github.com/DataDog/datadog-agent/pkg/util/pointer => ../pointer
	github.com/DataDog/datadog-agent/pkg/util/scrubber => ../scrubber
)

require (
	github.com/DataDog/datadog-agent/pkg/util/log v0.56.0-rc.3
	github.com/DataDog/datadog-agent/pkg/util/pointer v0.56.0-rc.3
	github.com/containerd/cgroups/v3 v3.0.4
	github.com/google/go-cmp v0.6.0
	github.com/karrick/godirwalk v1.17.0
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/DataDog/datadog-agent/pkg/util/scrubber v0.56.0-rc.3 // indirect
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/moby/sys/userns v0.1.0 // indirect
	github.com/opencontainers/runtime-spec v1.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/sys v0.27.0 // indirect
	google.golang.org/protobuf v1.35.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
