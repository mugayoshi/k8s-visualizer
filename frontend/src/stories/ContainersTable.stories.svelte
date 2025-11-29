<script module>
  import { defineMeta } from '@storybook/addon-svelte-csf';
  import ContainersTable from '../lib/components/ContainersTable.svelte';

  const sampleContainers = [
    { name: 'nginx', image: 'nginx:1.21' },
    { name: 'sidecar', image: 'busybox:latest' },
  ];
  const sampleStatuses = [
    { name: 'nginx', ready: true },
    { name: 'sidecar', ready: false },
  ];

  const { Story } = defineMeta({
    title: 'Components/ContainersTable',
    component: ContainersTable,
    tags: ['autodocs'],
    argTypes: {
      containers: { control: 'object' },
      containerStatuses: { control: 'object' },
    }
  });
</script>

<Story name="Default" args={{ containers: sampleContainers, containerStatuses: sampleStatuses }} />
