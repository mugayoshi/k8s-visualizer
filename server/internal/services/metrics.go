// internal/services/metrics.go
package services

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterMetrics represents overall cluster metrics
type ClusterMetrics struct {
	TotalNodes       int                `json:"total_nodes"`
	TotalPods        int                `json:"total_pods"`
	TotalNamespaces  int                `json:"total_namespaces"`
	CPUCapacity      string             `json:"cpu_capacity"`
	MemoryCapacity   string             `json:"memory_capacity"`
	CPUUsage         string             `json:"cpu_usage,omitempty"`
	MemoryUsage      string             `json:"memory_usage,omitempty"`
	NodeMetrics      []NodeMetrics      `json:"node_metrics"`
	NamespaceMetrics []NamespaceMetrics `json:"namespace_metrics,omitempty"`
}

// NodeMetrics represents metrics for a single node
type NodeMetrics struct {
	Name              string            `json:"name"`
	CPUCapacity       string            `json:"cpu_capacity"`
	MemoryCapacity    string            `json:"memory_capacity"`
	CPUAllocatable    string            `json:"cpu_allocatable"`
	MemoryAllocatable string            `json:"memory_allocatable"`
	PodCapacity       string            `json:"pod_capacity"`
	PodCount          int               `json:"pod_count"`
	Status            string            `json:"status"`
	Labels            map[string]string `json:"labels"`
}

// NamespaceMetrics represents metrics for a namespace
type NamespaceMetrics struct {
	Name        string `json:"name"`
	PodCount    int    `json:"pod_count"`
	CPUUsage    string `json:"cpu_usage,omitempty"`
	MemoryUsage string `json:"memory_usage,omitempty"`
}

// PodMetrics represents metrics for a single pod
type PodMetrics struct {
	Name            string `json:"name"`
	Namespace       string `json:"namespace"`
	Status          string `json:"status"`
	StatusDetail    string `json:"status_detail,omitempty"`
	ReadyContainers string `json:"ready_containers"`
	CPURequest      string `json:"cpu_request,omitempty"`
	CPULimit        string `json:"cpu_limit,omitempty"`
	MemoryRequest   string `json:"memory_request,omitempty"`
	MemoryLimit     string `json:"memory_limit,omitempty"`
	ContainerCount  int    `json:"container_count"`
	Age             string `json:"age,omitempty"`
	RestartCount    int32  `json:"restart_count"`
}

// GetClusterMetrics retrieves overall cluster metrics
func (k *K8sClient) GetClusterMetrics(ctx context.Context) (*ClusterMetrics, error) {
	clientset := k.GetClientset()

	// Get nodes
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	// Get all pods
	pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	// Get namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	// Calculate total capacity
	totalCPU := resource.NewQuantity(0, resource.DecimalSI)
	totalMemory := resource.NewQuantity(0, resource.BinarySI)
	nodeMetrics := make([]NodeMetrics, 0, len(nodes.Items))

	for _, node := range nodes.Items {
		// Add to totals
		cpu := node.Status.Capacity.Cpu()
		memory := node.Status.Capacity.Memory()
		totalCPU.Add(*cpu)
		totalMemory.Add(*memory)

		// Count pods on this node
		podCount := 0
		for _, pod := range pods.Items {
			if pod.Spec.NodeName == node.Name {
				podCount++
			}
		}

		// Get node status
		status := "Unknown"
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady {
				if condition.Status == corev1.ConditionTrue {
					status = "Ready"
				} else {
					status = "NotReady"
				}
				break
			}
		}

		nodeMetrics = append(nodeMetrics, NodeMetrics{
			Name:              node.Name,
			CPUCapacity:       node.Status.Capacity.Cpu().String(),
			MemoryCapacity:    node.Status.Capacity.Memory().String(),
			CPUAllocatable:    node.Status.Allocatable.Cpu().String(),
			MemoryAllocatable: node.Status.Allocatable.Memory().String(),
			PodCapacity:       node.Status.Capacity.Pods().String(),
			PodCount:          podCount,
			Status:            status,
			Labels:            node.Labels,
		})
	}

	// Calculate namespace metrics
	namespaceMetrics := make([]NamespaceMetrics, 0, len(namespaces.Items))
	namespacePodCount := make(map[string]int)

	for _, pod := range pods.Items {
		namespacePodCount[pod.Namespace]++
	}

	for _, ns := range namespaces.Items {
		namespaceMetrics = append(namespaceMetrics, NamespaceMetrics{
			Name:     ns.Name,
			PodCount: namespacePodCount[ns.Name],
		})
	}

	return &ClusterMetrics{
		TotalNodes:       len(nodes.Items),
		TotalPods:        len(pods.Items),
		TotalNamespaces:  len(namespaces.Items),
		CPUCapacity:      totalCPU.String(),
		MemoryCapacity:   totalMemory.String(),
		NodeMetrics:      nodeMetrics,
		NamespaceMetrics: namespaceMetrics,
	}, nil
}

