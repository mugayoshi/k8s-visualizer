<script module>
  import { defineMeta } from '@storybook/addon-svelte-csf';
  import PodHeader from '../lib/components/PodHeader.svelte';

  const samplePod = {
    metadata: {
      name: 'nginx-pod',
      namespace: 'default',
      uid: '1234-5678',
      creationTimestamp: '2025-11-29T09:00:00Z',
      labels: { app: 'nginx' }
    },
    spec: {
      nodeName: 'node-1'
    },
    status: {
      phase: 'Running'
    }
  };

  const { Story } = defineMeta({
    title: 'Components/PodHeader',
    component: PodHeader,
    tags: ['autodocs'],
    argTypes: {
      pod: { control: 'object' }
    }
  });
</script>

<Story name="Default" args={{ pod: samplePod }} />
