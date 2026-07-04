package client

import (
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	metricsapi "k8s.io/metrics/pkg/apis/metrics"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

type Client struct {
	metricsClient metricsclientset.Interface
	nodeClient    corev1client.CoreV1Interface
}

func NewClient() (*Client, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	metricsClient, err := metricsclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return NewClientWith(clientset, metricsClient)
}

func NewClientWith(clientset kubernetes.Interface, metricsClient metricsclientset.Interface) (*Client, error) {
	apiGroups, err := clientset.Discovery().ServerGroups()
	if err != nil {
		return nil, err
	}
	if !SupportedMetricsAPIVersionAvailable(apiGroups) {
		return nil, errors.New("Metrics API not available")
	}

	return &Client{
		nodeClient:    clientset.CoreV1(),
		metricsClient: metricsClient,
	}, nil
}

func SupportedMetricsAPIVersionAvailable(discoveredAPIGroups *metav1.APIGroupList) bool {
	supportedMetricsAPIVersions := []string{"v1beta1"}
	for _, discoveredAPIGroup := range discoveredAPIGroups.Groups {
		if discoveredAPIGroup.Name != metricsapi.GroupName {
			continue
		}
		for _, version := range discoveredAPIGroup.Versions {
			for _, supportedVersion := range supportedMetricsAPIVersions {
				if version.Version == supportedVersion {
					return true
				}
			}
		}
	}
	return false
}
