// src/lib/utils/colors.ts
export function getStatusColor(status: string): string {
  const statusLower = status.toLowerCase();
  
  switch (statusLower) {
    case 'running':
    case 'ready':
    case 'active':
      return 'text-green-600 bg-green-100';
    case 'pending':
    case 'creating':
      return 'text-yellow-600 bg-yellow-100';
    case 'failed':
    case 'error':
    case 'crashloopbackoff':
      return 'text-red-600 bg-red-100';
    case 'succeeded':
    case 'completed':
      return 'text-blue-600 bg-blue-100';
    case 'terminating':
      return 'text-orange-600 bg-orange-100';
    default:
      return 'text-gray-600 bg-gray-100';
  }
}