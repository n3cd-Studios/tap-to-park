import { getAuthHeader } from '$lib/auth';
import type { User } from '$lib/models';

export const load = ({ fetch }) => fetch("auth/info", { headers: getAuthHeader() }) as User;
