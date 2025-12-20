package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mugayoshi/k8s-visualizer/server/internal/services"
)

// mock that simulates a failing K8s client for metrics
type failingMetricsClient struct{}

func (m *failingMetricsClient) GetClusterMetrics(ctx context.Context) (*services.ClusterMetrics, error) {
	return nil, errors.New("simulated k8s failure")
}

// We only implement the single method we need for this test.

func TestMetricsEndpoint_K8sClientFails_Returns500(t *testing.T) {
	// Use gin in test mode
	gin.SetMode(gin.TestMode)
	r := gin.New()

	mock := &failingMetricsClient{}

	// Register the same handler shape as in main.go but using our mock
	r.GET("/api/metrics", func(c *gin.Context) {
		ctx := c.Request.Context()
		metrics, err := mock.GetClusterMetrics(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, metrics)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/metrics", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d, body: %s", w.Code, w.Body.String())
	}
}
