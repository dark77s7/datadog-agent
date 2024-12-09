// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package examples

import (
	"context"
	awskubernetes "github.com/DataDog/datadog-agent/test/new-e2e/pkg/environments/aws/kubernetes"
	"github.com/DataDog/test-infra-definitions/components/datadog/agentwithoperatorparams"
	"github.com/DataDog/test-infra-definitions/components/datadog/operatorparams"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"testing"
	"time"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/e2e"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/environments"
	"github.com/stretchr/testify/assert"
)

type kindOperatorSuite struct {
	e2e.BaseSuite[environments.Kubernetes]
}

func TestKindOperatorSuite(t *testing.T) {
	customDDA := agentwithoperatorparams.DDAConfig{
		Name: "apm-enabled",
		YamlConfig: `
apiVersion: datadoghq.com/v2alpha1
kind: DatadogAgent
spec:
  global:
    kubelet:
      tlsVerify: false
  features:
    apm:
      enabled: true
`,
	}

	e2e.Run(t, &kindOperatorSuite{}, e2e.WithProvisioner(awskubernetes.KindProvisioner(
		awskubernetes.WithOperator(),
		awskubernetes.WithOperatorOptions([]operatorparams.Option{
			operatorparams.WithNamespace("datadog"),
		}...),
		awskubernetes.WithOperatorDDAOptions([]agentwithoperatorparams.Option{
			agentwithoperatorparams.WithNamespace("datadog"),
			agentwithoperatorparams.WithDDAConfig(customDDA),
		}...),
	)),
		e2e.WithDevMode(), e2e.WithSkipDeleteOnFailure())
}

func (v *kindOperatorSuite) TestClusterAgentInstalled() {
	v.EventuallyWithT(func(c *assert.CollectT) {
		res, _ := v.Env().KubernetesCluster.Client().CoreV1().Pods("datadog").List(context.TODO(), v1.ListOptions{
			FieldSelector: "status.phase=Running",
		})
		containsClusterAgent := false
		for _, pod := range res.Items {
			if strings.Contains(pod.Name, "cluster-agent") {
				containsClusterAgent = true
				break
			}
		}
		assert.True(v.T(), containsClusterAgent, "Cluster Agent not found")
	}, 30*time.Second, 15*time.Second)
}
