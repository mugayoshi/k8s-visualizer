import type { PageLoad } from './$types';
import { apiClient } from '$lib/api/client';
import type { KubernetesPod } from '$lib/types/kubernetes';

export const load: PageLoad = async ({ params, url }) => {
  const name = params.name;
  const namespace = url.searchParams.get('namespace') || 'default';

  try {
    const pod = (await apiClient.getPod(namespace, name)) as unknown as KubernetesPod;
    return { pod };
  } catch (err) {
    return {
      pod: null,
      error: (err as Error)?.message || 'Failed to load pod',
    };
  }
};
