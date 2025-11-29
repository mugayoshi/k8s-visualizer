<script module>
  import { defineMeta } from '@storybook/addon-svelte-csf';
  import Metadata from '../lib/components/Metadata.svelte';

  const samplePod = {
    metadata: {
      name: 'nginx-pod',
      namespace: 'default',
      uid: '1234-5678',
      creationTimestamp: '2025-11-29T09:00:00Z',
      labels: { app: 'nginx', tier: 'frontend' },
      annotations: { 'deployment.kubernetes.io/revision': '3' }
    },
    spec: {
      nodeName: 'node-1'
    },
    status: {
      phase: 'Running'
    }
  };

  const { Story } = defineMeta({
    title: 'Components/Metadata',
    component: Metadata,
    tags: ['autodocs'],
    argTypes: {
      pod: { control: 'object' }
    }
  });
</script>

<Story name="Default" args={{ pod: samplePod }} />
