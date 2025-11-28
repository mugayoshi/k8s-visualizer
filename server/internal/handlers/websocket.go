// internal/handlers/websocket.go
package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mugayoshi/k8s-visualizer/server/internal/models"
	"github.com/mugayoshi/k8s-visualizer/server/internal/services"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins in development
		// In production, you should check the origin
		return true
	},
}

type WebSocketHandler struct {
	k8sClient services.K8sClientInterface
}

func NewWebSocketHandler(k8sClient services.K8sClientInterface) *WebSocketHandler {
	return &WebSocketHandler{k8sClient: k8sClient}
}

// HandleWebSocket handles WebSocket connections for real-time updates
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("WebSocket client connected from %s", conn.RemoteAddr())

	// Create context with cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel for sending messages to client
	send := make(chan models.WebSocketMessage, 256)

	// Start goroutine to write messages to WebSocket
	go h.writeMessages(conn, send, cancel)

	// Start goroutine to read messages from WebSocket
	go h.readMessages(conn, send, cancel)

	// Start watching Kubernetes resources
	go h.watchPods(ctx, send, "")

	// Keep connection alive until context is cancelled
	<-ctx.Done()
	log.Printf("WebSocket client disconnected from %s", conn.RemoteAddr())
}

// writeMessages writes messages from the send channel to the WebSocket
func (h *WebSocketHandler) writeMessages(conn *websocket.Conn, send chan models.WebSocketMessage, cancel context.CancelFunc) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-send:
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteJSON(message); err != nil {
				log.Printf("Error writing to WebSocket: %v", err)
				cancel()
				return
			}

		case <-ticker.C:
			// Send ping to keep connection alive
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error sending ping: %v", err)
				cancel()
				return
			}
		}
	}
}

// readMessages reads messages from the WebSocket
func (h *WebSocketHandler) readMessages(conn *websocket.Conn, send chan models.WebSocketMessage, cancel context.CancelFunc) {
	defer cancel()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var message models.WebSocketMessage
		err := conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		log.Printf("Received message: type=%s, action=%s, namespace=%s", message.Type, message.Action, message.Namespace)

		// Handle client messages
		h.handleClientMessage(message, send)
	}
}

// handleClientMessage handles messages received from the client
func (h *WebSocketHandler) handleClientMessage(message models.WebSocketMessage, send chan models.WebSocketMessage) {
	switch message.Action {
	case "subscribe_pods":
		// Client wants to subscribe to pod updates
		namespace := message.Namespace
		log.Printf("Client subscribed to pods in namespace: %s", namespace)
		// You could start a new watcher here based on client preferences

	case "subscribe_nodes":
		// Client wants to subscribe to node updates
		log.Printf("Client subscribed to node updates")

	case "get_metrics":
		// Client requests current metrics
		ctx := context.Background()
		metrics, err := h.k8sClient.GetClusterMetrics(ctx)
		if err != nil {
			log.Printf("Error getting metrics: %v", err)
			return
		}

		send <- models.WebSocketMessage{
			Type:      "metrics",
			Action:    "update",
			Data:      metrics,
			Timestamp: time.Now(),
		}

	default:
		log.Printf("Unknown action: %s", message.Action)
	}
}

// watchPods watches for pod changes and sends updates via WebSocket
func (h *WebSocketHandler) watchPods(ctx context.Context, send chan models.WebSocketMessage, namespace string) {
	clientset := h.k8sClient.GetClientset()

	// Create watcher
	watcher, err := clientset.CoreV1().Pods(namespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		log.Printf("Error creating pod watcher: %v", err)
		return
	}
	defer watcher.Stop()

	log.Printf("Started watching pods in namespace: %s", namespace)

	// Send initial pod list
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err == nil {
		result := make([]map[string]interface{}, 0, len(pods.Items))
		for _, pod := range pods.Items {
			// Get metrics for the pod
			metrics, err := h.k8sClient.GetPodMetrics(ctx, pod.Namespace, pod.Name)
			podData := map[string]interface{}{
				"name":      pod.Name,
				"namespace": pod.Namespace,
				"status":    string(pod.Status.Phase),
				"node":      pod.Spec.NodeName,
			}

			if err == nil {
				if metrics.CPURequest != "" {
					podData["cpu_request"] = metrics.CPURequest
				}
				if metrics.CPULimit != "" {
					podData["cpu_limit"] = metrics.CPULimit
				}
				if metrics.MemoryRequest != "" {
					podData["memory_request"] = metrics.MemoryRequest
				}
				if metrics.MemoryLimit != "" {
					podData["memory_limit"] = metrics.MemoryLimit
				}
				if metrics.Age != "" {
					podData["age"] = metrics.Age
				}
				podData["container_count"] = metrics.ContainerCount
				podData["restart_count"] = metrics.RestartCount
			}

			result = append(result, podData)
		}

		send <- models.WebSocketMessage{
			Type:      "pods",
			Action:    "initial",
			Namespace: namespace,
			Data:      result,
			Timestamp: time.Now(),
		}
	}

	// Watch for changes
	for {
		select {
		case <-ctx.Done():
			log.Printf("Stopped watching pods in namespace: %s", namespace)
			return

		case event, ok := <-watcher.ResultChan():
			if !ok {
				log.Printf("Pod watcher channel closed, restarting...")
				// Restart watcher
				time.Sleep(1 * time.Second)
				h.watchPods(ctx, send, namespace)
				return
			}

			// Convert event to JSON
			data, err := json.Marshal(event.Object)
			if err != nil {
				log.Printf("Error marshaling event: %v", err)
				continue
			}

			var podData map[string]interface{}
			if err := json.Unmarshal(data, &podData); err != nil {
				log.Printf("Error unmarshaling pod data: %v", err)
				continue
			}

			// Extract relevant fields
			metadata := podData["metadata"].(map[string]interface{})
			status := podData["status"].(map[string]interface{})
			spec := podData["spec"].(map[string]interface{})

			// Get metrics for the pod
			podNamespace := metadata["namespace"].(string)
			podName := metadata["name"].(string)
			metrics, err := h.k8sClient.GetPodMetrics(ctx, podNamespace, podName)

			simplifiedPod := map[string]interface{}{
				"name":      metadata["name"],
				"namespace": metadata["namespace"],
				"status":    status["phase"],
				"node":      spec["nodeName"],
			}

			if err == nil {
				if metrics.CPURequest != "" {
					simplifiedPod["cpu_request"] = metrics.CPURequest
				}
				if metrics.CPULimit != "" {
					simplifiedPod["cpu_limit"] = metrics.CPULimit
				}
				if metrics.MemoryRequest != "" {
					simplifiedPod["memory_request"] = metrics.MemoryRequest
				}
				if metrics.MemoryLimit != "" {
					simplifiedPod["memory_limit"] = metrics.MemoryLimit
				}
				if metrics.Age != "" {
					simplifiedPod["age"] = metrics.Age
				}
				simplifiedPod["container_count"] = metrics.ContainerCount
				simplifiedPod["restart_count"] = metrics.RestartCount
			}

			// Determine event type
			var eventType string
			switch event.Type {
			case watch.Added:
				eventType = "added"
			case watch.Modified:
				eventType = "modified"
			case watch.Deleted:
				eventType = "deleted"
			default:
				eventType = "unknown"
			}

			// Send update to client
			send <- models.WebSocketMessage{
				Type:      "pods",
				Action:    eventType,
				Namespace: namespace,
				Data:      simplifiedPod,
				Timestamp: time.Now(),
			}

			log.Printf("Pod %s: %v in namespace %s", eventType, metadata["name"], metadata["namespace"])
		}
	}
}
