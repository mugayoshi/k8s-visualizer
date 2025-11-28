<!-- src/routes/test-ws/+page.svelte -->
<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import type { WebSocketMessage } from '$lib/types/kubernetes';

  let ws: WebSocket | null = null;
  let connected = false;
  let messages: WebSocketMessage[] = [];
  let connectionStatus = 'Disconnected';
  let statusColor = 'bg-red-500';
  let wsUrl = 'ws://localhost:8080/api/ws';
  let messageToSend = '';
  let selectedAction = 'subscribe_pods';
  let selectedNamespace = '';
  let autoReconnect = true;
  let reconnectAttempts = 0;
  let maxReconnectAttempts = 5;

  function connect() {
    try {
      connectionStatus = 'Connecting...';
      statusColor = 'bg-yellow-500';
      
      ws = new WebSocket(wsUrl);

      ws.onopen = () => {
        console.log('WebSocket connected');
        connected = true;
        connectionStatus = 'Connected';
        statusColor = 'bg-green-500';
        reconnectAttempts = 0;

        // Send initial subscription
        sendMessage({
          action: 'subscribe_pods',
          namespace: '',
        });
      };

      ws.onmessage = (event) => {
        console.log('Received message:', event.data);
        try {
          const message: WebSocketMessage = JSON.parse(event.data);
          messages = [message, ...messages].slice(0, 50); // Keep last 50 messages
        } catch (error) {
          console.error('Error parsing message:', error);
          messages = [
            {
              type: 'error',
              action: 'parse_error',
              data: event.data,
              timestamp: new Date().toISOString(),
            },
            ...messages
          ].slice(0, 50);
        }
      };

      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        connectionStatus = 'Error';
        statusColor = 'bg-red-500';
      };

      ws.onclose = (event) => {
        console.log('WebSocket closed:', event.code, event.reason);
        connected = false;
        connectionStatus = 'Disconnected';
        statusColor = 'bg-red-500';
        ws = null;

        // Auto reconnect
        if (autoReconnect && reconnectAttempts < maxReconnectAttempts) {
          reconnectAttempts++;
          console.log(`Reconnecting... (${reconnectAttempts}/${maxReconnectAttempts})`);
          setTimeout(() => {
            connect();
          }, 3000);
        }
      };
    } catch (error) {
      console.error('Error connecting:', error);
      connectionStatus = 'Connection Failed';
      statusColor = 'bg-red-500';
    }
  }

  function disconnect() {
    autoReconnect = false;
    if (ws) {
      ws.close();
      ws = null;
    }
    connected = false;
    connectionStatus = 'Disconnected';
    statusColor = 'bg-red-500';
  }

  function sendMessage(data: any) {
    if (ws && connected) {
      const message = JSON.stringify(data);
      ws.send(message);
      console.log('Sent message:', message);
    } else {
      console.error('WebSocket is not connected');
    }
  }

  function sendCustomMessage() {
    if (messageToSend) {
      try {
        const data = JSON.parse(messageToSend);
        sendMessage(data);
      } catch (error) {
        console.error('Invalid JSON:', error);
        alert('Invalid JSON format');
      }
    }
  }

  function sendPredefinedMessage() {
    const data: any = {
      action: selectedAction,
    };

    if (selectedNamespace) {
      data.namespace = selectedNamespace;
    }

    sendMessage(data);
  }

  function clearMessages() {
    messages = [];
  }

  onMount(() => {
    // Auto-connect on mount
    connect();
  });

  onDestroy(() => {
    autoReconnect = false;
    disconnect();
  });

  function formatTimestamp(timestamp: string): string {
    return new Date(timestamp).toLocaleTimeString();
  }

  function getMessageColor(type: string): string {
    switch (type) {
      case 'pods':
        return 'border-blue-200 bg-blue-50';
      case 'nodes':
        return 'border-green-200 bg-green-50';
      case 'metrics':
        return 'border-purple-200 bg-purple-50';
      case 'error':
        return 'border-red-200 bg-red-50';
      default:
        return 'border-gray-200 bg-gray-50';
    }
  }
</script>

