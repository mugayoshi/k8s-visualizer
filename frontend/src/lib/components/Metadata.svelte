<script lang="ts">
  import type { KubernetesPod } from '$lib/types/kubernetes';
  import { statusClass } from '$lib/utils/statusClass';
  export let pod: KubernetesPod;
</script>

<div class="meta">
  <div><strong>Namespace:</strong> {pod.metadata?.namespace}</div>
  <div><strong>Status:</strong> <span class={statusClass(pod.status?.phase)}>{pod.status?.phase || 'Unknown'}</span></div>
  <div><strong>Node:</strong> {pod.spec?.nodeName || '-'}</div>
  <div><strong>Created:</strong> {pod.metadata?.creationTimestamp}</div>
</div>

<style>
  .meta { display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 0.5rem; margin-bottom: 1rem; }
  .status-running { color: green; font-weight:600 }
  .status-pending { color: orange; font-weight:600 }
  .status-failed { color: red; font-weight:600 }
  .status-succeeded { color: teal; font-weight:600 }
  .status-unknown { color: gray; font-weight:600 }
</style>
