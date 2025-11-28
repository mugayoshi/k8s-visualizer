package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"

	"github.com/gin-gonic/gin"
	"github.com/mugayoshi/k8s-visualizer/server/internal/services"
)

// mock that implements K8sClientInterface for logs test
type logsMock struct{}

func (l *logsMock) GetClientset() kubernetes.Interface { return nil }
func (l *logsMock) IsHealthy() bool                    { return true }
func (l *logsMock) GetClusterMetrics(ctx context.Context) (*services.ClusterMetrics, error) {
	return nil, nil
}
func (l *logsMock) GetNodeMetrics(ctx context.Context, nodeName string) (*services.NodeMetrics, error) {
	return nil, nil
}
func (l *logsMock) GetNamespaceMetrics(ctx context.Context, namespace string) (*services.NamespaceMetrics, error) {
	return nil, nil
}
func (l *logsMock) GetPodMetrics(ctx context.Context, namespace, podName string) (*services.PodMetrics, error) {
	return nil, nil
}
func (l *logsMock) GetPodLogs(ctx context.Context, namespace, podName string, opts *corev1.PodLogOptions) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader("log-line-1\nlog-line-2\n")), nil
}

// localMockK8s is a small mock wrapping a fake clientset for use in these tests
type localMockK8s struct{ cs kubernetes.Interface }

func (m *localMockK8s) GetClientset() kubernetes.Interface { return m.cs }
func (m *localMockK8s) IsHealthy() bool                    { return true }
func (m *localMockK8s) GetClusterMetrics(ctx context.Context) (*services.ClusterMetrics, error) {
	return nil, nil
}
func (m *localMockK8s) GetNodeMetrics(ctx context.Context, nodeName string) (*services.NodeMetrics, error) {
	return nil, nil
}
func (m *localMockK8s) GetNamespaceMetrics(ctx context.Context, namespace string) (*services.NamespaceMetrics, error) {
	return nil, nil
}
func (m *localMockK8s) GetPodMetrics(ctx context.Context, namespace, podName string) (*services.PodMetrics, error) {
	return nil, nil
}
func (m *localMockK8s) GetPodLogs(ctx context.Context, namespace, podName string, opts *corev1.PodLogOptions) (io.ReadCloser, error) {
	return nil, nil
}

func TestListDeployments_ReturnsDeployments(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "deploy1", Namespace: "default"},
		Spec:       appsv1.DeploymentSpec{Replicas: int32Ptr(2)},
	}

	cs := fake.NewSimpleClientset(dep)
	mock := &localMockK8s{cs: cs}

	handler := NewDeploymentHandler(mock)
	r := gin.New()
	r.GET("/api/deployments", handler.ListDeployments)

	req := httptest.NewRequest(http.MethodGet, "/api/deployments?namespace=default", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Deployments []interface{} `json:"deployments"`
		Count       int           `json:"count"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}
	if resp.Count != 1 || len(resp.Deployments) != 1 {
		t.Fatalf("expected 1 deployment, got count=%d len(deployments)=%d", resp.Count, len(resp.Deployments))
	}
}

func TestListDeployments_FieldsMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "deploy2", Namespace: "ns2"},
		Spec:       appsv1.DeploymentSpec{Replicas: int32Ptr(3)},
		Status:     appsv1.DeploymentStatus{ReadyReplicas: 1, AvailableReplicas: 2, UpdatedReplicas: 3},
	}

	cs := fake.NewSimpleClientset(dep)
	mock := &localMockK8s{cs: cs}

	handler := NewDeploymentHandler(mock)
	r := gin.New()
	r.GET("/api/deployments", handler.ListDeployments)

	req := httptest.NewRequest(http.MethodGet, "/api/deployments?namespace=ns2", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Deployments []map[string]interface{} `json:"deployments"`
		Count       int                      `json:"count"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}
	if resp.Count != 1 || len(resp.Deployments) != 1 {
		t.Fatalf("expected 1 deployment, got count=%d len(deployments)=%d", resp.Count, len(resp.Deployments))
	}
	d := resp.Deployments[0]
	if d["name"] != "deploy2" || int(d["replicas"].(float64)) != 3 {
		t.Fatalf("unexpected deployment fields: %+v", d)
	}
	if int(d["ready_replicas"].(float64)) != 1 || int(d["available_replicas"].(float64)) != 2 || int(d["updated_replicas"].(float64)) != 3 {
		t.Fatalf("replica counts mismatch: %+v", d)
	}
}

func TestGetPod_NotFound_Returns404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cs := fake.NewSimpleClientset() // no pods
	mock := &localMockK8s{cs: cs}

	handler := NewPodHandler(mock)
	r := gin.New()
	r.GET("/api/pods/:namespace/:name", handler.GetPod)

	req := httptest.NewRequest(http.MethodGet, "/api/pods/does/notexist", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestListDeployments_ListError_Returns500(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cs := fake.NewSimpleClientset()
	// cause list deployments to error
	// cs is a *fake.Clientset — use its Fake to inject reactor
	fc := cs
	fc.Fake.PrependReactor("list", "deployments", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("list error")
	})

	mock := &localMockK8s{cs: cs}
	handler := NewDeploymentHandler(mock)
	r := gin.New()
	r.GET("/api/deployments", handler.ListDeployments)

	req := httptest.NewRequest(http.MethodGet, "/api/deployments?namespace=default", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetPod_ReturnsPod(t *testing.T) {
	gin.SetMode(gin.TestMode)

	pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "mypod", Namespace: "ns1"}}
	cs := fake.NewSimpleClientset(pod)
	mock := &localMockK8s{cs: cs}

	handler := NewPodHandler(mock)
	r := gin.New()
	r.GET("/api/pods/:namespace/:name", handler.GetPod)

	req := httptest.NewRequest(http.MethodGet, "/api/pods/ns1/mypod", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp corev1.Pod
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Name != "mypod" || resp.Namespace != "ns1" {
		t.Fatalf("unexpected pod returned: %+v", resp)
	}
}

func TestGetPodLogs_ReturnsLogs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	lm := &logsMock{}
	handler := NewPodHandler(lm)
	r := gin.New()
	r.GET("/api/pods/:namespace/:name/logs", handler.GetPodLogs)

	req := httptest.NewRequest(http.MethodGet, "/api/pods/ns/log1/logs", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Logs string `json:"logs"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal logs response: %v", err)
	}
	if !strings.Contains(resp.Logs, "log-line-1") {
		t.Fatalf("unexpected logs: %s", resp.Logs)
	}
}

func int32Ptr(i int32) *int32 { return &i }
