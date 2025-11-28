// internal/models/types.go
package models

import "time"

// NodeResponse represents a simplified node response
type NodeResponse struct {
	Name              string            `json:"name"`
	Status            string            `json:"status"`
	Created           time.Time         `json:"created"`
	CPUCapacity       string            `json:"cpu_capacity"`
	MemoryCapacity    string            `json:"memory_capacity"`
	PodCapacity       string            `json:"pod_capacity"`
	CPUAllocatable    string            `json:"cpu_allocatable"`
	MemoryAllocatable string            `json:"memory_allocatable"`
	Labels            map[string]string `json:"labels"`
	Addresses         []NodeAddress     `json:"addresses"`
}

// NodeAddress represents a node address
type NodeAddress struct {
	Type    string `json:"type"`
	Address string `json:"address"`
}

// PodResponse represents a simplified pod response
type PodResponse struct {
	Name            string            `json:"name"`
	Namespace       string            `json:"namespace"`
	Status          string            `json:"status"`
	StatusDetail    string            `json:"status_detail,omitempty"`
	Phase           string            `json:"phase"`
	ReadyContainers string            `json:"ready_containers"`
	Node            string            `json:"node"`
	Created         time.Time         `json:"created"`
	Labels          map[string]string `json:"labels"`
	Containers      []ContainerInfo   `json:"containers"`
	RestartCount    int32             `json:"restart_count"`
	CPUUsage        string            `json:"cpu_usage,omitempty"`
	MemoryUsage     string            `json:"memory_usage,omitempty"`
}

// ContainerInfo represents container information
type ContainerInfo struct {
	Name         string `json:"name"`
	Image        string `json:"image"`
	Ready        bool   `json:"ready"`
	RestartCount int32  `json:"restart_count"`
	State        string `json:"state"`
}

// NamespaceResponse represents a simplified namespace response
type NamespaceResponse struct {
	Name    string            `json:"name"`
	Status  string            `json:"status"`
	Created time.Time         `json:"created"`
	Labels  map[string]string `json:"labels"`
}

// DeploymentResponse represents a simplified deployment response
type DeploymentResponse struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	Replicas          int32             `json:"replicas"`
	ReadyReplicas     int32             `json:"ready_replicas"`
	AvailableReplicas int32             `json:"available_replicas"`
	UpdatedReplicas   int32             `json:"updated_replicas"`
	Created           time.Time         `json:"created"`
	Labels            map[string]string `json:"labels"`
}

// ServiceResponse represents a simplified service response
type ServiceResponse struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Type      string            `json:"type"`
	ClusterIP string            `json:"cluster_ip"`
	Ports     []ServicePort     `json:"ports"`
	Created   time.Time         `json:"created"`
	Labels    map[string]string `json:"labels"`
}

// ServicePort represents a service port
type ServicePort struct {
	Name       string `json:"name"`
	Protocol   string `json:"protocol"`
	Port       int32  `json:"port"`
	TargetPort string `json:"target_port"`
	NodePort   int32  `json:"node_port,omitempty"`
}

// EventResponse represents a Kubernetes event
type EventResponse struct {
	Type      string    `json:"type"`
	Reason    string    `json:"reason"`
	Message   string    `json:"message"`
	Source    string    `json:"source"`
	FirstTime time.Time `json:"first_time"`
	LastTime  time.Time `json:"last_time"`
	Count     int32     `json:"count"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status       string `json:"status"`
	K8sConnected bool   `json:"k8s_connected"`
	Error        string `json:"error,omitempty"`
}

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	Type      string      `json:"type"`
	Action    string      `json:"action"`
	Namespace string      `json:"namespace,omitempty"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// WatchEvent represents a Kubernetes watch event
type WatchEvent struct {
	Type   string      `json:"type"` // ADDED, MODIFIED, DELETED
	Object interface{} `json:"object"`
}
