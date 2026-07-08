<script lang="ts">
	import { Button } from "$lib/components/ui/button";
	import PageLayout from "$lib/components/page-layout.svelte";
	import KnowledgeInfo from "$lib/features/knowledge/components/knowledge-info.svelte";
	import KnowledgeItemTable from "$lib/features/knowledge/components/tables/knowledge-item-table.svelte";
	import { knowledgeService } from "$lib/features/knowledge/service";
	import ArrowLeftIcon from "@lucide/svelte/icons/arrow-left";
	import Loader2Icon from "@lucide/svelte/icons/loader-2";
	import RefreshCwIcon from "@lucide/svelte/icons/refresh-cw";
	import { toast } from "svelte-sonner";
	import type { PageProps } from "./$types";

	type Props = PageProps;

	let { params }: Props = $props();
	const id = $derived(params.id);

	const knowledgeQuery = $derived(knowledgeService.findById(id));
	const knowledgeItemsQuery = $derived(knowledgeService.listItems(id));
	const syncMutation = $derived(knowledgeService.sync(id));

	const knowledge = $derived(knowledgeQuery.data);
	const knowledgeItems = $derived(knowledgeItemsQuery.data ?? []);
	const isSyncing = $derived(
		syncMutation.isPending || knowledge?.sync_status === "in_progress",
	);

	function handleSync() {
		syncMutation.mutate(undefined, {
			onSuccess: () => {
				toast.success("Knowledge sync started.");
			},
			onError: (error) => {
				toast.error(error.message);
			},
		});
	}
</script>

<PageLayout
	title={knowledge?.name ?? "Knowledge"}
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
			<KnowledgeInfo knowledge={knowledge} isLoading={knowledgeQuery.isLoading} />
		</div>

		<div class="h-full min-w-0 flex-1 basis-0">
			<KnowledgeItemTable
				items={knowledgeItems}
				isLoading={knowledgeItemsQuery.isLoading}
				class="h-full min-w-0"
			/>
		</div>
	</div>
</PageLayout>
