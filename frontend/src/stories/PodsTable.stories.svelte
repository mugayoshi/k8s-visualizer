<script module>
  import { defineMeta } from '@storybook/addon-svelte-csf';
  import PodsTable from '../lib/components/PodsTable.svelte';

  const samplePods = [
    {
      name: 'web-6f7c8d9b6f-abcde',
      namespace: 'default',
      status: 'Running',
      status_detail: 'Ready: 1/1',
      ready_containers: '1/1',
      node: 'node-1',
      created: '2025-11-27T10:00:00Z',
      labels: { app: 'web' },
      containers: [{ name: 'web', image: 'example/web:latest', ready: true }],
      restart_count: 0,
      cpu_request: '100m',
      cpu_limit: '250m',
      memory_request: '128Mi',
      memory_limit: '256Mi',
      container_count: 1,
      age: '2h'
    },
    {
      name: 'db-0',
      namespace: 'databases',
      status: 'Pending',
      status_detail: 'ContainerCreating',
      ready_containers: '0/1',
      node: 'node-2',
      created: '2025-11-27T11:30:00Z',
      labels: { app: 'db' },
      containers: [{ name: 'db', image: 'example/db:1.0', ready: false }],
      restart_count: 1,
      cpu_request: '200m',
      cpu_limit: '500m',
      memory_request: '256Mi',
      memory_limit: '512Mi',
      container_count: 1,
      age: '30m'
    }
  ];

  const { Story } = defineMeta({
    title: 'components/PodsTable',
    component: PodsTable,
    tags: ['autodocs'],
    argTypes: {
      pods: { control: 'object' },
      loading: { control: 'boolean' },
      error: { control: 'text' }
    }
  });
</script>

<Story name="Default" args={{ pods: samplePods }} />

<Story name="Loading" args={{ pods: [], loading: true }} />

<Story name="Empty" args={{ pods: [] }} />
