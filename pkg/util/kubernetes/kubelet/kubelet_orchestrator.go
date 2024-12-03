// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build kubelet && orchestrator

package kubelet

import (
	"context"

	v1 "k8s.io/api/core/v1"
	kubeletv1alpha1 "k8s.io/kubelet/pkg/apis/stats/v1alpha1"
)

// KubeUtilInterface defines the interface for kubelet api
// and includes extra functions for the orchestrator build flag
type KubeUtilInterface interface {
	GetNodeInfo(ctx context.Context) (string, string, error)
	GetNodename(ctx context.Context) (string, error)
	GetLocalPodList(ctx context.Context) ([]*Pod, error)
	GetLocalPodListWithMetadata(ctx context.Context) (*PodList, error)
	ForceGetLocalPodList(ctx context.Context) (*PodList, error)
	GetPodForContainerID(ctx context.Context, containerID string) (*Pod, error)
	QueryKubelet(ctx context.Context, path string) ([]byte, int, error)
	GetRawConnectionInfo() map[string]string
	GetRawMetrics(ctx context.Context) ([]byte, error)
	GetRawLocalPodList(ctx context.Context) ([]*v1.Pod, error)
	GetLocalStatsSummary(ctx context.Context) (*kubeletv1alpha1.Summary, error)
}
