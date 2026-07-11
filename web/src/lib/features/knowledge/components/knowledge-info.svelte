<script lang="ts">
	import * as Card from "$lib/components/ui/card";
	import { Skeleton } from "$lib/components/ui/skeleton";
	import TableTimeField from "$lib/components/flex-table/table-time-field.svelte";
	import { Badge, type BadgeVariant } from "$lib/components/ui/badge";
	import { JsonConfigCard } from "$lib/components/json-config";
	import type { Knowledge, KnowledgeSyncStatus } from "../schemas";
	import Button from "$lib/components/ui/button/button.svelte";
	import { LoaderIcon, PencilIcon, RefreshCwIcon } from "@lucide/svelte";
	import KnowledgeUpdateDialog from "./knowledge-update-dialog.svelte";
	import { knowledgeService } from "../service";
	import { logger } from "$lib/logger";

	type Props = {
		knowledge?: Knowledge;
		isLoading?: boolean;
	};

	let { knowledge, isLoading = false }: Props = $props();

	let updateKnowledgeDialogOpen = $state(false);

	function syncStatusVariant(status: KnowledgeSyncStatus | null | undefined): BadgeVariant {
		if (status === "success") return "default";
		if (status === "in_progress") return "secondary";
		return "destructive";
	}

	function handleUpdateKnowledge() {
		updateKnowledgeDialogOpen = true;
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
	<Card.Header class="flex justify-between items-center">
		<div>
			<Card.Title>
				{#if isLoading}
					<Skeleton class="h-6 w-3/4" />
				{:else}
					{knowledge?.name}
				{/if}
			</Card.Title>
			{#if isLoading}
				<Skeleton class="h-4 w-full" />
			{:else if knowledge?.description}
				<Card.Description>{knowledge.description}</Card.Description>
			{/if}
		</div>
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
	</Card.Content>
</Card.Root>

<KnowledgeUpdateDialog
  bind:open={updateKnowledgeDialogOpen}
	knowledge={knowledge}
/>