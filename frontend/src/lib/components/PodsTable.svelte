<script lang="ts">
  import type { Pod } from '$lib/types/kubernetes';

  export let pods: Pod[] = [];
  export let loading: boolean = false;
  export let error: string | null = null;

  function statusClass(status: string) {
    const s = (status || '').toLowerCase();
    if (s.includes('running')) return 'status-running';
    if (s.includes('pending')) return 'status-pending';
    if (s.includes('failed')) return 'status-failed';
    if (s.includes('succeeded') || s.includes('completed') || s.includes('success')) return 'status-succeeded';
    return 'status-unknown';
  }
</script>

<div class="pods-table">
  {#if loading}
    <p>Loading pods...</p>
  {:else if error}
    <p class="error">Error: {error}</p>
  {:else if pods.length === 0}
    <p>No pods found</p>
  {:else}
    <table>
      <thead>
        <tr>
          <th>Name</th>
          <th>Namespace</th>
          <th>Status</th>
          <th>Restarts</th>
          <th>Age</th>
          <th>Containers</th>
          <th>CPU Request</th>
          <th>CPU Limit</th>
          <th>Memory Request</th>
          <th>Memory Limit</th>
        </tr>
      </thead>
      <tbody>
        {#each pods as pod (pod.name)}
          <tr>
            <td>{pod.name}</td>
            <td>{pod.namespace}</td>
            <td>
              <span class="status-badge {statusClass(pod.status)}" title={pod.status_detail ? `${pod.status_detail} · ${pod.ready_containers || ''}` : (pod.ready_containers || pod.status)}>
                {pod.status}
              </span>
              {#if pod.status_detail}
                <small>{pod.status_detail}</small>
              {/if}
            </td>
            <td>{pod.restart_count || 0}</td>
            <td>{pod.age || '-'}</td>
            <td>{pod.container_count || '-'}</td>
            <td>{pod.cpu_request || '-'}</td>
            <td>{pod.cpu_limit || '-'}</td>
            <td>{pod.memory_request || '-'}</td>
            <td>{pod.memory_limit || '-'}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>

<style>
  .pods-table {
    padding: 12px;
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th,
  td {
    border: 1px solid #ddd;
    padding: 10px;
    text-align: left;
  }

  th {
    background-color: #f5f5f5;
    font-weight: bold;
  }

  tr:hover {
    background-color: #f9f9f9;
  }

  .error {
    color: red;
  }

  .status-badge {
    display: inline-block;
    padding: 4px 8px;
    border-radius: 12px;
    color: white;
    font-weight: 600;
    font-size: 0.85em;
  }
  .status-running { background: #16a34a; }
  .status-pending { background: #f59e0b; }
  .status-failed { background: #ef4444; }
  .status-succeeded { background: #6b7280; }
  .status-unknown { background: #7c3aed; }
</style>
