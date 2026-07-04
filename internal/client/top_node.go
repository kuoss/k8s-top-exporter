package client

import (
	"context"
	"errors"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	metricsapi "k8s.io/metrics/pkg/apis/metrics"
	metricsv1beta1api "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func (c *Client) NodeMetricsList() ([]NodeMetrics, error) {
	metrics, err := c.nodeMetricsFromAPI()
	if err != nil {
		return nil, err
	}
	if len(metrics.Items) == 0 {
		return nil, errors.New("metrics not available yet")
	}
	var nodes []v1.Node
	nodeList, err := c.nodeClient.Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	nodes = append(nodes, nodeList.Items...)

	allocatable := make(map[string]v1.ResourceList)
	for _, n := range nodes {
		allocatable[n.Name] = n.Status.Allocatable
	}

	NodeMetricsList := nodeMetrics(metrics.Items, allocatable)
	return NodeMetricsList, nil
}

func (c *Client) nodeMetricsFromAPI() (*metricsapi.NodeMetricsList, error) {
	versionedMetrics, err := c.metricsClient.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	metrics := &metricsapi.NodeMetricsList{}
	err = metricsv1beta1api.Convert_v1beta1_NodeMetricsList_To_metrics_NodeMetricsList(versionedMetrics, metrics, nil)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func nodeMetrics(metrics []metricsapi.NodeMetrics, availableResources map[string]v1.ResourceList) []NodeMetrics {
	var NodeMetricsList []NodeMetrics
	for _, m := range metrics {
		available := availableResources[m.Name]
		cpuQuantity := m.Usage[v1.ResourceCPU]
		cpuAvailable := available[v1.ResourceCPU]
		memoryQuantity := m.Usage[v1.ResourceMemory]
		memoryAvailable := available[v1.ResourceMemory]
		NodeMetricsList = append(NodeMetricsList, NodeMetrics{
			Name:                   m.Name,
			CPUCores:               float64(cpuQuantity.MilliValue()) / 1000,
			MemoryBytes:            memoryQuantity.Value(),
			AllocatableCPUCores:    float64(cpuAvailable.MilliValue()) / 1000,
			AllocatableMemoryBytes: memoryAvailable.Value(),
		})
	}
	return NodeMetricsList
}
