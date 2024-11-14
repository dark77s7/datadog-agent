// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package integrationslogs

import (
	_ "embed"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/e2e"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/environments"
	awshost "github.com/DataDog/datadog-agent/test/new-e2e/pkg/environments/aws/host"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/utils/e2e/client/agentclient"
	"github.com/DataDog/datadog-agent/test/new-e2e/tests/agent-metrics-logs/log-agent/utils"
	"github.com/DataDog/test-infra-definitions/components/datadog/agentparams"
)

type IntegrationsLogsSuite struct {
	e2e.BaseSuite[environments.Host]
}

//go:embed fixtures/tenLogs.py
var writeTenLogsCheck string

//go:embed fixtures/tenLogs.yaml
var writeTenLogsConfig string

//go:embed fixtures/maxSize.py
var maxSizeCheck string

//go:embed fixtures/maxSize.yaml
var maxSizeConfig string

// TestLinuxFakeIntakeSuite
func TestIntegrationsLogsSuite(t *testing.T) {
	suiteParams := []e2e.SuiteOption{
		e2e.WithProvisioner(awshost.Provisioner(awshost.WithAgentOptions(
			agentparams.WithLogs(),
			// set the integration log file max size to 1MB
			agentparams.WithAgentConfig("logs_config.integrations_logs_files_max_size: 2"),
			agentparams.WithFile("/etc/datadog-agent/checks.d/writeTenLogs.py", writeTenLogsCheck, true),
			agentparams.WithFile("/etc/datadog-agent/conf.d/writeTenLogs.yaml", writeTenLogsConfig, true),
			agentparams.WithFile("/etc/datadog-agent/checks.d/maxSize.py", maxSizeCheck, true),
			agentparams.WithFile("/etc/datadog-agent/conf.d/maxSize.yaml", maxSizeConfig, true))))}

	e2e.Run(t, &IntegrationsLogsSuite{}, suiteParams...)
}

// TestWriteTenLogsCheck ensures a check that logs are written to the file ten
// logs at a time
func (v *IntegrationsLogsSuite) TestWriteTenLogsCheck() {
	writeTenLogs := v.Env().Agent.Client.Check(agentclient.WithArgs([]string{"writeTenLogs"}))
	assert.Contains(v.T(), writeTenLogs, "writeTenLogs")

	v.EventuallyWithT(func(c *assert.CollectT) {
		output := v.Env().RemoteHost.MustExecute("sudo cat /opt/datadog-agent/run/integrations/writeTenLogs*.log")
		newLineCount := strings.Count(output, "\n")
		assert.Equal(c, newLineCount%10, 0)
		assert.GreaterOrEqual(c, newLineCount, 10)
	}, 1*time.Minute, 5*time.Second)

	utils.CheckLogsExpected(v.T(), v.Env().FakeIntake, "ten_logs_service", "Custom log message", []string{"env:dev", "bar:foo"})
}

// TestIntegrationLogFileMaxSize ensures integration log files don't exceed the max file size
func (v *IntegrationsLogsSuite) TestIntegrationLogFileMaxSize() {
	maxSizeCheck := v.Env().Agent.Client.Check(agentclient.WithArgs([]string{"maxSize"}))
	assert.Contains(v.T(), maxSizeCheck, "maxSize")

	// Since it's not yet possible to write to the integration log file by calling
	// the agent check command, we can test if the file size limits are being
	// respected using the following method:
	// 1. Wait until the integration log file reaches the maximum allowable size
	// (it won't be deleted until the next log that exceeds the maximum allowable
	// size is written).
	// 2. Check immediately that on the subsequent write, the file is smaller than
	// in step 1, indicating the log file has been deleted and remade, and thus
	// respects the set size.
	v.EventuallyWithT(func(c *assert.CollectT) {
		output := v.Env().RemoteHost.MustExecute("sudo cat /opt/datadog-agent/run/integrations/maxSize*.log")
		integrationLogFileSize := len(output)
		assert.Equal(c, 2*1024*1024, integrationLogFileSize)
	}, 1*time.Minute, 5*time.Second)

	v.EventuallyWithT(func(c *assert.CollectT) {
		output := v.Env().RemoteHost.MustExecute("sudo cat /opt/datadog-agent/run/integrations/maxSize*.log")
		integrationLogFileSize := len(output)
		assert.Equal(c, 1*1024*1024, integrationLogFileSize)
	}, 1*time.Minute, 5*time.Second)
}
