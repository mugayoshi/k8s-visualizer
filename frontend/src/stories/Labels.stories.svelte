<script module>
  import { defineMeta } from '@storybook/addon-svelte-csf';
  import Labels from '../lib/components/Labels.svelte';

  const sampleLabels = {
    app: 'nginx',
    tier: 'frontend',
    version: '1.21'
  };

  const { Story } = defineMeta({
    title: 'Components/Labels',
    component: Labels,
    tags: ['autodocs'],
    argTypes: {
      labels: { control: 'object' }
    }
  });
</script>

<Story name="Default" args={{ labels: sampleLabels }} />
