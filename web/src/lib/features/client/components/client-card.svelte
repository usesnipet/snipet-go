<script lang="ts">
	import * as Card from "$lib/components/ui/card";
	import * as DropdownMenu from "$lib/components/ui/dropdown-menu";
	import { Button } from "$lib/components/ui/button";
	import Building2Icon from "@lucide/svelte/icons/building-2";
	import EllipsisIcon from "@lucide/svelte/icons/ellipsis";
	import PencilIcon from "@lucide/svelte/icons/pencil";
	import TrashIcon from "@lucide/svelte/icons/trash";
	import { cn } from "$lib/utils";
	import type { Client } from "../schemas";
	import { CopyIcon } from "@lucide/svelte";
	import { toast } from "svelte-sonner";

	type Props = {
		client: Client;
		onEdit: () => void;
		onDelete: () => void;
		class?: string;
	};

	let { client, onEdit, onDelete, class: className }: Props = $props();

	let menuOpen = $state(false);

	function handleCopyCode() {
		navigator.clipboard.writeText(client.code);
		toast.success("Client code copied.");
	}
</script>

<Card.Root class={cn("group", className)}>
	<Card.Header class="grid-cols-1">
		<div class="flex items-start justify-between gap-2">
			<div class="bg-muted flex size-10 shrink-0 items-center justify-center rounded-xl">
				<Building2Icon class="text-muted-foreground" />
			</div>
			<DropdownMenu.Root bind:open={menuOpen}>
				<DropdownMenu.Trigger>
					{#snippet child({ props })}
						<Button
							{...props}
							variant="ghost"
							size="icon-sm"
							class={cn(
								"opacity-0 transition-opacity group-hover:opacity-100",
								menuOpen && "opacity-100",
							)}
							aria-label="Client actions"
						>
							<EllipsisIcon />
						</Button>
					{/snippet}
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="end">
					<DropdownMenu.Group>
						<DropdownMenu.Item onSelect={onEdit}>
							<PencilIcon />
							Edit
						</DropdownMenu.Item>
						<DropdownMenu.Item variant="destructive" onSelect={onDelete}>
							<TrashIcon />
							Delete
						</DropdownMenu.Item>
					</DropdownMenu.Group>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		</div>
		<Card.Title class="truncate">{client.name}</Card.Title>
		<Card.Description class="truncate font-mono text-xs flex items-center gap-2">
			{client.code}
			<Button
			  variant="outline"
				size="icon-sm"
				class="opacity-0 group-hover:opacity-100 transition-opacity"
				onclick={handleCopyCode}
				aria-label="Copy client code"
			>
				<CopyIcon class="size-3" />
			</Button>
		</Card.Description>
	</Card.Header>
</Card.Root>
