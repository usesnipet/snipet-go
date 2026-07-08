<script lang="ts">
	import PageLayout from "$lib/components/page-layout.svelte";
	import ApiKeyTable from "$lib/features/api-key/components/api-key-table.svelte";
	import { apiKeyService } from "$lib/features/api-key/service";

  const listQuery = apiKeyService.list();
  const apiKeys = $derived(listQuery.data ?? []);
</script>

<PageLayout
	title="API Keys"
	description="Manage your API keys."
>
  <ApiKeyTable {apiKeys} isLoading={listQuery.isLoading} />
</PageLayout>
