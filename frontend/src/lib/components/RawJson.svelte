<script lang="ts">
  import type { KubernetesPod } from '$lib/types/kubernetes';
  export let pod: KubernetesPod;

  let showRawJSON = false;
  function toggleRaw() {
    showRawJSON = !showRawJSON;
  }
</script>

<div class="raw-json">
  <h2>Raw JSON</h2>
  <button on:click={toggleRaw} class="toggle-btn">{showRawJSON ? '▼' : '▶'} {showRawJSON ? 'Hide' : 'Show'} Raw JSON</button>
  {#if showRawJSON}
    <pre class="json">{JSON.stringify(pod, null, 2)}</pre>
  {/if}
</div>

<style>
  .json { background:#0b1220; color:#e6eef8; padding:1rem; border-radius:6px; overflow:auto; }
  .toggle-btn { background: #3b82f6; color: white; border: none; padding: 0.5rem 1rem; border-radius: 4px; cursor: pointer; font-weight: 600; margin-bottom: 0.5rem; }
  .toggle-btn:hover { background: #2563eb; }
</style>
