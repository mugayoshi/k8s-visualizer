// src/lib/api/client.ts
import axios, { type AxiosInstance } from 'axios';
import type { Node, Pod, Namespace, Deployment, ClusterMetrics } from '$lib/types/kubernetes';

class ApiClient {
  private client: AxiosInstance;

  constructor(baseURL: string = 'http://localhost:8080') {
    this.client = axios.create({
      baseURL,
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json',
      },
    });
  }

  // Health check
  async healthCheck() {
    const response = await this.client.get('/health');
    return response.data;
  }

  // Metrics
  async getClusterMetrics(): Promise<ClusterMetrics> {
    const response = await this.client.get('/api/metrics');
    return response.data;
  }

  // Nodes
  async listNodes(): Promise<{ nodes: Node[]; count: number }> {
    const response = await this.client.get('/api/nodes');
    return response.data;
  }

  async getNode(name: string): Promise<Node> {
    const response = await this.client.get(`/api/nodes/${name}`);
    return response.data;
  }

  // Pods
  async listPods(namespace: string = 'default'): Promise<{ pods: Pod[]; count: number }> {
    const response = await this.client.get('/api/pods', {
      params: { namespace },
    });
    return response.data;
  }

  async getPod(namespace: string, name: string): Promise<Pod> {
    const response = await this.client.get(`/api/pods/${namespace}/${name}`);
    return response.data;
  }

  async getPodLogs(namespace: string, name: string, container?: string): Promise<{ logs: string }> {
    const response = await this.client.get(`/api/pods/${namespace}/${name}/logs`, {
      params: { container },
    });
    return response.data;
  }

  // Namespaces
  async listNamespaces(): Promise<{ namespaces: Namespace[]; count: number }> {
    const response = await this.client.get('/api/namespaces');
    return response.data;
  }

  // Deployments
  async listDeployments(namespace: string = 'default'): Promise<{ deployments: Deployment[]; count: number }> {
    const response = await this.client.get('/api/deployments', {
      params: { namespace },
    });
    return response.data;
  }
}

// Export a singleton instance
export const apiClient = new ApiClient(
  import.meta.env.VITE_API_URL || 'http://localhost:8080'
);

// Reusable fetchPods function for use in Svelte pages
export async function fetchPods(namespace: string = 'all'): Promise<Pod[]> {
  const response = await apiClient.listPods(namespace);
  return response.pods;
}