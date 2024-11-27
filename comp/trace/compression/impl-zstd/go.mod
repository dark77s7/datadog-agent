module github.com/DataDog/datadog-agent/comp/trace/compression/impl-zstd

go 1.22.0

replace github.com/DataDog/datadog-agent/comp/trace/compression/def => ../../../../comp/trace/compression/def/

require (
	github.com/DataDog/datadog-agent/comp/trace/compression/def v0.59.1-rc.1
	github.com/DataDog/zstd v1.5.5
)
