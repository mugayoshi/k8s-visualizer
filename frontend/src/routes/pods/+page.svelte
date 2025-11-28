<script lang="ts">
    import { onMount } from 'svelte';
    import type { Pod } from '$lib/types/kubernetes';



    let pods: Pod[] = [];
    let loading = true;
    let error: string | null = null;
    let ws: WebSocket | null = null;
    let connected = false;
    let wsUrl = 'ws://localhost:8080/api/ws';
    let autoReconnect = true;
    let reconnectAttempts = 0;
    let maxReconnectAttempts = 5;

    function connect() {
        try {
            ws = new WebSocket(wsUrl);

            ws.onopen = () => {
                connected = true;
                reconnectAttempts = 0;
                // Subscribe to pods updates
                ws?.send(JSON.stringify({ action: 'subscribe_pods', namespace: '' }));
            };

            ws.onmessage = (event) => {
                try {
                    const message = JSON.parse(event.data);
                    if (message.type === 'pods') {
                        error = null;
                        loading = false;

                        // Handle initial pod list or individual updates
                        if (message.action === 'initial' && Array.isArray(message.data)) {
                            pods = message.data;
                        } else if (message.action === 'added' && message.data) {
                            // Add new pod if not already exists
                            const exists = pods.some(p => p.name === message.data.name && p.namespace === message.data.namespace);
                            if (!exists) {
                                pods = [...pods, message.data];
                            }
                        } else if (message.action === 'modified' && message.data) {
                            // Update existing pod
                            pods = pods.map(p =>
                                p.name === message.data.name && p.namespace === message.data.namespace
                                    ? { ...p, ...message.data }
                                    : p
                            );
                        } else if (message.action === 'deleted' && message.data) {
                            // Remove deleted pod
                            pods = pods.filter(p =>
                                !(p.name === message.data.name && p.namespace === message.data.namespace)
                            );
                        }
                    }
                } catch (err) {
                    error = 'Error parsing WebSocket message';
                }
            };

            ws.onerror = () => {
                error = 'WebSocket error';
            };

            ws.onclose = () => {
                connected = false;
                ws = null;
                if (autoReconnect && reconnectAttempts < maxReconnectAttempts) {
                    reconnectAttempts++;
                    setTimeout(connect, 3000);
                }
            };
        } catch (err) {
            error = 'WebSocket connection failed';
        }
    }

    function statusClass(status: string) {
        const s = (status || '').toLowerCase();
        if (s.includes('running')) return 'status-running';
        if (s.includes('pending')) return 'status-pending';
        if (s.includes('failed')) return 'status-failed';
        if (s.includes('succeeded') || s.includes('completed') || s.includes('success')) return 'status-succeeded';
        return 'status-unknown';
    }

    onMount(() => {
        connect();
        return () => {
            autoReconnect = false;
            if (ws) ws.close();
        };
    });
</script>

<div class="container">
    <h1>Kubernetes Pods</h1>

    {#if loading}
        <p>Loading pods...</p>
    {:else if error}
        <p class="error">Error: {error}</p>
    {:else if pods.length === 0}
        <p>No pods found</p>
    {:else}
        <table>
            <thead>
                <tr>
                    <th>Name</th>
                    <th>Namespace</th>
                    <th>Status</th>
                    <!-- <th>Ready Containers</th> -->
                    <th>Restarts</th>
                    <th>Age</th>
                    <th>Containers</th>
                    <th>CPU Request</th>
                    <th>CPU Limit</th>
                    <th>Memory Request</th>
                    <th>Memory Limit</th>
                </tr>
            </thead>
            <tbody>
                {#each pods as pod (pod.name)}
                    <tr>
                        <td>{pod.name}</td>
                        <td>{pod.namespace}</td>
                        <td>
                            <span class="status-badge {statusClass(pod.status)}" title={pod.status_detail ? `${pod.status_detail} · ${pod.ready_containers || ''}` : (pod.ready_containers || pod.status)}>
                                {pod.status}
                            </span>
                            {#if pod.status_detail}
                                <small>{pod.status_detail}</small>
                            {/if}
                        </td>
                        <!-- <td>{pod.ready_containers || '-'}</td> -->
                        <td>{pod.restart_count || 0}</td>
                        <td>{pod.age || '-'}</td>
                        <td>{pod.container_count || '-'}</td>
                        <td>{pod.cpu_request || '-'}</td>
                        <td>{pod.cpu_limit || '-'}</td>
                        <td>{pod.memory_request || '-'}</td>
                        <td>{pod.memory_limit || '-'}</td>
                    </tr>
                {/each}
            </tbody>
        </table>
    {/if}
</div>

<style>
    .container {
        padding: 20px;
    }

    table {
        width: 100%;
        border-collapse: collapse;
    }

    th,
    td {
        border: 1px solid #ddd;
        padding: 12px;
        text-align: left;
    }

    th {
        background-color: #f5f5f5;
        font-weight: bold;
    }

    tr:hover {
        background-color: #f9f9f9;
    }

    .error {
        color: red;
    }

    /* Status badges */
    .status-badge {
        display: inline-block;
        padding: 4px 8px;
        border-radius: 12px;
        color: white;
        font-weight: 600;
        font-size: 0.85em;
    }
    .status-running { background: #16a34a; }
    .status-pending { background: #f59e0b; }
    .status-failed { background: #ef4444; }
    .status-succeeded { background: #6b7280; }
    .status-unknown { background: #7c3aed; }
</style>