import { getAuthHeader, getUserInfo } from '$lib/auth';
import type { User } from '$lib/models';

export const load = ({ fetch }) => getUserInfo(fetch);
