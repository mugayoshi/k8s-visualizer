<!-- src/routes/test/+page.svelte -->
<script lang="ts">
  import { onMount } from 'svelte';
  import { apiClient } from '$lib/api/client';
  import type { Pod } from '$lib/types/kubernetes';

  let loading = false;
  let error: string | null = null;
  let pods: Pod[] = [];
  let namespace = 'all';
  let healthStatus: any = null;

  async function checkHealth() {
    try {
      loading = true;
      error = null;
      healthStatus = await apiClient.healthCheck();
      console.log('Health check:', healthStatus);
    } catch (e: any) {
      error = e.message;
      console.error('Health check failed:', e);
    } finally {
      loading = false;
    }
  }

  async function fetchPods() {
    try {
      loading = true;
      error = null;
      const response = await apiClient.listPods(namespace);
      pods = response.pods;
      console.log('Fetched pods:', response);
    } catch (e: any) {
      error = e.message;
      console.error('Failed to fetch pods:', e);
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    // Check health first
    checkHealth();
  });
</script>

<div class="max-w-7xl mx-auto p-6">
  <h1 class="text-3xl font-bold mb-6">API Client Test</h1>

  <!-- Health Check Section -->
  <div class="bg-white rounded-lg shadow p-6 mb-6">
    <h2 class="text-xl font-semibold mb-4">Health Check</h2>
    
    <button
      on:click={checkHealth}
      disabled={loading}
      class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:bg-gray-400"
    >
      {loading ? 'Checking...' : 'Check Health'}
    </button>

    {#if healthStatus}
      <div class="mt-4 p-4 bg-green-50 rounded">
        <pre class="text-sm">{JSON.stringify(healthStatus, null, 2)}</pre>
      </div>
    {/if}
  </div>

  <!-- Pod Fetch Section -->
  <div class="bg-white rounded-lg shadow p-6 mb-6">
    <h2 class="text-xl font-semibold mb-4">Fetch Pods</h2>

    <div class="flex gap-4 mb-4">
      <div>
        <label for="namespace" class="block text-sm font-medium mb-2">
          Namespace:
        </label>
        <select
          id="namespace"
          bind:value={namespace}
          class="px-4 py-2 border border-gray-300 rounded"
        >
          <option value="all">All Namespaces</option>
          <option value="default">default</option>
          <option value="kube-system">kube-system</option>
          <option value="kube-public">kube-public</option>
          <option value="k8s-visualizer">k8s-visualizer</option>
        </select>
      </div>

      <div class="flex items-end">
        <button
          on:click={fetchPods}
          disabled={loading}
          class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:bg-gray-400"
        >
          {loading ? 'Loading...' : 'Fetch Pods'}
        </button>
      </div>
    </div>

    {#if error}
      <div class="mt-4 p-4 bg-red-50 border border-red-200 rounded text-red-800">
        <strong>Error:</strong> {error}
      </div>
    {/if}

    {#if pods.length > 0}
      <div class="mt-4">
        <p class="text-sm text-gray-600 mb-2">Found {pods.length} pod(s)</p>
        
        <div class="space-y-4">
          {#each pods as pod}
            <div class="border border-gray-200 rounded p-4">
              <div class="flex justify-between items-start mb-2">
                <div>
                  <h3 class="font-semibold text-lg">{pod.name}</h3>
                  <p class="text-sm text-gray-600">Namespace: {pod.namespace}</p>
                </div>
                <span
                  class="px-3 py-1 rounded-full text-sm font-semibold"
                  class:bg-green-100={pod.status.toLowerCase() === 'running'}
                  class:text-green-800={pod.status.toLowerCase() === 'running'}
                  class:bg-yellow-100={pod.status.toLowerCase() === 'pending'}
                  class:text-yellow-800={pod.status.toLowerCase() === 'pending'}
                  class:bg-red-100={pod.status.toLowerCase() === 'failed'}
                  class:text-red-800={pod.status.toLowerCase() === 'failed'}
                >
                  {pod.status}
                </span>
              </div>

              <div class="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <span class="text-gray-600">Node:</span>
                  <span class="ml-2 font-mono">{pod.node || 'N/A'}</span>
                </div>
                <div>
                  <span class="text-gray-600">Created:</span>
                  <span class="ml-2">{new Date(pod.created).toLocaleString()}</span>
                </div>
              </div>

              {#if pod.containers && pod.containers.length > 0}
                <div class="mt-3">
                  <p class="text-sm font-medium text-gray-700 mb-1">Containers:</p>
                  <div class="space-y-1">
                    {#each pod.containers as container}
                      <div class="flex items-center gap-2 text-sm">
                        <span
                          class="w-2 h-2 rounded-full"
                          class:bg-green-500={container.ready}
                          class:bg-gray-400={!container.ready}
                        ></span>
                        <span class="font-mono">{container.name}</span>
                        <span class="text-gray-500">({container.image})</span>
                      </div>
                    {/each}
                  </div>
                </div>
              {/if}

              {#if pod.labels && Object.keys(pod.labels).length > 0}
                <div class="mt-3">
                  <p class="text-sm font-medium text-gray-700 mb-1">Labels:</p>
                  <div class="flex flex-wrap gap-1">
                    {#each Object.entries(pod.labels) as [key, value]}
                      <span class="px-2 py-0.5 bg-gray-100 text-gray-700 text-xs rounded">
                        {key}: {value}
                      </span>
                    {/each}
                  </div>
                </div>
              {/if}
            </div>
          {/each}
        </div>
      </div>
    {:else if !loading && !error}
      <p class="mt-4 text-gray-500 text-center py-8">
        No pods fetched yet. Click "Fetch Pods" to load data.
      </p>
    {/if}

    {#if loading}
      <div class="mt-4 flex justify-center">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    {/if}
  </div>

  <!-- Raw Response Section -->
  {#if pods.length > 0}
    <div class="bg-white rounded-lg shadow p-6">
      <h2 class="text-xl font-semibold mb-4">Raw JSON Response</h2>
      <div class="bg-gray-50 rounded p-4 overflow-auto">
        <pre class="text-xs">{JSON.stringify(pods, null, 2)}</pre>
      </div>
    </div>
  {/if}
</div>

<style>
  pre {
    max-height: 400px;
    overflow: auto;
  }
</style>