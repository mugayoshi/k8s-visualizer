// src/lib/stores/websocket.ts
import { writable } from 'svelte/store';
import type { WebSocketMessage } from '$lib/types/kubernetes';

export const wsConnected = writable(false);
export const wsMessages = writable<WebSocketMessage[]>([]);

class WebSocketClient {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 3000;

  connect(url: string = import.meta.env.VITE_WEBSOCKET_URL || 'ws://localhost:8080/api/ws') {
    try {
      this.ws = new WebSocket(url);

      this.ws.onopen = () => {
        console.log('WebSocket connected');
        wsConnected.set(true);
        this.reconnectAttempts = 0;

        // Subscribe to pod updates
        this.send({
          action: 'subscribe_pods',
          namespace: '',
        });
      };

      this.ws.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data);
          wsMessages.update(messages => [...messages, message]);
          
          // Handle different message types
          this.handleMessage(message);
        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
        }
      };

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
      };

      this.ws.onclose = () => {
        console.log('WebSocket disconnected');
        wsConnected.set(false);
        this.attemptReconnect(url);
      };
    } catch (error) {
      console.error('Error connecting to WebSocket:', error);
    }
  }

  private handleMessage(message: WebSocketMessage) {
    // Handle different message types
    switch (message.type) {
      case 'pods':
        // Update pod store based on action
        console.log('Pod update:', message.action, message.data);
        break;
      case 'metrics':
        // Update metrics store
        console.log('Metrics update:', message.data);
        break;
      default:
        console.log('Unknown message type:', message.type);
    }
  }

  private attemptReconnect(url: string) {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`);
      
      setTimeout(() => {
        this.connect(url);
      }, this.reconnectDelay);
    } else {
      console.error('Max reconnection attempts reached');
    }
  }

  send(data: any) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data));
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}

export const wsClient = new WebSocketClient();