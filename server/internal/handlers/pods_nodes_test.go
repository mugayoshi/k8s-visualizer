package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "io"

    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/resource"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/kubernetes/fake"

    "github.com/gin-gonic/gin"
    "github.com/mugayoshi/k8s-visualizer/server/internal/services"
)

type mockK8s struct{
    cs kubernetes.Interface
}

func (m *mockK8s) GetClientset() kubernetes.Interface { return m.cs }
func (m *mockK8s) IsHealthy() bool { return true }
func (m *mockK8s) GetClusterMetrics(ctx context.Context) (*services.ClusterMetrics, error) { return nil, nil }
func (m *mockK8s) GetNodeMetrics(ctx context.Context, nodeName string) (*services.NodeMetrics, error) { return nil, nil }
func (m *mockK8s) GetNamespaceMetrics(ctx context.Context, namespace string) (*services.NamespaceMetrics, error) { return nil, nil }
func (m *mockK8s) GetPodMetrics(ctx context.Context, namespace, podName string) (*services.PodMetrics, error) { return nil, nil }
func (m *mockK8s) GetPodLogs(ctx context.Context, namespace, podName string, opts *corev1.PodLogOptions) (io.ReadCloser, error) { return nil, nil }

func TestListPods_ReturnsPods(t *testing.T) {
    gin.SetMode(gin.TestMode)

    pod := &corev1.Pod{
        ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "default"},
        Spec: corev1.PodSpec{Containers: []corev1.Container{{Name: "c1", Image: "img"}}},
        Status: corev1.PodStatus{Phase: corev1.PodRunning},
    }

    cs := fake.NewSimpleClientset(pod)
    mock := &mockK8s{cs: cs}

    handler := NewPodHandler(mock)
    r := gin.New()
    r.GET("/api/pods", handler.ListPods)

    req := httptest.NewRequest(http.MethodGet, "/api/pods?namespace=default", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
    }

    var resp struct{
        Pods []interface{} `json:"pods"`
        Count int `json:"count"`
    }
    if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
        t.Fatalf("failed to unmarshal body: %v", err)
    }
    if resp.Count != 1 || len(resp.Pods) != 1 {
        t.Fatalf("expected 1 pod, got count=%d len(pods)=%d", resp.Count, len(resp.Pods))
    }
}

func TestListNodes_ReturnsNodes(t *testing.T) {
    gin.SetMode(gin.TestMode)

    node := &corev1.Node{
        ObjectMeta: metav1.ObjectMeta{Name: "node1"},
        Status: corev1.NodeStatus{
            Capacity: corev1.ResourceList{
                corev1.ResourceCPU: resource.MustParse("2"),
                corev1.ResourceMemory: resource.MustParse("1Gi"),
                corev1.ResourcePods: resource.MustParse("110"),
            },
            Conditions: []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionTrue}},
        },
    }

    cs := fake.NewSimpleClientset(node)
    mock := &mockK8s{cs: cs}

    handler := NewNodeHandler(mock)
    r := gin.New()
    r.GET("/api/nodes", handler.ListNodes)

    req := httptest.NewRequest(http.MethodGet, "/api/nodes", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
    }

    var resp struct{
        Nodes []interface{} `json:"nodes"`
        Count int `json:"count"`
    }
    if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
        t.Fatalf("failed to unmarshal body: %v", err)
    }
    if resp.Count != 1 || len(resp.Nodes) != 1 {
        t.Fatalf("expected 1 node, got count=%d len(nodes)=%d", resp.Count, len(resp.Nodes))
    }
}
