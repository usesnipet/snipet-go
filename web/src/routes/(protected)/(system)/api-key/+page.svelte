<script lang="ts">
	import PageLayout from "$lib/components/page-layout.svelte";
	import ApiKeyCreateDialog from "$lib/features/api-key/components/api-key-create-dialog.svelte";
	import ApiKeyTable from "$lib/features/api-key/components/api-key-table.svelte";
	import { apiKeyService } from "$lib/features/api-key/service";
	import PlusIcon from "@lucide/svelte/icons/plus";

	const listQuery = apiKeyService.list();
	const apiKeys = $derived(listQuery.data ?? []);
</script>

<PageLayout title="API Keys" description="Manage your API keys.">
	{#snippet actionsRight()}
		<ApiKeyCreateDialog>
			{#snippet trigger()}
				<PlusIcon />
				Create API key
			{/snippet}
		</ApiKeyCreateDialog>
	{/snippet}
	<ApiKeyTable {apiKeys} isLoading={listQuery.isLoading} />
</PageLayout>
