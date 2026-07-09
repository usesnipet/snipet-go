<script lang="ts">
	import { Button } from "$lib/components/ui/button";
	import PageLayout from "$lib/components/page-layout.svelte";
	import KnowledgeInfo from "$lib/features/knowledge/components/knowledge-info.svelte";
	import KnowledgeItemTable from "$lib/features/knowledge/components/tables/knowledge-item-table.svelte";
	import { knowledgeService } from "$lib/features/knowledge/service";
	import ArrowLeftIcon from "@lucide/svelte/icons/arrow-left";
	import Loader2Icon from "@lucide/svelte/icons/loader-2";
	import RefreshCwIcon from "@lucide/svelte/icons/refresh-cw";
	import type { PageProps } from "./$types";

	let { params }: PageProps = $props();

	const knowledgeQuery = $derived(knowledgeService.findById(params.id));
	const knowledgeItemsQuery = $derived(knowledgeService.listItems(params.id));
	const syncMutation = $derived(knowledgeService.sync(params.id));

	const isSyncing = $derived(
		syncMutation.isPending || knowledgeQuery.data?.sync_status === "in_progress",
	);

	function handleSync() {
		syncMutation.mutate(undefined);
	}
</script>

<PageLayout
	title={knowledgeQuery.data?.name ?? "Knowledge"}
	description="View knowledge details and stored items."
>
	{#snippet actionsLeft()}
		<Button variant="outline" size="icon" href="/(protected)/knowledge" aria-label="Back to knowledge list">
			<ArrowLeftIcon />
		</Button>
	{/snippet}

	{#snippet actionsRight()}
		<Button onclick={handleSync} disabled={isSyncing}>
			{#if isSyncing}
				<Loader2Icon class="animate-spin" />
			{:else}
				<RefreshCwIcon />
			{/if}
			Sync
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
