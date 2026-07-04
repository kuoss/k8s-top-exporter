package collector

import (
	"log"

	topclient "github.com/jmnote/k8s-top-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type collector struct {
	topclient *topclient.Client

	nodeCPUCoresDesc               *prometheus.Desc
	nodeMemoryBytesDesc            *prometheus.Desc
	nodeAllocatableCPUCoresDesc    *prometheus.Desc
	nodeAllocatableMemoryBytesDesc *prometheus.Desc

	podCPUCoresDesc          *prometheus.Desc
	podMemoryBytesDesc       *prometheus.Desc
	containerCPUCoresDesc    *prometheus.Desc
	containerMemoryBytesDesc *prometheus.Desc
}

func NewCollector() (*collector, error) {
	topclient, err := topclient.NewClient()
	if err != nil {
		return nil, err
	}
	return &collector{
		topclient: topclient,

		nodeCPUCoresDesc:               prometheus.NewDesc("k8s_top_node_cpu_cores", "CPU usage of the node in cores.", []string{"name"}, nil),
		nodeMemoryBytesDesc:            prometheus.NewDesc("k8s_top_node_memory_bytes", "Memory usage of the node in bytes.", []string{"name"}, nil),
		nodeAllocatableCPUCoresDesc:    prometheus.NewDesc("k8s_top_node_allocatable_cpu_cores", "Allocatable CPU of the node in cores.", []string{"name"}, nil),
		nodeAllocatableMemoryBytesDesc: prometheus.NewDesc("k8s_top_node_allocatable_memory_bytes", "Allocatable memory of the node in bytes.", []string{"name"}, nil),

		podCPUCoresDesc:          prometheus.NewDesc("k8s_top_pod_cpu_cores", "CPU usage of the pod in cores.", []string{"namespace", "name"}, nil),
		podMemoryBytesDesc:       prometheus.NewDesc("k8s_top_pod_memory_bytes", "Memory usage of the pod in bytes.", []string{"namespace", "name"}, nil),
		containerCPUCoresDesc:    prometheus.NewDesc("k8s_top_pod_container_cpu_cores", "CPU usage of the container in cores.", []string{"namespace", "pod", "name"}, nil),
		containerMemoryBytesDesc: prometheus.NewDesc("k8s_top_pod_container_memory_bytes", "Memory usage of the container in bytes.", []string{"namespace", "pod", "name"}, nil),
	}, nil
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	c.collectNodeMetrics(ch)
	c.collectPodAndContainerMetrics(ch)
}

func (c *collector) collectNodeMetrics(ch chan<- prometheus.Metric) {
	nodeMetricsList, err := c.topclient.GetNodeMetricsList()
	if err != nil {
		log.Println(err)
		return
	}
	for _, m := range nodeMetricsList {
		ch <- prometheus.MustNewConstMetric(c.nodeCPUCoresDesc, prometheus.GaugeValue, m.CPUCores, []string{m.Name}...)
		ch <- prometheus.MustNewConstMetric(c.nodeMemoryBytesDesc, prometheus.GaugeValue, float64(m.MemoryBytes), []string{m.Name}...)
		ch <- prometheus.MustNewConstMetric(c.nodeAllocatableCPUCoresDesc, prometheus.GaugeValue, m.AllocatableCPUCores, []string{m.Name}...)
		ch <- prometheus.MustNewConstMetric(c.nodeAllocatableMemoryBytesDesc, prometheus.GaugeValue, float64(m.AllocatableMemoryBytes), []string{m.Name}...)
	}
}

func (c *collector) collectPodAndContainerMetrics(ch chan<- prometheus.Metric) {
	podAndContainerMetricsList, err := c.topclient.GetPodAndContainerMetricsList()
	if err != nil {
		log.Println(err)
		return
	}
	for _, m := range podAndContainerMetricsList.PodMetricsList {
		ch <- prometheus.MustNewConstMetric(c.podCPUCoresDesc, prometheus.GaugeValue, m.CPUCores, []string{m.Namespace, m.Name}...)
		ch <- prometheus.MustNewConstMetric(c.podMemoryBytesDesc, prometheus.GaugeValue, float64(m.MemoryBytes), []string{m.Namespace, m.Name}...)
	}
	for _, m := range podAndContainerMetricsList.ContainerMetricsList {
		ch <- prometheus.MustNewConstMetric(c.containerCPUCoresDesc, prometheus.GaugeValue, m.CPUCores, []string{m.Namespace, m.Pod, m.Name}...)
		ch <- prometheus.MustNewConstMetric(c.containerMemoryBytesDesc, prometheus.GaugeValue, float64(m.MemoryBytes), []string{m.Namespace, m.Pod, m.Name}...)
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.nodeCPUCoresDesc
	ch <- c.nodeMemoryBytesDesc
	ch <- c.nodeAllocatableCPUCoresDesc
	ch <- c.nodeAllocatableMemoryBytesDesc

	ch <- c.podCPUCoresDesc
	ch <- c.podMemoryBytesDesc
	ch <- c.containerCPUCoresDesc
	ch <- c.containerMemoryBytesDesc
}

