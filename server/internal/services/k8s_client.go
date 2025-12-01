// internal/services/k8s_client.go
package services

import (
	"context"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8sClient struct {
	clientset kubernetes.Interface
	config    *rest.Config
}

// K8sClientInterface defines the subset of methods used by handlers and tests.
type K8sClientInterface interface {
	GetClientset() kubernetes.Interface
	IsHealthy() bool
	GetClusterMetrics(ctx context.Context) (*ClusterMetrics, error)
	GetNodeMetrics(ctx context.Context, nodeName string) (*NodeMetrics, error)
	GetNamespaceMetrics(ctx context.Context, namespace string) (*NamespaceMetrics, error)
	GetPodMetrics(ctx context.Context, namespace, podName string) (*PodMetrics, error)
	GetPodLogs(ctx context.Context, namespace, podName string, opts *corev1.PodLogOptions) (io.ReadCloser, error)
}

func NewK8sClient() (*K8sClient, error) {
	var config *rest.Config
	var err error

	// Try in-cluster config first
	config, err = rest.InClusterConfig()
	if err != nil {
		// Fallback to kubeconfig
		var kubeconfig string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
		if envKubeconfig := os.Getenv("KUBECONFIG"); envKubeconfig != "" {
			kubeconfig = envKubeconfig
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}

		// If running in Docker, replace localhost with host.docker.internal
		if os.Getenv("RUNNING_IN_DOCKER") == "true" {
			u, err := url.Parse(config.Host)
			if err == nil {
				originalHost := u.Host
				if strings.HasPrefix(originalHost, "127.0.0.1") || strings.HasPrefix(originalHost, "localhost") {
					// The address to connect to
					u.Host = strings.Replace(originalHost, "127.0.0.1", "host.docker.internal", 1)
					u.Host = strings.Replace(u.Host, "localhost", "host.docker.internal", 1)
					config.Host = u.String()

					// The server name to use for TLS verification
					config.TLSClientConfig.ServerName = strings.Split(originalHost, ":")[0]
				}
			}
		}
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &K8sClient{
		clientset: clientset,
		config:    config,
	}, nil
}

func (k *K8sClient) GetClientset() kubernetes.Interface {
	return k.clientset
}

func (k *K8sClient) IsHealthy() bool {
	ctx := context.Background()
	_, err := k.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{Limit: 1})
	return err == nil
}

// GetPodLogs returns a stream (io.ReadCloser) for pod logs using provided PodLogOptions
func (k *K8sClient) GetPodLogs(ctx context.Context, namespace, podName string, opts *corev1.PodLogOptions) (io.ReadCloser, error) {
	req := k.clientset.CoreV1().Pods(namespace).GetLogs(podName, opts)
	return req.Stream(ctx)
}
