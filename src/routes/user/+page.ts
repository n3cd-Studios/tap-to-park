import { getAuthHeader } from '$lib/auth';
import type { User } from '$lib/models';

export const load = ({ fetch }) =>
  fetch("auth/info", { headers: getAuthHeader() })
    .then(res => res.json() as User)
    .catch(_ => null);
