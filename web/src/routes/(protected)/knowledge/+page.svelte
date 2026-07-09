<script lang="ts">
	import KnowledgeCreateDialog from "$lib/features/knowledge/components/knowledge-create-dialog.svelte";
	import PageLayout from "$lib/components/page-layout.svelte";
	import KnowledgeTable from "$lib/features/knowledge/components/tables/knowledge-table.svelte";
	import { knowledgeService } from "$lib/features/knowledge/service";

  const listQuery = knowledgeService.list();
  const knowledges = $derived(listQuery.data ?? []);
</script>

<PageLayout
	title="Knowledge"
	description="Manage your knowledge."
>
	{#snippet actionsRight()}
		<KnowledgeCreateDialog />
	{/snippet}
  <KnowledgeTable {knowledges} isLoading={listQuery.isLoading} />
</PageLayout>
