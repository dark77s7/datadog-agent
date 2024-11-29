// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2022-present Datadog, Inc.

package util

import (
	"context"
	"strings"

	"github.com/DataDog/datadog-agent/pkg/config/env"
	"github.com/DataDog/datadog-agent/pkg/config/model"
	pkgconfigsetup "github.com/DataDog/datadog-agent/pkg/config/setup"
	configUtils "github.com/DataDog/datadog-agent/pkg/config/utils"
	"github.com/DataDog/datadog-agent/pkg/util/fargate"
	"github.com/DataDog/datadog-agent/pkg/util/flavor"
	"github.com/DataDog/datadog-agent/pkg/util/kubernetes/clustername"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

// GetStaticTagsSlice gets the "static tags" for this agent.  These are tags
// that are attached directly to everything the agent produces, but _not_
// included in host tags.  In environments with no host metadata (such as where
// the hostname is empty), tags that would otherwise be included in host
// metadata are generated by this function.
func GetStaticTagsSlice(ctx context.Context) []string {
	// fargate (ECS or EKS) does not have host tags, so we need to
	// add static tags to each container manually

	if !fargate.IsFargateInstance() {
		return nil
	}

	tags := []string{}

	// DD_TAGS / DD_EXTRA_TAGS
	tags = append(tags, configUtils.GetConfiguredTags(pkgconfigsetup.Datadog(), false)...)

	// EKS Fargate specific tags
	if env.IsFeaturePresent(env.EKSFargate) {
		// eks_fargate_node
		node, err := fargate.GetEKSFargateNodename()
		if err != nil {
			log.Infof("Couldn't build the 'eks_fargate_node' tag: %v", err)
		} else {
			tags = append(tags, "eks_fargate_node:"+node)
		}

		// kube_cluster_name
		clusterTagNamePrefix := "kube_cluster_name:"
		var tag string
		var found bool
		for _, tag = range tags {
			if strings.HasPrefix(tag, clusterTagNamePrefix) {
				found = true
				break
			}
		}
		if found {
			log.Infof("'%s' was set manually via DD_TAGS, not changing it", clusterTagNamePrefix+tag)
		} else {
			cluster := clustername.GetClusterNameTagValue(ctx, "")
			if cluster == "" {
				log.Infof("Couldn't build the %q.. tag, DD_CLUSTER_NAME can be used to set it", clusterTagNamePrefix)
			} else {
				tags = append(tags, clusterTagNamePrefix+cluster)
			}
		}
	}

	return tags
}

// GetStaticTags is similar to GetStaticTagsSlice, but returning a map[string][]string containing
// <key>:<value> pairs for tags.  Tags not matching this pattern are omitted.
func GetStaticTags(ctx context.Context) map[string][]string {
	tags := GetStaticTagsSlice(ctx)
	if tags == nil {
		return nil
	}
	return sliceToMap(tags)
}

// GetGlobalEnvTags is similar to GetStaticTags, but returning a map[string][]string containing
// <key>:<value> pairs for all global environment tags on the cluster agent. This includes:
// DD_TAGS, DD_EXTRA_TAGS, DD_CLUSTER_CHECKS_EXTRA_TAGS, and DD_ORCHESTRATOR_EXPLORER_EXTRA_TAGS
func GetGlobalEnvTags(config model.Reader) map[string][]string {
	if flavor.GetFlavor() != flavor.ClusterAgent {
		return nil
	}

	// DD_TAGS / DD_EXTRA_TAGS
	tags := configUtils.GetConfiguredTags(config, false)

	// DD_CLUSTER_CHECKS_EXTRA_TAGS / DD_ORCHESTRATOR_EXPLORER_EXTRA_TAGS
	tags = append(tags, configUtils.GetConfiguredDCATags(config)...)

	if tags == nil {
		return nil
	}
	return sliceToMap(tags)
}

func sliceToMap(tags []string) map[string][]string {
	rv := make(map[string][]string, len(tags))
	for _, t := range tags {
		tagParts := strings.SplitN(t, ":", 2)
		if len(tagParts) == 2 {
			if _, ok := rv[tagParts[0]]; !ok {
				rv[tagParts[0]] = []string{}
			}
			rv[tagParts[0]] = append(rv[tagParts[0]], tagParts[1])
		}
	}
	return rv
}
