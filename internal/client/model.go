package client

type NodeMetrics struct {
	Name                   string
	CPUCores               float64
	MemoryBytes            int64
	AllocatableCPUCores    float64
	AllocatableMemoryBytes int64
}

type PodMetrics struct {
	Namespace   string
	Name        string
	CPUCores    float64
	MemoryBytes int64
}

type ContainerMetrics struct {
	Namespace   string
	Pod         string
	Name        string
	CPUCores    float64
	MemoryBytes int64
}

type PodAndContainerMetricsList struct {
	PodMetricsList       []PodMetrics
	ContainerMetricsList []ContainerMetrics
}
