// internal/handlers/deployments.go
package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mugayoshi/k8s-visualizer/server/internal/services"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentHandler struct {
	k8sClient services.K8sClientInterface
}

func NewDeploymentHandler(k8sClient services.K8sClientInterface) *DeploymentHandler {
	return &DeploymentHandler{k8sClient: k8sClient}
}

func (h *DeploymentHandler) ListDeployments(c *gin.Context) {
	ctx := context.Background()
	namespace := c.DefaultQuery("namespace", "default")
	clientset := h.k8sClient.GetClientset()

	deployments, err := clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]gin.H, 0, len(deployments.Items))
	for _, deploy := range deployments.Items {
		readyReplicas := int32(0)
		availableReplicas := int32(0)

		if deploy.Status.ReadyReplicas != 0 {
			readyReplicas = deploy.Status.ReadyReplicas
		}
		if deploy.Status.AvailableReplicas != 0 {
			availableReplicas = deploy.Status.AvailableReplicas
		}

		result = append(result, gin.H{
			"name":               deploy.Name,
			"namespace":          deploy.Namespace,
			"replicas":           *deploy.Spec.Replicas,
			"ready_replicas":     readyReplicas,
			"available_replicas": availableReplicas,
			"created":            deploy.CreationTimestamp.Time,
			"labels":             deploy.Labels,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"deployments": result,
		"count":       len(result),
	})
}
