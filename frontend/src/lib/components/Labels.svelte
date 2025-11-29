<script lang="ts">
  export let labels: Record<string, string> | undefined = {};
</script>

{#if !labels || Object.keys(labels).length === 0}
  <p>No labels</p>
{:else}
  <ul class="labels">
    {#each Object.entries(labels) as [k, v]}
      <li><code>{k}</code>: {v}</li>
    {/each}
  </ul>
{/if}

<style>
  .labels { list-style: none; padding: 0; display: flex; gap: 0.5rem; flex-wrap: wrap; }
  .labels li { background:#f3f4f6; padding:0.25rem 0.5rem; border-radius:4px; }
</style>
