<script module>
  import { defineMeta } from '@storybook/addon-svelte-csf';
  import RawJson from '../lib/components/RawJson.svelte';

  const samplePod = {
    metadata: {
      name: 'nginx-pod',
      namespace: 'default',
      uid: '1234-5678',
      creationTimestamp: '2025-11-29T09:00:00Z',
      labels: { app: 'nginx' }
    },
    spec: {
      nodeName: 'node-1',
      containers: [{ name: 'nginx', image: 'nginx:1.21' }]
    },
    status: {
      phase: 'Running',
      containerStatuses: [{ name: 'nginx', ready: true }]
    }
  };

  const { Story } = defineMeta({
    title: 'Components/RawJson',
    component: RawJson,
    tags: ['autodocs'],
    argTypes: {
      pod: { control: 'object' }
    }
  });
</script>

<Story name="Default" args={{ pod: samplePod }} />
