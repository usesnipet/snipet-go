<script lang="ts">
	import PageLayout from "$lib/components/page-layout.svelte";
	import KnowledgeTable from "$lib/features/knowledge/components/table/knowledge-table.svelte";
	import { knowledgeService } from "$lib/features/knowledge/service";

  const listQuery = knowledgeService.list();
  const knowledges = $derived(listQuery.data ?? []);
</script>

<PageLayout
	title="Knowledge"
	description="Manage your knowledge."
>
  <KnowledgeTable {knowledges} isLoading={listQuery.isLoading} />
</PageLayout>
