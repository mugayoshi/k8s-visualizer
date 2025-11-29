<script lang="ts">
  interface Container {
    name: string;
    image: string;
  }

  interface ContainerStatus {
    name: string;
    ready: boolean;
  }

  export let containers: Container[] = [];
  export let containerStatuses: ContainerStatus[] | undefined = [];
</script>

<table class="containers">
  <thead>
    <tr><th>Name</th><th>Image</th><th>Ready</th></tr>
  </thead>
  <tbody>
    {#each containers as c}
      <tr>
        <td>{c.name}</td>
        <td>{c.image}</td>
        <td>{containerStatuses?.find((cs) => cs.name === c.name)?.ready ? 'yes' : 'no'}</td>
      </tr>
    {/each}
  </tbody>
</table>

<style>
  .containers { width:100%; border-collapse: collapse; margin-bottom:1rem; }
  .containers th, .containers td { border:1px solid #e5e7eb; padding:0.5rem; text-align:left; }
</style>
