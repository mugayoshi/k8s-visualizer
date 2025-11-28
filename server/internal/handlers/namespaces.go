// internal/handlers/namespaces.go
package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mugayoshi/k8s-visualizer/server/internal/services"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceHandler struct {
	k8sClient services.K8sClientInterface
}

func NewNamespaceHandler(k8sClient services.K8sClientInterface) *NamespaceHandler {
	return &NamespaceHandler{k8sClient: k8sClient}
}

func (h *NamespaceHandler) ListNamespaces(c *gin.Context) {
	ctx := context.Background()
	clientset := h.k8sClient.GetClientset()

	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]gin.H, 0, len(namespaces.Items))
	for _, ns := range namespaces.Items {
		result = append(result, gin.H{
			"name":    ns.Name,
			"status":  string(ns.Status.Phase),
			"created": ns.CreationTimestamp.Time,
			"labels":  ns.Labels,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"namespaces": result,
		"count":      len(result),
	})
}
