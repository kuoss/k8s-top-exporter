package client

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsapi "k8s.io/metrics/pkg/apis/metrics"
)

func TestNodeMetrics(t *testing.T) {
	t.Parallel()

	metrics := []metricsapi.NodeMetrics{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "node-a"},
			Usage: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("250m"),
				v1.ResourceMemory: resource.MustParse("2Gi"),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "node-b"},
			Usage: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("4"),
				v1.ResourceMemory: resource.MustParse("512Mi"),
			},
		},
	}

	available := map[string]v1.ResourceList{
		"node-a": {
			v1.ResourceCPU:    resource.MustParse("4"),
			v1.ResourceMemory: resource.MustParse("8Gi"),
		},
		"node-b": {
			v1.ResourceCPU:    resource.MustParse("8"),
			v1.ResourceMemory: resource.MustParse("16Gi"),
		},
	}

	got := nodeMetrics(metrics, available)
	want := []NodeMetrics{
		{
			Name:                   "node-a",
			CPUCores:               0.25,
			MemoryBytes:            int64(2 * 1024 * 1024 * 1024),
			AllocatableCPUCores:    4,
			AllocatableMemoryBytes: int64(8 * 1024 * 1024 * 1024),
		},
		{
			Name:                   "node-b",
			CPUCores:               4,
			MemoryBytes:            int64(512 * 1024 * 1024),
			AllocatableCPUCores:    8,
			AllocatableMemoryBytes: int64(16 * 1024 * 1024 * 1024),
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("nodeMetrics() = %#v, want %#v", got, want)
	}
}
