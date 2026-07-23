<script lang="ts">
	import KnowledgeCreateDialog from "$lib/features/knowledge/components/knowledge-create-dialog.svelte";
	import PageLayout from "$lib/components/page-layout.svelte";
	import KnowledgeTable from "$lib/features/knowledge/components/tables/knowledge-table.svelte";
	import { knowledgeService } from "$lib/features/knowledge/service";
	import { PlusIcon } from "@lucide/svelte";

  const listQuery = knowledgeService.list();
  const knowledges = $derived(listQuery.data ?? []);
</script>

<PageLayout
	title="Knowledge"
	description="Manage your knowledge."
>
	{#snippet actionsRight()}
		<KnowledgeCreateDialog>
			{#snippet trigger()}
				<PlusIcon />
				Create knowledge
			{/snippet}
		</KnowledgeCreateDialog>
	{/snippet}
  <KnowledgeTable {knowledges} isLoading={listQuery.isLoading} />
</PageLayout>
