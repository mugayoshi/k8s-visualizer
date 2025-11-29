<script lang="ts">
  // Logs are provided by the parent component/page via props.
  // This component no longer performs network requests.
  export let logs: string = '';
  export let loading: boolean = false;
  export let error: string | null = null;

  let showLogs = false;

  function toggleLogs() {
    showLogs = !showLogs;
  }
</script>

<div class="logs-block">
  <h2>Logs</h2>
  <button on:click={toggleLogs} class="toggle-btn">{showLogs ? '▼' : '▶'} {showLogs ? 'Hide' : 'Show'} Logs</button>
  {#if showLogs}
    {#if loading}
      <p>Loading logs...</p>
    {:else if error}
      <p class="error">Error: {error}</p>
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