<div class="max-w-7xl mx-auto p-6">
  <h1 class="text-3xl font-bold mb-6">WebSocket Test</h1>

  <!-- Connection Status -->
  <div class="bg-white rounded-lg shadow p-6 mb-6">
    <div class="flex items-center justify-between mb-4">
      <div class="flex items-center gap-3">
        <div class="flex items-center gap-2">
          <div class="w-3 h-3 rounded-full {statusColor}"></div>
          <span class="font-semibold">{connectionStatus}</span>
        </div>
        {#if reconnectAttempts > 0}
          <span class="text-sm text-gray-600">
            (Reconnect attempt: {reconnectAttempts}/{maxReconnectAttempts})
          </span>
        {/if}
      </div>

      <div class="flex gap-2">
        {#if !connected}
          <button
            on:click={connect}
            class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
          >
            Connect
          </button>
        {:else}
          <button
            on:click={disconnect}
            class="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700"
          >
            Disconnect
          </button>
        {/if}
      </div>
    </div>

    <div class="space-y-2">
      <div>
        <label for="wsUrl" class="block text-sm font-medium mb-1">
          WebSocket URL:
        </label>
        <input
          id="wsUrl"
          type="text"
          bind:value={wsUrl}
          disabled={connected}
          class="w-full px-3 py-2 border border-gray-300 rounded disabled:bg-gray-100"
          placeholder="ws://localhost:8080/api/ws"
        />
      </div>

      <div class="flex items-center gap-2">
        <input
          id="autoReconnect"
          type="checkbox"
          bind:checked={autoReconnect}
          class="w-4 h-4"
        />
        <label for="autoReconnect" class="text-sm">
          Auto-reconnect on disconnect
        </label>
      </div>
    </div>
  </div>

  <!-- Send Messages -->
  <div class="bg-white rounded-lg shadow p-6 mb-6">
    <h2 class="text-xl font-semibold mb-4">Send Messages</h2>

    <!-- Predefined Messages -->
    <div class="mb-4">
      <h3 class="text-sm font-medium mb-2">Quick Actions:</h3>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <label for="action" class="block text-sm mb-1">Action:</label>
          <select
            id="action"
            bind:value={selectedAction}
            class="w-full px-3 py-2 border border-gray-300 rounded"
          >
            <option value="subscribe_pods">Subscribe to Pods</option>
            <option value="subscribe_nodes">Subscribe to Nodes</option>
            <option value="get_metrics">Get Metrics</option>
          </select>
        </div>

        <div>
          <label for="namespace" class="block text-sm mb-1">Namespace (optional):</label>
          <input
            id="namespace"
            type="text"
            bind:value={selectedNamespace}
            placeholder="default, kube-system, etc."
            class="w-full px-3 py-2 border border-gray-300 rounded"
          />
        </div>

        <div class="flex items-end">
          <button
            on:click={sendPredefinedMessage}
            disabled={!connected}
            class="w-full px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:bg-gray-400"
          >
            Send
          </button>
        </div>
      </div>
    </div>

    <!-- Custom JSON Message -->
    <div>
      <h3 class="text-sm font-medium mb-2">Custom JSON Message:</h3>
      <div class="flex gap-2">
        <textarea
          bind:value={messageToSend}
          placeholder='default'
          class="flex-1 px-3 py-2 border border-gray-300 rounded font-mono text-sm"
          rows="3"
        ></textarea>
        <button
          on:click={sendCustomMessage}
          disabled={!connected}
          class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:bg-gray-400"
        >
          Send Custom
        </button>
      </div>
    </div>
  </div>

  <!-- Messages Log -->
  <div class="bg-white rounded-lg shadow p-6">
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-xl font-semibold">Messages Log ({messages.length})</h2>
      <button
        on:click={clearMessages}
        class="px-3 py-1 text-sm bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
      >
        Clear
      </button>
    </div>

    {#if messages.length === 0}
      <div class="text-center py-12 text-gray-500">
        <p>No messages received yet.</p>
        <p class="text-sm mt-2">Connect to WebSocket and messages will appear here.</p>
      </div>
    {:else}
      <div class="space-y-3 max-h-[600px] overflow-y-auto">
        {#each messages as message, index}
          <div class="border rounded p-4 {getMessageColor(message.type)}">
            <div class="flex justify-between items-start mb-2">
              <div class="flex gap-2">
                <span class="px-2 py-0.5 bg-white rounded text-xs font-semibold">
                  {message.type}
                </span>
                <span class="px-2 py-0.5 bg-white rounded text-xs font-semibold">
                  {message.action}
                </span>
                {#if message.namespace}
                  <span class="px-2 py-0.5 bg-white rounded text-xs">
                    {message.namespace}
                  </span>
                {/if}
              </div>
              <span class="text-xs text-gray-600">
                {formatTimestamp(message.timestamp)}
              </span>
            </div>

            <div class="bg-white rounded p-3 overflow-auto">
              <pre class="text-xs font-mono">{JSON.stringify(message.data, null, 2)}</pre>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  pre {
    max-height: 300px;
    white-space: pre-wrap;
    word-wrap: break-word;
  }
</style>