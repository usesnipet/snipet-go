import { loadAuth } from "$lib/features/api-key/stores/api-key-auth.svelte";

export const prerender = false;
export const ssr = false;

export function load() {
  loadAuth();
}