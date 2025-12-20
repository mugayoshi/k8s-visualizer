package handlers

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/mugayoshi/k8s-visualizer/server/internal/models"
	"github.com/mugayoshi/k8s-visualizer/server/internal/services"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// mock that returns an error for GetClusterMetrics
type failingK8s struct{}

func (f *failingK8s) GetClientset() kubernetes.Interface { return nil }
func (f *failingK8s) IsHealthy() bool                    { return false }
func (f *failingK8s) GetClusterMetrics(ctx context.Context) (*services.ClusterMetrics, error) {
	return nil, errMock
}
func (f *failingK8s) GetNodeMetrics(ctx context.Context, nodeName string) (*services.NodeMetrics, error) {
	return nil, errMock
}
func (f *failingK8s) GetNamespaceMetrics(ctx context.Context, namespace string) (*services.NamespaceMetrics, error) {
	return nil, errMock
}
func (f *failingK8s) GetPodMetrics(ctx context.Context, namespace, podName string) (*services.PodMetrics, error) {
	return nil, errMock
}
func (f *failingK8s) GetPodLogs(ctx context.Context, namespace, podName string, opts *corev1.PodLogOptions) (io.ReadCloser, error) {
	return nil, errMock
}

// mock that returns success for GetClusterMetrics
type successK8s struct{}

func (s *successK8s) GetClientset() kubernetes.Interface { return nil }
func (s *successK8s) IsHealthy() bool                    { return true }
func (s *successK8s) GetClusterMetrics(ctx context.Context) (*services.ClusterMetrics, error) {
	return &services.ClusterMetrics{TotalNodes: 3}, nil
}
func (s *successK8s) GetNodeMetrics(ctx context.Context, nodeName string) (*services.NodeMetrics, error) {
	return &services.NodeMetrics{}, nil
}
func (s *successK8s) GetNamespaceMetrics(ctx context.Context, namespace string) (*services.NamespaceMetrics, error) {
	return &services.NamespaceMetrics{}, nil
}
func (s *successK8s) GetPodMetrics(ctx context.Context, namespace, podName string) (*services.PodMetrics, error) {
	return &services.PodMetrics{}, nil
}
func (s *successK8s) GetPodLogs(ctx context.Context, namespace, podName string, opts *corev1.PodLogOptions) (io.ReadCloser, error) {
	return nil, nil
}

var errMock = &mockError{}

type mockError struct{}

func (m *mockError) Error() string { return "mock error" }

func TestHandleClientMessage_GetMetrics_FailureDoesNotPanic(t *testing.T) {
	h := &WebSocketHandler{k8sClient: &failingK8s{}}
	send := make(chan models.WebSocketMessage, 1)

	h.handleClientMessage(models.WebSocketMessage{Action: "get_metrics"}, send, context.Background())

	select {
	case msg := <-send:
		t.Fatalf("expected no message on send channel, got: %v", msg)
	case <-time.After(10 * time.Millisecond):
		// success: no message
	}
}

func TestHandleClientMessage_GetMetrics_SuccessSendsMetrics(t *testing.T) {
	h := &WebSocketHandler{k8sClient: &successK8s{}}
	send := make(chan models.WebSocketMessage, 1)

	h.handleClientMessage(models.WebSocketMessage{Action: "get_metrics"}, send, context.Background())

	select {
	case msg := <-send:
		if msg.Type != "metrics" {
			t.Fatalf("expected message type 'metrics', got %s", msg.Type)
		}
		// Data should be a *services.ClusterMetrics
		if _, ok := msg.Data.(*services.ClusterMetrics); !ok {
			t.Fatalf("expected data to be *services.ClusterMetrics, got %T", msg.Data)
		}
	case <-time.After(50 * time.Millisecond):
		t.Fatalf("expected a metrics message but none received")
	}
}
