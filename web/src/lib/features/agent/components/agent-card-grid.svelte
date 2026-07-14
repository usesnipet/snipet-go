<script lang="ts">
	import ConfirmDialog from "$lib/components/confirm-dialog.svelte";
	import { Skeleton } from "$lib/components/ui/skeleton";
	import { cn } from "$lib/utils";
	import type { Agent } from "../schemas";
	import { agentService } from "../service";
	import AgentCard from "./agent-card.svelte";
	import AgentUpdateDialog from "./agent-update-dialog.svelte";

	type Props = {
		agents: Agent[];
		isLoading?: boolean;
		class?: string;
	};

	let { agents, isLoading = false, class: className }: Props = $props();

	let deleteConfirmDialogOpen = $state(false);
	let deleteConfirmDialogId = $state<string | null>(null);
	let updateDialogOpen = $state(false);
	let updateDialogAgent: Agent | null = $state(null);

	const deleteMutation = agentService.delete();

	function handleOpenUpdateDialog(agent: Agent) {
		updateDialogAgent = agent;
		updateDialogOpen = true;
	}

	function handleOpenDeleteConfirm(id: string) {
		deleteConfirmDialogId = id;
		deleteConfirmDialogOpen = true;
	}

	function handleDelete() {
		if (!deleteConfirmDialogId) return;
		deleteMutation.mutate(deleteConfirmDialogId);
		deleteConfirmDialogId = null;
	}
</script>

{#if isLoading}
	<div
		class={cn(
			"grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4",
			className,
		)}
	>
		<Skeleton class="h-40 w-full" />
		<Skeleton class="h-40 w-full" />
		<Skeleton class="h-40 w-full" />
		<Skeleton class="h-40 w-full" />
	</div>
{:else if agents.length === 0}
	<p class="text-muted-foreground text-sm">No agents yet.</p>
{:else}
	<div
		class={cn(
			"grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4",
			className,
		)}
	>
		{#each agents as agent (agent.id)}
			<AgentCard
				{agent}
				onEdit={() => handleOpenUpdateDialog(agent)}
				onDelete={() => handleOpenDeleteConfirm(agent.id)}
			/>
		{/each}
	</div>
{/if}

<AgentUpdateDialog bind:open={updateDialogOpen} agent={updateDialogAgent ?? undefined} />

<ConfirmDialog
	bind:open={deleteConfirmDialogOpen}
	title="Delete agent"
	description="Are you sure you want to delete this agent? This action cannot be undone."
	danger
	onConfirm={handleDelete}
/>
