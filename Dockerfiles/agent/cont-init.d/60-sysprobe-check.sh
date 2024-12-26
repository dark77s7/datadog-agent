#!/bin/bash

if grep -Eq '^ *enable_tcp_queue_length *: *true' /etc/datadog-agent/system-probe.yaml || [[ "$DD_SYSTEM_PROBE_CONFIG_ENABLE_TCP_QUEUE_LENGTH" == "true" ]]; then
  if [ -f /etc/datadog-agent/conf.d/tcp_queue_length.d/conf.yaml.example ]; then
    mv /etc/datadog-agent/conf.d/tcp_queue_length.d/conf.yaml.example \
       /etc/datadog-agent/conf.d/tcp_queue_length.d/conf.yaml.default
  fi
fi

if grep -Eq '^ *enable_oom_kill *: *true' /etc/datadog-agent/system-probe.yaml || [[ "$DD_SYSTEM_PROBE_CONFIG_ENABLE_OOM_KILL" == "true" ]]; then
  if [ -f /etc/datadog-agent/conf.d/oom_kill.d/conf.yaml.example ]; then
    mv /etc/datadog-agent/conf.d/oom_kill.d/conf.yaml.example \
       /etc/datadog-agent/conf.d/oom_kill.d/conf.yaml.default
  fi
fi

# Match the key gpu_monitoring.enabled: true, allowing for other keys to be present below gpu_monitoring.
# regex breakdown:
# gpu_monitoring:\s*\n - match the gpu_monitoring parent key line
# (\s+.*\n)? - match any number of child keys indented under gpu_monitoring. Will stop the match if we find another parent key at the same level as gpu_monitoring
# \s+enabled\s*:\s*true - match the enabled: true key-value pair
# We use perl to read the whole file at once (-0777) and exit with 0 if the regex matches, 1 otherwise.
if perl -0777 -ne 'exit 0 if /gpu_monitoring:\s*\n(\s+.*\n)?\s+enabled\s*:\s*true/; exit 1' /etc/datadog-agent/system-probe.yaml || [[ "$DD_GPU_MONITORING_ENABLED" == "true" ]]; then
  if [ -f /etc/datadog-agent/conf.d/gpu.d/conf.yaml.example ]; then
    mv /etc/datadog-agent/conf.d/gpu.d/conf.yaml.example \
       /etc/datadog-agent/conf.d/gpu.d/conf.yaml.default
  fi
fi
