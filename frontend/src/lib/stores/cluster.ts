// src/lib/stores/cluster.ts
import { writable, derived } from 'svelte/store';
import type { Node, Pod, Namespace, ClusterMetrics } from '$lib/types/kubernetes';

// Cluster data stores
export const nodes = writable<Node[]>([]);
export const pods = writable<Pod[]>([]);
export const namespaces = writable<Namespace[]>([]);
export const metrics = writable<ClusterMetrics | null>(null);

// Loading states
export const loading = writable({
  nodes: false,
  pods: false,
  namespaces: false,
  metrics: false,
});

// Error states
export const errors = writable({
  nodes: null as string | null,
  pods: null as string | null,
  namespaces: null as string | null,
  metrics: null as string | null,
});

// Selected namespace filter
export const selectedNamespace = writable<string>('all');

// Derived stores
export const filteredPods = derived(
  [pods, selectedNamespace],
  ([$pods, $selectedNamespace]) => {
    if ($selectedNamespace === 'all') {
      return $pods;
    }
    return $pods.filter(pod => pod.namespace === $selectedNamespace);
  }
);

export const podsByStatus = derived(pods, ($pods) => {
  return $pods.reduce((acc, pod) => {
    const status = pod.status.toLowerCase();
    acc[status] = (acc[status] || 0) + 1;
    return acc;
  }, {} as Record<string, number>);
});