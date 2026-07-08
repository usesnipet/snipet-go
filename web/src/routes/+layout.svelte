<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { ModeWatcher } from "mode-watcher";
	import { QueryClientProvider, QueryClient } from '@tanstack/svelte-query'
	import { browser } from '$app/environment';

  const queryClient = new QueryClient({
		defaultOptions: {
      queries: {
        enabled: browser,
      },
    },
	})

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<ModeWatcher />
<QueryClientProvider client={queryClient}>
	{@render children()}
</QueryClientProvider>
