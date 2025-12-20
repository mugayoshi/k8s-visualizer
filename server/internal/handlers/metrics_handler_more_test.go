package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mugayoshi/k8s-visualizer/server/internal/services"
)

// mock success client
type successMetricsClient struct{}

func (m *successMetricsClient) GetClusterMetrics(ctx context.Context) (*services.ClusterMetrics, error) {
	return &services.ClusterMetrics{
		TotalNodes:      2,
		TotalPods:       5,
		TotalNamespaces: 1,
		CPUCapacity:     "2",
		MemoryCapacity:  "4Gi",
	}, nil
}

func TestMetricsEndpoint_Success_Returns200AndBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	mock := &successMetricsClient{}

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

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var got services.ClusterMetrics
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if got.TotalNodes != 2 || got.TotalPods != 5 || got.TotalNamespaces != 1 {
		t.Fatalf("unexpected metrics payload: %+v", got)
	}
}
