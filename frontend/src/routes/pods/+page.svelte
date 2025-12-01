<script lang="ts">
    import { onMount } from 'svelte';
    import type { Pod } from '$lib/types/kubernetes';
    import PodsTable from '$lib/components/PodsTable.svelte';



    let pods: Pod[] = [];
    let loading = true;
    let error: string | null = null;
    let ws: WebSocket | null = null;
    let connected = false;
    let wsUrl = import.meta.env.VITE_WEBSOCKET_URL || 'ws://localhost:8080/api/ws';
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

    <PodsTable {pods} {loading} {error} />
</div>

 