<script lang="ts">
  import { apiClient } from '$lib/api/client';
  import type { KubernetesPod } from '$lib/types/kubernetes';
  export let pod: KubernetesPod;

  let showLogs = false;
  let logs = '';
  let logsLoading = false;
  let logsError: string | null = null;

  async function fetchLogs() {
    if (logs) return;
    logsLoading = true;
    logsError = null;
    try {
      const namespace = pod?.metadata?.namespace || 'default';
      const name = pod?.metadata?.name || '';
      const resp = await apiClient.getPodLogs(namespace, name);
      logs = resp.logs || '';
    } catch (err) {
      logsError = (err as Error)?.message || 'Failed to load logs';
    } finally {
      logsLoading = false;
    }
  }

  function toggleLogs() {
    showLogs = !showLogs;
    if (showLogs && !logs && !logsError) fetchLogs();
  }
</script>

<div class="logs-block">
  <h2>Logs</h2>
  <button on:click={toggleLogs} class="toggle-btn">{showLogs ? '▼' : '▶'} {showLogs ? 'Hide' : 'Show'} Logs</button>
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
</div>

<style>
  .logs { background:#0b1220; color:#e6eef8; padding:1rem; border-radius:6px; overflow:auto; max-height:400px; }
  .toggle-btn { background: #3b82f6; color: white; border: none; padding: 0.5rem 1rem; border-radius: 4px; cursor: pointer; font-weight: 600; margin-bottom: 0.5rem; }
  .toggle-btn:hover { background: #2563eb; }
  .error { color: red }
</style>
