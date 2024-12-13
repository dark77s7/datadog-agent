module github.com/DataDog/datadog-agent/pkg/orchestrator/model

go 1.22.0

replace (
	github.com/DataDog/datadog-agent/pkg/util/log => ../../util/log/
	github.com/DataDog/datadog-agent/pkg/util/scrubber => ../../util/scrubber/
)

require (
	github.com/DataDog/datadog-agent/pkg/util/log v0.60.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
)

require (
	github.com/DataDog/datadog-agent/pkg/util/scrubber v0.60.0 // indirect
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