// GetNodeMetrics retrieves metrics for a specific node
func (k *K8sClient) GetNodeMetrics(ctx context.Context, nodeName string) (*NodeMetrics, error) {
	clientset := k.GetClientset()

	node, err := clientset.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	// Count pods on this node
	pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	status := "Unknown"
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			if condition.Status == corev1.ConditionTrue {
				status = "Ready"
			} else {
				status = "NotReady"
			}
			break
		}
	}

	return &NodeMetrics{
		Name:              node.Name,
		CPUCapacity:       node.Status.Capacity.Cpu().String(),
		MemoryCapacity:    node.Status.Capacity.Memory().String(),
		CPUAllocatable:    node.Status.Allocatable.Cpu().String(),
		MemoryAllocatable: node.Status.Allocatable.Memory().String(),
		PodCapacity:       node.Status.Capacity.Pods().String(),
		PodCount:          len(pods.Items),
		Status:            status,
		Labels:            node.Labels,
	}, nil
}

// GetNamespaceMetrics retrieves metrics for a specific namespace
func (k *K8sClient) GetNamespaceMetrics(ctx context.Context, namespace string) (*NamespaceMetrics, error) {
	clientset := k.GetClientset()

	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	return &NamespaceMetrics{
		Name:     namespace,
		PodCount: len(pods.Items),
	}, nil
}

// GetPodMetrics retrieves metrics for a specific pod
// Note: This requires metrics-server to be installed in the cluster
func (k *K8sClient) GetPodMetrics(ctx context.Context, namespace, podName string) (*PodMetrics, error) {
	clientset := k.GetClientset()

	// Try to get metrics from metrics API
	// First, check if the pod exists
	pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod: %w", err)
	}

	metrics := &PodMetrics{
		Name:           pod.Name,
		Namespace:      pod.Namespace,
		ContainerCount: len(pod.Spec.Containers),
	}

	// Calculate total restarts
	for _, containerStatus := range pod.Status.ContainerStatuses {
		metrics.RestartCount += containerStatus.RestartCount
	}

	// Calculate total CPU and memory requests/limits
	var totalCPURequest, totalCPULimit, totalMemoryRequest, totalMemoryLimit resource.Quantity

	for _, container := range pod.Spec.Containers {
		if cpuReq := container.Resources.Requests.Cpu(); cpuReq != nil {
			totalCPURequest.Add(*cpuReq)
		}
		if cpuLim := container.Resources.Limits.Cpu(); cpuLim != nil {
			totalCPULimit.Add(*cpuLim)
		}
		if memReq := container.Resources.Requests.Memory(); memReq != nil {
			totalMemoryRequest.Add(*memReq)
		}
		if memLim := container.Resources.Limits.Memory(); memLim != nil {
			totalMemoryLimit.Add(*memLim)
		}
	}

	// Set metrics if values exist
	if !totalCPURequest.IsZero() {
		metrics.CPURequest = totalCPURequest.String()
	}
	if !totalCPULimit.IsZero() {
		metrics.CPULimit = totalCPULimit.String()
	}
	if !totalMemoryRequest.IsZero() {
		metrics.MemoryRequest = totalMemoryRequest.String()
	}
	if !totalMemoryLimit.IsZero() {
		metrics.MemoryLimit = totalMemoryLimit.String()
	}

	// Calculate pod age
	if !pod.CreationTimestamp.IsZero() {
		age := metav1.Now().Sub(pod.CreationTimestamp.Time)
		metrics.Age = formatAge(age)
	}

	return metrics, nil
}

// formatAge formats a duration into a human-readable age string
func formatAge(d time.Duration) string {
	if d < 0 {
		return "0s"
	}

	totalSeconds := int64(d.Seconds())

	if totalSeconds < 60 {
		return fmt.Sprintf("%ds", totalSeconds)
	} else if totalSeconds < 3600 {
		return fmt.Sprintf("%dm", totalSeconds/60)
	} else if totalSeconds < 86400 {
		return fmt.Sprintf("%dh", totalSeconds/3600)
	} else {
		return fmt.Sprintf("%dd", totalSeconds/86400)
	}
}
