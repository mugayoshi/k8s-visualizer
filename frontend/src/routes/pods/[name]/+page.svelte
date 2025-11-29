<script lang="ts">
  import { onMount } from 'svelte';
  import type { KubernetesPod } from '$lib/types/kubernetes';
  import { apiClient } from '$lib/api/client';
  import PodHeader from '$lib/components/PodHeader.svelte';
  import Metadata from '$lib/components/Metadata.svelte';
  import Labels from '$lib/components/Labels.svelte';
  import ContainersTable from '$lib/components/ContainersTable.svelte';
  import RawJson from '$lib/components/RawJson.svelte';
  import Logs from '$lib/components/Logs.svelte';
  export let data: { pod: KubernetesPod | null; error?: string };

  const refreshInterval = 10 * 1000; // 10 seconds

  let pod = data.pod;
  let error: string | null = data.error || null;
  let lastFetched: string | null = pod ? new Date().toISOString() : null;
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
</script>

<div class="container">
  {#if error}
    <p class="error">Error: {error}</p>
  {:else if !pod}
    <p>Pod not found or still loading.</p>
  {:else}
    <PodHeader {pod} lastFetched={lastFetched} {isRefreshing} refreshIntervalSeconds={refreshInterval / 1000} />
    <Metadata {pod} />

    <h2>Labels</h2>
    <Labels labels={pod.metadata?.labels} />

    <h2>Containers ({pod.spec?.containers?.length ?? 0})</h2>
    <ContainersTable containers={pod.spec?.containers || []} containerStatuses={pod.status?.containerStatuses || []} />

    <RawJson {pod} />
    <Logs {pod} />
  {/if}
</div>

<style>
  .container { padding: 1rem; }
  .error { color: red }
</style>
