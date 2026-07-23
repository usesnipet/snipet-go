<script lang="ts">
	import ChatSidebar from "$lib/features/chat/components/chat-sidebar.svelte";
	import { appService } from "$lib/features/app/service";
	import InheritClientDisabled from "$lib/features/chat/components/inherit-client-disabled.svelte";
	import Loading from "$lib/components/loading.svelte";

	let { children } = $props();

	const configQuery = appService.config();
	const config = $derived(configQuery.data);
	const enabled = $derived(config?.inherit_client === true);
</script>
<div class="h-screen">
	{#if configQuery.isPending}
		<Loading />
	{:else if enabled}
		<div class="flex h-svh overflow-hidden">
			<ChatSidebar clientName={config?.inherit_client_name ?? "Snipet"} />
			<main class="min-h-0 min-w-0 flex-1 overflow-hidden">
				{@render children()}
			</main>
		</div>
	{:else}
		<InheritClientDisabled />
	{/if}
</div>
