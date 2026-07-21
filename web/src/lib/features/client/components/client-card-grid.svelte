<script lang="ts">
	import ConfirmDialog from "$lib/components/confirm-dialog.svelte";
	import { Skeleton } from "$lib/components/ui/skeleton";
	import { cn } from "$lib/utils";
	import type { Client } from "../schemas";
	import { clientService } from "../service";
	import ClientCard from "./client-card.svelte";
	import ClientUpdateDialog from "./client-update-dialog.svelte";

	type Props = {
		clients: Client[];
		isLoading?: boolean;
		class?: string;
	};

	let { clients, isLoading = false, class: className }: Props = $props();

	let deleteConfirmDialogOpen = $state(false);
	let deleteConfirmDialogCode = $state<string | null>(null);
	let updateDialogOpen = $state(false);
	let updateDialogClient: Client | null = $state(null);

	const deleteMutation = clientService.delete();

	function handleOpenUpdateDialog(client: Client) {
		updateDialogClient = client;
		updateDialogOpen = true;
	}

	function handleOpenDeleteConfirm(code: string) {
		deleteConfirmDialogCode = code;
		deleteConfirmDialogOpen = true;
	}

	function handleDelete() {
		if (!deleteConfirmDialogCode) return;
		deleteMutation.mutate(deleteConfirmDialogCode);
		deleteConfirmDialogCode = null;
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
{:else if clients.length === 0}
	<p class="text-muted-foreground text-sm">No clients yet.</p>
{:else}
	<div
		class={cn(
			"grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4",
			className,
		)}
	>
		{#each clients as client (client.id)}
			<ClientCard
				{client}
				onEdit={() => handleOpenUpdateDialog(client)}
				onDelete={() => handleOpenDeleteConfirm(client.code)}
			/>
		{/each}
	</div>
{/if}

<ClientUpdateDialog bind:open={updateDialogOpen} client={updateDialogClient ?? undefined} />

<ConfirmDialog
	bind:open={deleteConfirmDialogOpen}
	title="Delete client"
	description="Are you sure you want to delete this client? This action cannot be undone."
	danger
	onConfirm={handleDelete}
/>
