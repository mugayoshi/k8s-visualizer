<script module>
  import { defineMeta } from '@storybook/addon-svelte-csf';
  import Logs from '../lib/components/Logs.svelte';

  const sampleLogs = [
    '2025-11-29T10:00:00Z Starting nginx...',
    '2025-11-29T10:00:01Z Listening on port 80',
    '2025-11-29T10:00:02Z GET /index.html 200'
  ].join('\n');

  const { Story } = defineMeta({
    title: 'Components/Logs',
    component: Logs,
    tags: ['autodocs'],
    argTypes: {
      logs: { control: 'text' },
      loading: { control: 'boolean' },
      error: { control: 'text' }
    }
  });
</script>

<Story name="Default" args={{ logs: sampleLogs }} />
