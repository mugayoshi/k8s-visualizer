// cmd/server/main.go
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mugayoshi/k8s-visualizer/server/internal/handlers"
	"github.com/mugayoshi/k8s-visualizer/server/internal/middleware"
	"github.com/mugayoshi/k8s-visualizer/server/internal/services"
)

func main() {
	// Initialize Kubernetes client
	k8sClient, err := services.NewK8sClient()
	if err != nil {
		log.Fatalf("Failed to create K8s client: %v", err)
	}

	// Create Gin router
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())
	// TODO
	// r.Use(middleware.Logger())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":        "healthy",
			"k8s_connected": k8sClient.IsHealthy(),
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Node endpoints
		nodeHandler := handlers.NewNodeHandler(k8sClient)
		api.GET("/nodes", nodeHandler.ListNodes)
		api.GET("/nodes/:name", nodeHandler.GetNode)

		// Pod endpoints
		podHandler := handlers.NewPodHandler(k8sClient)
		api.GET("/pods", podHandler.ListPods)
		api.GET("/pods/:namespace/:name", podHandler.GetPod)
		api.GET("/pods/:namespace/:name/logs", podHandler.GetPodLogs)

		// Namespace endpoints
		namespaceHandler := handlers.NewNamespaceHandler(k8sClient)
		api.GET("/namespaces", namespaceHandler.ListNamespaces)

		// Deployment endpoints
		deploymentHandler := handlers.NewDeploymentHandler(k8sClient)
		api.GET("/deployments", deploymentHandler.ListDeployments)

		// TODO
		// WebSocket endpoint
		// wsHandler := handlers.NewWebSocketHandler(k8sClient)
		// api.GET("/ws", wsHandler.HandleWebSocket)
	}

	// Start server
	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
