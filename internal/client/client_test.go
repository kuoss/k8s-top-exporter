package client

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery/fake"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	metricsfake "k8s.io/metrics/pkg/client/clientset/versioned/fake"
)

func TestSupportedMetricsAPIVersionAvailable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		groups   *metav1.APIGroupList
		expected bool
	}{
		{
			name: "supported version present",
			groups: &metav1.APIGroupList{
				Groups: []metav1.APIGroup{
					{
						Name: "metrics.k8s.io",
						Versions: []metav1.GroupVersionForDiscovery{
							{Version: "v1beta1"},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "supported version absent",
			groups: &metav1.APIGroupList{
				Groups: []metav1.APIGroup{
					{
						Name: "metrics.k8s.io",
						Versions: []metav1.GroupVersionForDiscovery{
							{Version: "v1"},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "metrics group absent",
			groups: &metav1.APIGroupList{
				Groups: []metav1.APIGroup{
					{
						Name: "apps",
						Versions: []metav1.GroupVersionForDiscovery{
							{Version: "v1"},
						},
					},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := SupportedMetricsAPIVersionAvailable(tt.groups); got != tt.expected {
				t.Fatalf("SupportedMetricsAPIVersionAvailable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewClientWith(t *testing.T) {
	t.Parallel()

	kubeClient := k8sfake.NewSimpleClientset()
	discoveryClient := kubeClient.Discovery().(*fake.FakeDiscovery)
	discoveryClient.Resources = []*metav1.APIResourceList{
		{
			GroupVersion: "metrics.k8s.io/v1beta1",
		},
	}

	metricsClient := metricsfake.NewSimpleClientset()

	c, err := NewClientWith(kubeClient, metricsClient)
	if err != nil {
		t.Fatalf("NewClientWith() unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("NewClientWith() returned nil client")
	}
}

func TestNewClientWithRejectsMissingMetricsAPI(t *testing.T) {
	t.Parallel()

	kubeClient := k8sfake.NewSimpleClientset()
	discoveryClient := kubeClient.Discovery().(*fake.FakeDiscovery)
	discoveryClient.Resources = []*metav1.APIResourceList{
		{
			GroupVersion: "apps/v1",
		},
	}

	metricsClient := metricsfake.NewSimpleClientset()

	_, err := NewClientWith(kubeClient, metricsClient)
	if err == nil {
		t.Fatal("NewClientWith() error = nil, want error")
	}
}
