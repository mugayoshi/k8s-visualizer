// internal/handlers/nodes.go
package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mugayoshi/k8s-visualizer/server/internal/services"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodeHandler struct {
	k8sClient *services.K8sClient
}

func NewNodeHandler(k8sClient *services.K8sClient) *NodeHandler {
	return &NodeHandler{k8sClient: k8sClient}
}

func (h *NodeHandler) ListNodes(c *gin.Context) {
	ctx := context.Background()
	clientset := h.k8sClient.GetClientset()

	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Transform to simpler response
	result := make([]gin.H, 0, len(nodes.Items))
	for _, node := range nodes.Items {
		nodeInfo := gin.H{
			"name":            node.Name,
			"created":         node.CreationTimestamp.Time,
			"cpu_capacity":    node.Status.Capacity.Cpu().String(),
			"memory_capacity": node.Status.Capacity.Memory().String(),
			"pod_capacity":    node.Status.Capacity.Pods().String(),
			"labels":          node.Labels,
			"addresses":       node.Status.Addresses,
		}

		// Get status
		if len(node.Status.Conditions) > 0 {
			for _, condition := range node.Status.Conditions {
				if condition.Type == "Ready" {
					nodeInfo["status"] = string(condition.Status)
					break
				}
			}
		}

		result = append(result, nodeInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes": result,
		"count": len(result),
	})
}

func (h *NodeHandler) GetNode(c *gin.Context) {
	ctx := context.Background()
	nodeName := c.Param("name")
	clientset := h.k8sClient.GetClientset()

	node, err := clientset.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Node not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":            node.Name,
		"created":         node.CreationTimestamp.Time,
		"cpu_capacity":    node.Status.Capacity.Cpu().String(),
		"memory_capacity": node.Status.Capacity.Memory().String(),
		"conditions":      node.Status.Conditions,
		"labels":          node.Labels,
	})
}
