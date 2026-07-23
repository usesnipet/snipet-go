<script lang="ts">
	import { Button } from "$lib/components/ui/button";
	import PageLayout from "$lib/components/page-layout.svelte";
	import KnowledgeInfo from "$lib/features/knowledge/components/knowledge-info.svelte";
	import KnowledgeItemTable from "$lib/features/knowledge/components/tables/knowledge-item-table.svelte";
	import { knowledgeService } from "$lib/features/knowledge/service";
	import ArrowLeftIcon from "@lucide/svelte/icons/arrow-left";
	import type { PageProps } from "./$types";

	let { params }: PageProps = $props();

	const knowledgeQuery = $derived(knowledgeService.findById(params.id));
	const knowledgeItemsQuery = $derived(knowledgeService.listItems(params.id));
</script>

<PageLayout
	title={`Knowledge • ${knowledgeQuery.data?.name}`}
	description="View knowledge details and stored items."
>
	{#snippet actionsLeft()}
		<Button variant="outline" size="icon" href="/manage/(protected)/(system)/knowledge" aria-label="Back to knowledge list">
			<ArrowLeftIcon />
		</Button>
	{/snippet}

	<div class="flex h-full min-h-0 w-full gap-4">
		<div class="h-full w-[30%] shrink-0 basis-[30%]">
			<KnowledgeInfo knowledge={knowledgeQuery.data} isLoading={knowledgeQuery.isLoading} />
		</div>

		<div class="h-full min-w-0 flex-1 basis-0">
			<KnowledgeItemTable
				items={knowledgeItemsQuery.data ?? []}
				isLoading={knowledgeItemsQuery.isLoading}
				class="h-full min-w-0"
			/>
		</div>
	</div>
</PageLayout>
