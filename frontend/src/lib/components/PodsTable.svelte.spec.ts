import { page } from 'vitest/browser';
import { describe, expect, it } from 'vitest';
import { render } from 'vitest-browser-svelte';
import PodsTable from './PodsTable.svelte';

describe('PodsTable component', () => {
  it('renders table with pods when data is provided', async () => {
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
      }
    ];

    render(PodsTable, { props: { pods: samplePods } });

    const table = page.getByRole('table');
    await expect.element(table).toBeInTheDocument();

    const podName = page.getByText('web-6f7c8d9b6f-abcde');
    await expect.element(podName).toBeInTheDocument();

    const status = page.getByText('Running');
    await expect.element(status).toBeInTheDocument();
  });

  it('renders loading state when loading is true', async () => {
    render(PodsTable, { props: { pods: [], loading: true } });

    const loadingText = page.getByText('Loading pods...');
    await expect.element(loadingText).toBeInTheDocument();
  });

  it('renders error state when error is provided', async () => {
    render(PodsTable, { props: { pods: [], error: 'Connection failed' } });

    const errorText = page.getByText('Error: Connection failed');
    await expect.element(errorText).toBeInTheDocument();
  });
});
