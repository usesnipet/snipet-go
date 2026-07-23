<script lang="ts">
	import ChatComposer from "$lib/features/chat/components/chat-composer.svelte";
	import { appService } from "$lib/features/app/service";

	const configQuery = appService.config();
	const config = $derived(configQuery.data);
</script>

<ChatComposer title={config?.inherit_client_name ?? "Snipet"} />
