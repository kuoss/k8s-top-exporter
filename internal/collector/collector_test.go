package collector

import (
	"testing"

	topclient "github.com/jmnote/k8s-top-exporter/internal/client"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type stubTopClient struct {
	nodeMetrics []topclient.NodeMetrics
	podMetrics  *topclient.PodAndContainerMetricsList
}

func (s stubTopClient) NodeMetricsList() ([]topclient.NodeMetrics, error) {
	return s.nodeMetrics, nil
}

func (s stubTopClient) PodAndContainerMetricsList() (*topclient.PodAndContainerMetricsList, error) {
	return s.podMetrics, nil
}

func TestCollector(t *testing.T) {
	t.Parallel()

	c := &Collector{
		topclient: stubTopClient{
			nodeMetrics: []topclient.NodeMetrics{
				{
					Name:                   "node-a",
					CPUCores:               1.25,
					MemoryBytes:            1024,
					AllocatableCPUCores:    2.5,
					AllocatableMemoryBytes: 4096,
				},
			},
			podMetrics: &topclient.PodAndContainerMetricsList{
				PodMetricsList: []topclient.PodMetrics{
					{
						Namespace:   "default",
						Name:        "pod-a",
						CPUCores:    0.75,
						MemoryBytes: 2048,
					},
				},
				ContainerMetricsList: []topclient.ContainerMetrics{
					{
						Namespace:   "default",
						Pod:         "pod-a",
						Name:        "c1",
						CPUCores:    0.25,
						MemoryBytes: 512,
					},
				},
			},
		},
		nodeCPUCoresDesc:               prometheus.NewDesc("k8s_top_node_cpu_cores", "CPU usage of the node in cores.", []string{"name"}, nil),
		nodeMemoryBytesDesc:            prometheus.NewDesc("k8s_top_node_memory_bytes", "Memory usage of the node in bytes.", []string{"name"}, nil),
		nodeAllocatableCPUCoresDesc:    prometheus.NewDesc("k8s_top_node_allocatable_cpu_cores", "Allocatable CPU of the node in cores.", []string{"name"}, nil),
		nodeAllocatableMemoryBytesDesc: prometheus.NewDesc("k8s_top_node_allocatable_memory_bytes", "Allocatable memory of the node in bytes.", []string{"name"}, nil),
		podCPUCoresDesc:                prometheus.NewDesc("k8s_top_pod_cpu_cores", "CPU usage of the pod in cores.", []string{"namespace", "name"}, nil),
		podMemoryBytesDesc:             prometheus.NewDesc("k8s_top_pod_memory_bytes", "Memory usage of the pod in bytes.", []string{"namespace", "name"}, nil),
		containerCPUCoresDesc:          prometheus.NewDesc("k8s_top_container_cpu_cores", "CPU usage of the container in cores.", []string{"namespace", "pod", "name"}, nil),
		containerMemoryBytesDesc:       prometheus.NewDesc("k8s_top_container_memory_bytes", "Memory usage of the container in bytes.", []string{"namespace", "pod", "name"}, nil),
	}

	reg := prometheus.NewRegistry()
	if err := reg.Register(c); err != nil {
		t.Fatalf("Register() error: %v", err)
	}

	families, err := reg.Gather()
	if err != nil {
		t.Fatalf("Gather() error: %v", err)
	}

	assertGauge(t, families, "k8s_top_node_cpu_cores", map[string]string{"name": "node-a"}, 1.25)
	assertGauge(t, families, "k8s_top_node_memory_bytes", map[string]string{"name": "node-a"}, 1024)
	assertGauge(t, families, "k8s_top_node_allocatable_cpu_cores", map[string]string{"name": "node-a"}, 2.5)
	assertGauge(t, families, "k8s_top_node_allocatable_memory_bytes", map[string]string{"name": "node-a"}, 4096)
	assertGauge(t, families, "k8s_top_pod_cpu_cores", map[string]string{"namespace": "default", "name": "pod-a"}, 0.75)
	assertGauge(t, families, "k8s_top_pod_memory_bytes", map[string]string{"namespace": "default", "name": "pod-a"}, 2048)
	assertGauge(t, families, "k8s_top_container_cpu_cores", map[string]string{"namespace": "default", "pod": "pod-a", "name": "c1"}, 0.25)
	assertGauge(t, families, "k8s_top_container_memory_bytes", map[string]string{"namespace": "default", "pod": "pod-a", "name": "c1"}, 512)
}

func assertGauge(t *testing.T, families []*dto.MetricFamily, name string, labels map[string]string, want float64) {
	t.Helper()

	for _, family := range families {
		if family.GetName() != name {
			continue
		}
		for _, metric := range family.GetMetric() {
			if hasLabels(metric, labels) {
				if got := metric.GetGauge().GetValue(); got != want {
					t.Fatalf("%s labels=%v = %v, want %v", name, labels, got, want)
				}
				return
			}
		}
		t.Fatalf("metric %s with labels %v not found", name, labels)
	}

	t.Fatalf("metric family %s not found", name)
}

func hasLabels(metric *dto.Metric, want map[string]string) bool {
	if len(metric.GetLabel()) != len(want) {
		return false
	}
	for _, label := range metric.GetLabel() {
		if want[label.GetName()] != label.GetValue() {
			return false
		}
	}
	return true
}
