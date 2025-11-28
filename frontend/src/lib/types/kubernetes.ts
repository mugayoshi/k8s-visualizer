// src/lib/types/kubernetes.ts
export interface Node {
  name: string;
  status: string;
  created: string;
  cpu_capacity: string;
  memory_capacity: string;
  pod_capacity: string;
  cpu_allocatable?: string;
  memory_allocatable?: string;
  labels: Record<string, string>;
  addresses: Array<{
    type: string;
    address: string;
  }>;
}

export interface Container {
  name: string;
  image: string;
  ready: boolean;
  restart_count?: number;
  state?: string;
}

export interface Pod {
  name: string;
  namespace: string;
  status: string;
  status_detail?: string;
  ready_containers?: string;
  phase?: string;
  node: string;
  created: string;
  labels: Record<string, string>;
  containers: Container[];
  restart_count?: number;
  cpu_request?: string;
  cpu_limit?: string;
  memory_request?: string;
  memory_limit?: string;
  container_count?: number;
  age?: string;
}

export interface Namespace {
  name: string;
  status: string;
  created: string;
  labels: Record<string, string>;
}

export interface Deployment {
  name: string;
  namespace: string;
  replicas: number;
  ready_replicas: number;
  available_replicas: number;
  updated_replicas?: number;
  created: string;
  labels: Record<string, string>;
}

export interface ClusterMetrics {
  total_nodes: number;
  total_pods: number;
  total_namespaces: number;
  cpu_capacity: string;
  memory_capacity: string;
  node_metrics: Array<{
    name: string;
    cpu_capacity: string;
    memory_capacity: string;
    pod_count: number;
    status: string;
  }>;
}

export interface WebSocketMessage {
  type: string;
  action: string;
  namespace?: string;
  data: any;
  timestamp: string;
}