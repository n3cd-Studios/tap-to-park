import { getAuthHeader, getUserInfo } from '$lib/auth';

export const load = ({ fetch }) => getUserInfo(fetch);
