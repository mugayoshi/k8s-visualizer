// internal/handlers/pods.go
package handlers

import (

	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mugayoshi/k8s-visualizer/server/internal/services"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodHandler struct {
	k8sClient services.K8sClientInterface
}

func NewPodHandler(k8sClient services.K8sClientInterface) *PodHandler {
	return &PodHandler{k8sClient: k8sClient}
}

func (h *PodHandler) ListPods(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.DefaultQuery("namespace", "default")
	clientset := h.k8sClient.GetClientset()

	var pods *corev1.PodList
	var err error

	if namespace == "all" {
		pods, err = clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	} else {
		pods, err = clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]gin.H, 0, len(pods.Items))
	for _, pod := range pods.Items {
		containers := make([]gin.H, 0, len(pod.Spec.Containers))
		for _, container := range pod.Spec.Containers {
			ready := false
			for _, cs := range pod.Status.ContainerStatuses {
				if cs.Name == container.Name {
					ready = cs.Ready
					break
				}
			}
			containers = append(containers, gin.H{
				"name":  container.Name,
				"image": container.Image,
				"ready": ready,
			})
		}

		result = append(result, gin.H{
			"name":       pod.Name,
			"namespace":  pod.Namespace,
			"status":     string(pod.Status.Phase),
			"node":       pod.Spec.NodeName,
			"created":    pod.CreationTimestamp.Time,
			"labels":     pod.Labels,
			"containers": containers,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"pods":  result,
		"count": len(result),
	})
}

func (h *PodHandler) GetPod(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	name := c.Param("name")
	clientset := h.k8sClient.GetClientset()

	pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pod not found"})
		return
	}

	c.JSON(http.StatusOK, pod)
}

func (h *PodHandler) GetPodLogs(c *gin.Context) {
	ctx := c.Request.Context()
	namespace := c.Param("namespace")
	name := c.Param("name")
	container := c.DefaultQuery("container", "")
	tailLines := int64(100)

	opts := &corev1.PodLogOptions{
		TailLines: &tailLines,
	}
	if container != "" {
		opts.Container = container
	}

	logStream, err := h.k8sClient.GetPodLogs(ctx, namespace, name, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer logStream.Close()

	// Read all logs and return as string
	logBytes, err := io.ReadAll(logStream)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": string(logBytes),
	})
}
