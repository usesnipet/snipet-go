<script lang="ts">
	import * as Card from "$lib/components/ui/card";
	import { Skeleton } from "$lib/components/ui/skeleton";
	import TableTimeField from "$lib/components/flex-table/table-time-field.svelte";
	import { Badge, type BadgeVariant } from "$lib/components/ui/badge";
	import { JsonConfigCard } from "$lib/components/json-config";
	import type { Knowledge, KnowledgeIndex, KnowledgeSyncStatus } from "../schemas";
	import Button from "$lib/components/ui/button/button.svelte";
	import { LoaderIcon, PencilIcon, PlusIcon, RefreshCwIcon } from "@lucide/svelte";
	import KnowledgeUpdateDialog from "./knowledge-update-dialog.svelte";
	import { knowledgeService } from "../service";
	import { logger } from "$lib/logger";
	import KnowledgeIndexCreateDialog from "./knowledge-index-create-dialog.svelte";
	import KnowledgeIndexUpdateDialog from "./knowledge-index-update-dialog.svelte";
	import { ScrollArea } from "$lib/components/ui/scroll-area";

	type Props = {
		knowledge?: Knowledge;
		isLoading?: boolean;
	};

	let { knowledge, isLoading = false }: Props = $props();

	const knowledgeIndexesQuery = $derived(knowledgeService.listIndexes(knowledge?.id));
	const knowledgeIndexes = $derived(knowledgeIndexesQuery.data ?? []);
	const isIndexesLoading = $derived(isLoading || knowledgeIndexesQuery.isLoading);

	let updateKnowledgeDialogOpen = $state(false);
	let createIndexDialogOpen = $state(false);
	let editIndexDialogOpen = $state(false);
	let editingIndex = $state<KnowledgeIndex | undefined>(undefined);

	function syncStatusVariant(status: KnowledgeSyncStatus | null | undefined): BadgeVariant {
		if (status === "success") return "success";
		if (status === "in_progress") return "secondary";
		return "destructive";
	}

	function handleUpdateKnowledge() {
		updateKnowledgeDialogOpen = true;
	}

	function handleCreateIndex() {
		createIndexDialogOpen = true;
	}

	function handleEditIndex(index: KnowledgeIndex) {
		editingIndex = index;
		editIndexDialogOpen = true;
	}

	const syncMutation = knowledgeService.sync();

	const isSyncing = $derived(
		syncMutation.isPending || knowledge?.sync_status === "in_progress",
	);

	function handleSync() {
		if (!knowledge) {
			logger.warn("Knowledge not found when syncing on knowledge info");
			return;
		}

		syncMutation.mutate({ id: knowledge.id });
	}
</script>

<Card.Root class="h-full w-full min-w-0">
	<Card.Header>
		<div class="flex justify-between items-center">
			<Card.Title>
				{#if isLoading}
					<Skeleton class="h-6 w-3/4" />
				{:else}
					{knowledge?.name}
				{/if}
			</Card.Title>
			<div>
				<Button size="sm" onclick={handleUpdateKnowledge}>
					<PencilIcon />
					Edit
				</Button>
				<Button size="sm" onclick={handleSync} disabled={isSyncing}>
					{#if isSyncing}
						<LoaderIcon class="animate-spin" />
					{:else}
						<RefreshCwIcon />
					{/if}
					Sync
				</Button>
			</div>
		</div>
		<div>
			{#if isLoading}
				<Skeleton class="h-4 w-full" />
			{:else if knowledge?.description}
				<Card.Description>{knowledge.description}</Card.Description>
			{/if}
		</div>
	</Card.Header>

	<Card.Content class="space-y-4">
		<div class="space-y-1">
			<p class="text-muted-foreground text-sm">Driver</p>
			{#if isLoading}
				<Skeleton class="h-4 w-1/2" />
			{:else}
				<p class="text-sm font-medium">{knowledge?.driver}</p>
			{/if}
		</div>
		<div class="flex flex-col gap-1">
			<p class="text-muted-foreground text-sm">Driver Configuration</p>
			{#if isLoading}
				<Skeleton class="h-16 w-full" />
			{:else}
				<JsonConfigCard data={knowledge?.configuration} />
			{/if}
		</div>

		<div class="space-y-1">
			<p class="text-muted-foreground text-sm">Last Synced At</p>
			{#if isLoading}
				<Skeleton class="h-4 w-2/3" />
			{:else if knowledge?.last_synced_at}
				<p class="text-sm font-medium">
					<TableTimeField date={knowledge.last_synced_at} />
				</p>
			{:else}
				<p class="text-muted-foreground text-sm">Never synced</p>
			{/if}
		</div>

		<div class="space-y-1">
			<p class="text-muted-foreground text-sm">Sync Status</p>
			{#if isLoading}
				<Skeleton class="h-5 w-24" />
			{:else if knowledge?.sync_status}
				<Badge variant={syncStatusVariant(knowledge.sync_status)}>{knowledge.sync_status}</Badge>
			{:else}
				<p class="text-muted-foreground text-sm">—</p>
			{/if}
		</div>

		{#if !isLoading && knowledge?.sync_error}
			<div class="space-y-1">
				<p class="text-muted-foreground text-sm">Sync Error</p>
				<p class="text-destructive text-sm">{knowledge.sync_error}</p>
			</div>
		{/if}

		<div class="space-y-2">
			<div class="flex items-center justify-between">
				<p class="text-muted-foreground text-sm">Indexes</p>
				<Button size="sm" onclick={handleCreateIndex} disabled={isLoading}>
					<PlusIcon />
					Create Index
				</Button>
			</div>

			{#if isIndexesLoading}
				<div class="space-y-2">
					<Skeleton class="h-14 w-full rounded-lg" />
					<Skeleton class="h-14 w-full rounded-lg" />
				</div>
			{:else}
				<ScrollArea class="max-h-52 rounded-lg border">
					<div class="divide-y">
						{#each knowledgeIndexes as index (index.id)}
							<div class="flex items-center gap-3 px-3 py-2.5">
								<div class="min-w-0 flex-1 space-y-0.5">
									<div class="flex items-center gap-2">
										<p class="truncate text-sm font-medium">{index.name}</p>
										<Badge variant="secondary">{index.driver}</Badge>
									</div>
								</div>
								<Button
									variant="ghost"
									size="icon-sm"
									onclick={() => handleEditIndex(index)}
									aria-label="Edit index {index.name}"
								>
									<PencilIcon />
								</Button>
							</div>
						{:else}
							<p class="text-muted-foreground px-3 py-6 text-center text-sm">No indexes yet</p>
						{/each}
					</div>
				</ScrollArea>
			{/if}
		</div>
	</Card.Content>
</Card.Root>

<KnowledgeUpdateDialog
  bind:open={updateKnowledgeDialogOpen}
	knowledge={knowledge}
/>
{#if knowledge}
	<KnowledgeIndexCreateDialog
		knowledge={knowledge}
		bind:open={createIndexDialogOpen}
	/>
	{#if editingIndex}
		<KnowledgeIndexUpdateDialog
			knowledge={knowledge}
			index={editingIndex}
			bind:open={editIndexDialogOpen}
		/>
	{/if}
{/if}