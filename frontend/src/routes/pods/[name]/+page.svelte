<script lang="ts">
  import { onMount } from 'svelte';
  import type { KubernetesPod } from '$lib/types/kubernetes';
  import { apiClient } from '$lib/api/client';
  export let data: { pod: KubernetesPod | null; error?: string };

  const refreshInterval = 10 * 1000; // 10 seconds

  let pod = data.pod;
  let error: string | null = data.error || null;
  let lastFetched: string | null = pod ? new Date().toISOString() : null;
  let showRawJSON = false;
  let showLogs = false;
  let logs = '';
  let logsLoading = false;
  let logsError: string | null = null;
  let refreshIntervalId: number | null = null;
  let isRefreshing = false;

  async function refreshPodData() {
    if (!pod) return;
    isRefreshing = true;
    try {
      const namespace = pod.metadata?.namespace || 'default';
      const name = pod.metadata?.name;
      pod = (await apiClient.getPod(namespace, name)) as unknown as KubernetesPod;
      error = null;
      lastFetched = new Date().toISOString();
    } catch (err) {
      error = (err as Error)?.message || 'Failed to refresh pod';
    } finally {
      isRefreshing = false;
    }
  }

  onMount(() => {
    refreshIntervalId = window.setInterval(refreshPodData, refreshInterval);
    return () => {
      if (refreshIntervalId) clearInterval(refreshIntervalId);
    };
  });

  async function fetchLogs() {
    if (logs) return; // Already fetched
    logsLoading = true;
    logsError = null;
    try {
      const namespace = pod?.metadata?.namespace || 'default';
      const name = pod?.metadata?.name || '';
      const response = await apiClient.getPodLogs(namespace, name);
      logs = response.logs || '';
    } catch (err) {
      logsError = (err as Error)?.message || 'Failed to load logs';
    } finally {
      logsLoading = false;
    }
  }

  function toggleLogs() {
    showLogs = !showLogs;
    if (showLogs && !logs && !logsError) {
      fetchLogs();
    }
  }

  function statusClass(status: string | undefined): string {
    const s = (status || '').toLowerCase();
    if (s.includes('running')) return 'status-running';
    if (s.includes('pending')) return 'status-pending';
    if (s.includes('failed')) return 'status-failed';
    if (s.includes('succeeded') || s.includes('completed') || s.includes('success')) return 'status-succeeded';
    return 'status-unknown';
  }
</script>

<div class="container">
  {#if error}
    <p class="error">Error: {error}</p>
  {:else if !pod}
    <p>Pod not found or still loading.</p>
  {:else}
    <div class="header">
      <a href="/pods" class="back-btn">← Back to Pods</a>
      <h1>Pod: {pod.metadata?.name}</h1>
      <div>
        <div class="refresh-status" class:refreshing={isRefreshing}>
          {isRefreshing ? '⟳ Refreshing...' : `✓ Auto-refresh every ${refreshInterval / 1000}s`}
        </div>
        <div class="last-fetched">Last fetched: {lastFetched ? new Date(lastFetched).toLocaleString() : '-'}</div>
      </div>
    </div>

    <div class="meta">
      <div><strong>Namespace:</strong> {pod.metadata?.namespace}</div>
      <div><strong>Status:</strong> <span class={statusClass(pod.status?.phase)}>{pod.status?.phase || 'Unknown'}</span></div>
      <div><strong>Node:</strong> {pod.spec?.nodeName || '-'}</div>
      <div><strong>Created:</strong> {pod.metadata?.creationTimestamp}</div>
    </div>

    <h2>Labels</h2>
    {#if Object.keys(pod.metadata?.labels || {}).length === 0}
      <p>No labels</p>
    {:else}
      <ul class="labels">
        {#each Object.entries(pod.metadata?.labels || {}) as [k, v]}
          <li><code>{k}</code>: {v}</li>
        {/each}
      </ul>
    {/if}

    <h2>Containers ({pod.spec?.containers?.length ?? 0})</h2>
    <table class="containers">
      <thead>
        <tr><th>Name</th><th>Image</th><th>Ready</th></tr>
      </thead>
      <tbody>
        {#each pod.spec?.containers || [] as c}
          <tr>
            <td>{c.name}</td>
            <td>{c.image}</td>
            <td>{pod.status?.containerStatuses?.find((cs) => cs.name === c.name)?.ready ? 'yes' : 'no'}</td>
          </tr>
        {/each}
      </tbody>
    </table>

    <h2>Raw JSON</h2>
    <button on:click={() => (showRawJSON = !showRawJSON)} class="toggle-btn">
      {showRawJSON ? '▼' : '▶'} {showRawJSON ? 'Hide' : 'Show'} Raw JSON
    </button>
    {#if showRawJSON}
      <pre class="json">{JSON.stringify(pod, null, 2)}</pre>
    {/if}

    <h2>Logs</h2>
    <button on:click={toggleLogs} class="toggle-btn">
      {showLogs ? '▼' : '▶'} {showLogs ? 'Hide' : 'Show'} Logs
    </button>
    {#if showLogs}
      {#if logsLoading}
        <p>Loading logs...</p>
      {:else if logsError}
        <p class="error">Error: {logsError}</p>
      {:else if logs}
        <pre class="logs">{logs}</pre>
      {:else}
        <p>No logs available</p>
      {/if}
    {/if}
  {/if}
</div>

<style>
  .container { padding: 1rem; }
  .header { display: flex; align-items: center; gap: 1rem; margin-bottom: 1rem; }
  .back-btn { display: inline-block; padding: 0.5rem 1rem; background: #e5e7eb; color: #1f2937; text-decoration: none; border-radius: 4px; font-weight: 600; transition: background 0.2s; }
  .back-btn:hover { background: #d1d5db; }
  .refresh-status { font-size: 0.875rem; color: #6b7280; padding: 0.5rem 1rem; background: #f3f4f6; border-radius: 4px; }
  .refresh-status.refreshing { color: #0066cc; font-weight: 600; animation: pulse 1s ease-in-out infinite; }
  @keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.6; } }
  .meta { display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 0.5rem; margin-bottom: 1rem; }
  .labels { list-style: none; padding: 0; display: flex; gap: 0.5rem; flex-wrap: wrap; }
  .labels li { background:#f3f4f6; padding:0.25rem 0.5rem; border-radius:4px; }
  .containers { width:100%; border-collapse: collapse; margin-bottom:1rem; }
  .containers th, .containers td { border:1px solid #e5e7eb; padding:0.5rem; text-align:left; }
  .json { background:#0b1220; color:#e6eef8; padding:1rem; border-radius:6px; overflow:auto; }
  .logs { background:#0b1220; color:#e6eef8; padding:1rem; border-radius:6px; overflow:auto; max-height:400px; }
  .toggle-btn { background: #3b82f6; color: white; border: none; padding: 0.5rem 1rem; border-radius: 4px; cursor: pointer; font-weight: 600; margin-bottom: 0.5rem; }
  .toggle-btn:hover { background: #2563eb; }
  .status-running { color: green; font-weight:600 }
  .status-pending { color: orange; font-weight:600 }
  .status-failed { color: red; font-weight:600 }
  .status-succeeded { color: teal; font-weight:600 }
  .status-unknown { color: gray; font-weight:600 }
  .error { color: red }
</style>
