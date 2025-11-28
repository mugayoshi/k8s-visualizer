// src/lib/utils/formatters.ts
import { formatDistanceToNow, format } from 'date-fns';

export function formatDate(date: string | Date): string {
  return format(new Date(date), 'MMM dd, yyyy HH:mm:ss');
}

export function formatRelativeTime(date: string | Date): string {
  return formatDistanceToNow(new Date(date), { addSuffix: true });
}

export function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

export function formatMemory(memStr: string): string {
  // Parse Kubernetes memory format (e.g., "1024Ki", "2Gi")
  const units: Record<string, number> = {
    'Ki': 1024,
    'Mi': 1024 * 1024,
    'Gi': 1024 * 1024 * 1024,
    'Ti': 1024 * 1024 * 1024 * 1024,
  };

  for (const [unit, multiplier] of Object.entries(units)) {
    if (memStr.endsWith(unit)) {
      const value = parseFloat(memStr.slice(0, -unit.length));
      return formatBytes(value * multiplier);
    }
  }

  return memStr;
}

export function formatCPU(cpuStr: string): string {
  // Parse Kubernetes CPU format (e.g., "500m", "2")
  if (cpuStr.endsWith('m')) {
    const value = parseFloat(cpuStr.slice(0, -1));
    return `${value}m`;
  }
  return `${parseFloat(cpuStr) * 1000}m`;
}