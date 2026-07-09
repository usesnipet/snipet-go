<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { ModeWatcher } from "mode-watcher";
	import { QueryClientProvider } from '@tanstack/svelte-query'
  import { Toaster } from "$lib/components/ui/sonner/index.js";
	import { queryClient } from '$lib/query-client';

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<ModeWatcher />
<Toaster />
<QueryClientProvider client={queryClient}>
	{@render children()}
</QueryClientProvider>
