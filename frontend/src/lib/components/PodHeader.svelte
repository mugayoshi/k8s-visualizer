<script lang="ts">
  import type { KubernetesPod } from '$lib/types/kubernetes';
  export let pod: KubernetesPod;
  export let lastFetched: string | null;
  export let isRefreshing: boolean;
  export let refreshIntervalSeconds: number = 10;
</script>

<div class="header">
  <a href="/pods" class="back-btn">← Back to Pods</a>
  <h1>Pod: {pod.metadata?.name}</h1>
  <div>
    <div class="refresh-status" class:refreshing={isRefreshing}>
      {isRefreshing ? '⟳ Refreshing...' : `✓ Auto-refresh every ${refreshIntervalSeconds}s`}
    </div>
    <div class="last-fetched">Last fetched: {lastFetched ? new Date(lastFetched).toLocaleString() : '-'}</div>
  </div>
</div>

<style>
  .header { display: flex; align-items: center; gap: 1rem; margin-bottom: 1rem; }
  .back-btn { display: inline-block; padding: 0.5rem 1rem; background: #e5e7eb; color: #1f2937; text-decoration: none; border-radius: 4px; font-weight: 600; transition: background 0.2s; }
  .back-btn:hover { background: #d1d5db; }
  .refresh-status { font-size: 0.875rem; color: #6b7280; padding: 0.5rem 1rem; background: #f3f4f6; border-radius: 4px; }
  .refresh-status.refreshing { color: #0066cc; font-weight: 600; animation: pulse 1s ease-in-out infinite; }
  @keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.6; } }
  .last-fetched { font-size: 0.85rem; color: #6b7280; margin-top: 4px; }
</style>
