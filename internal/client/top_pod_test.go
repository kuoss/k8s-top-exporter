package client

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsapi "k8s.io/metrics/pkg/apis/metrics"
)

func TestGetPodMetrics(t *testing.T) {
	t.Parallel()

	metrics := []metricsapi.PodMetrics{
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "pod-a",
			},
			Containers: []metricsapi.ContainerMetrics{
				{
					Name: "c1",
					Usage: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("150m"),
						v1.ResourceMemory: resource.MustParse("128Mi"),
					},
				},
				{
					Name: "c2",
					Usage: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("350m"),
						v1.ResourceMemory: resource.MustParse("256Mi"),
					},
				},
			},
		},
	}

	got := (&Client{}).getPodMetrics(metrics)
	want := []PodMetrics{
		{
			Namespace:   "default",
			Name:        "pod-a",
			CPUCores:    0.5,
			MemoryBytes: int64(384 * 1024 * 1024),
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("getPodMetrics() = %#v, want %#v", got, want)
	}
}

func TestGetContainerMetrics(t *testing.T) {
	t.Parallel()

	metrics := []metricsapi.PodMetrics{
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "pod-a",
			},
			Containers: []metricsapi.ContainerMetrics{
				{
					Name: "c1",
					Usage: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("150m"),
						v1.ResourceMemory: resource.MustParse("128Mi"),
					},
				},
				{
					Name: "c2",
					Usage: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("350m"),
						v1.ResourceMemory: resource.MustParse("256Mi"),
					},
				},
			},
		},
	}

	got := (&Client{}).getContainerMetrics(metrics)
	want := []ContainerMetrics{
		{
			Namespace:   "default",
			Pod:         "pod-a",
			Name:        "c1",
			CPUCores:    0.15,
			MemoryBytes: int64(128 * 1024 * 1024),
		},
		{
			Namespace:   "default",
			Pod:         "pod-a",
			Name:        "c2",
			CPUCores:    0.35,
			MemoryBytes: int64(256 * 1024 * 1024),
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("getContainerMetrics() = %#v, want %#v", got, want)
	}
}
