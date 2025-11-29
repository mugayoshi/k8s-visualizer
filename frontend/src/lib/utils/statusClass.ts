export function statusClass(status: string | undefined): string {
  const s = (status || '').toLowerCase();
  if (s.includes('running')) return 'status-running';
  if (s.includes('pending')) return 'status-pending';
  if (s.includes('failed')) return 'status-failed';
  if (s.includes('succeeded') || s.includes('completed') || s.includes('success')) return 'status-succeeded';
  return 'status-unknown';
}
