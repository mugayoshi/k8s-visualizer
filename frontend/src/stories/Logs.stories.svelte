<script module>
  import { defineMeta } from '@storybook/addon-svelte-csf';
  import Logs from '../lib/components/Logs.svelte';

  const samplePod = {
    metadata: { name: 'nginx-pod', namespace: 'default' },
    spec: { nodeName: 'node-1' },
    status: { phase: 'Running' }
  };

  const { Story } = defineMeta({
    title: 'Components/Logs',
    component: Logs,
    tags: ['autodocs'],
    argTypes: {
      pod: { control: 'object' }
    }
  });
</script>

<Story name="Default" args={{ pod: samplePod }} />
