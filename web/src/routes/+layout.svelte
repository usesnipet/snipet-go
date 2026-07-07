<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { ModeWatcher } from "mode-watcher";
	import { onMount } from 'svelte';
	import { loadAuth } from '$lib/store/auth.svelte';

	let { children } = $props();

	onMount(loadAuth);
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<ModeWatcher />
{@render children()}
