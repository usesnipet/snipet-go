<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import ConfirmDialog from "$lib/components/confirm-dialog.svelte";
	import { Button } from "$lib/components/ui/button";
	import * as DropdownMenu from "$lib/components/ui/dropdown-menu";
	import { cn } from "$lib/utils";
	import EllipsisIcon from "@lucide/svelte/icons/ellipsis";
	import PencilIcon from "@lucide/svelte/icons/pencil";
	import TrashIcon from "@lucide/svelte/icons/trash";
	import type { MockSession } from "../stores/mock-sessions.svelte";
	import { mockSessions } from "../stores/mock-sessions.svelte";
	import RenameSessionDialog from "./rename-session-dialog.svelte";
	import { resolve } from "$app/paths";

	type Props = {
		session: MockSession;
		active?: boolean;
	};

	let { session, active = false }: Props = $props();

	let menuOpen = $state(false);
	let renameOpen = $state(false);
	let deleteOpen = $state(false);

	const label = $derived(String(session.metadata.name ?? "New chat"));

	function handleSelect() {
		goto(resolve("/(chat)/session/[session_id]", { session_id: session.id }));
	}

	function handleDelete() {
		const wasActive = page.params.session_id === session.id;
		mockSessions.delete(session.id);
		if (wasActive) {
			goto(resolve("/(chat)"));
		}
	}
</script>

<div
	class={cn(
		"group relative flex w-full items-center rounded-lg text-sm",
		active ? "bg-muted font-medium" : "hover:bg-muted/60",
	)}
>
	<button
		type="button"
		class="min-w-0 flex-1 truncate py-2 pr-9 pl-3 text-left"
		onclick={handleSelect}
	>
		{label}
	</button>

	<DropdownMenu.Root bind:open={menuOpen}>
		<DropdownMenu.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					variant="ghost"
					size="icon-sm"
					class={cn(
						"absolute right-1 opacity-0 transition-opacity group-hover:opacity-100",
						menuOpen && "opacity-100",
					)}
					aria-label="Session actions"
				>
					<EllipsisIcon />
				</Button>
			{/snippet}
		</DropdownMenu.Trigger>
		<DropdownMenu.Content align="end">
			<DropdownMenu.Group>
				<DropdownMenu.Item
					onSelect={() => {
						renameOpen = true;
					}}
				>
					<PencilIcon />
					Rename
				</DropdownMenu.Item>
				<DropdownMenu.Item
					variant="destructive"
					onSelect={() => {
						deleteOpen = true;
					}}
				>
					<TrashIcon />
					Delete
				</DropdownMenu.Item>
			</DropdownMenu.Group>
		</DropdownMenu.Content>
	</DropdownMenu.Root>
</div>

{#if renameOpen}
	<RenameSessionDialog
		bind:open={renameOpen}
		sessionId={session.id}
		initialName={label}
	/>
{/if}

<ConfirmDialog
	bind:open={deleteOpen}
	title="Delete session"
	description="Are you sure you want to delete this session? This action cannot be undone."
	danger
	onConfirm={handleDelete}
/>
