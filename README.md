# k8s-top-exporter

`k8s-top-exporter` is a Prometheus exporter that collects resource usage data from the Kubernetes Metrics API, equivalent to the data shown by the `kubectl top` command.

It queries the Kubernetes Metrics Server to provide real-time CPU and memory usage data for Nodes, Pods, and Containers.

## Key Features

- **Standardized Base Units**: All metrics are collected and exposed in standard base units in accordance with Prometheus guidelines:
  - CPU: `cores` (float64)
  - Memory: `bytes` (int64)
- **Container Security**: Employs multi-stage builds and uses `gcr.io/distroless/static:nonroot` as the base image to minimize container vulnerability exposure.
- **Granular Monitoring**: Collects resource usage data at the Node, Pod, and Container levels.
- **Raw Metric Exposure**: Instead of exporting pre-calculated percentage values, it exposes raw usage metrics alongside allocatable resource totals. This allows users to calculate percentages and free resources flexibly using PromQL.

## Exported Metrics

### 1. Node Metrics

| Metric Name | Type | Description |
| :--- | :--- | :--- |
| `k8s_top_node_cpu_cores` | Gauge | Current CPU usage of the node (unit: Cores) |
| `k8s_top_node_memory_bytes` | Gauge | Current memory usage of the node (unit: Bytes) |
| `k8s_top_node_allocatable_cpu_cores` | Gauge | Total allocatable CPU of the node (unit: Cores) |
| `k8s_top_node_allocatable_memory_bytes` | Gauge | Total allocatable memory of the node (unit: Bytes) |

### 2. Pod Metrics

| Metric Name | Type | Labels | Description |
| :--- | :--- | :--- | :--- |
| `k8s_top_pod_cpu_cores` | Gauge | `namespace`, `name` | Total current CPU usage of the pod (unit: Cores) |
| `k8s_top_pod_memory_bytes` | Gauge | `namespace`, `name` | Total current memory usage of the pod (unit: Bytes) |

### 3. Container Metrics

| Metric Name | Type | Labels | Description |
| :--- | :--- | :--- | :--- |
| `k8s_top_pod_container_cpu_cores` | Gauge | `namespace`, `pod`, `name` | Current CPU usage of the individual container (unit: Cores) |
| `k8s_top_pod_container_memory_bytes` | Gauge | `namespace`, `pod`, `name` | Current memory usage of the individual container (unit: Bytes) |

## PromQL Examples

Since metrics are exported in base units, you can write PromQL queries to compute resource usage percentages or free capacity as needed.

### Calculate Node Resource Usage Percentage (%)
```promql
# CPU usage percentage per node
k8s_top_node_cpu_cores / k8s_top_node_allocatable_cpu_cores * 100

# Memory usage percentage per node
k8s_top_node_memory_bytes / k8s_top_node_allocatable_memory_bytes * 100
```

### Calculate Node Free Resource Capacity
```promql
# Free memory in bytes per node
k8s_top_node_allocatable_memory_bytes - k8s_top_node_memory_bytes
```

### Query Top Pods in a Namespace
```promql
# Top 5 pods using the most memory in the "production" namespace
topk(5, k8s_top_pod_memory_bytes{namespace="production"})
```

## Deployment

### Prerequisites
- The Metrics Server must be installed and running in your Kubernetes cluster.

### Deploying the Exporter
Apply all resources in the `deploy/` directory, which includes the ServiceAccount, RBAC permissions, Deployment, and Service:

```bash
kubectl apply -f deploy/
```

For local development, run the binary from the command package:

```bash
go run ./cmd/k8s-top-exporter
```

## License

This project is licensed under the [Apache License 2.0](LICENSE).
