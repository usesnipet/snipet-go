<script lang="ts">
	import ApiKeyGuard from "$lib/features/api-key/components/api-key-guard.svelte";

	let { children } = $props();
</script>

<ApiKeyGuard>
	{@render children()}
</ApiKeyGuard>
