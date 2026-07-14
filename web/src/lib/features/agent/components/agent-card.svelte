<script lang="ts">
	import * as Card from "$lib/components/ui/card";
	import * as DropdownMenu from "$lib/components/ui/dropdown-menu";
	import { Button } from "$lib/components/ui/button";
	import BotIcon from "@lucide/svelte/icons/bot";
	import EllipsisIcon from "@lucide/svelte/icons/ellipsis";
	import PencilIcon from "@lucide/svelte/icons/pencil";
	import TrashIcon from "@lucide/svelte/icons/trash";
	import { cn } from "$lib/utils";
	import type { Agent } from "../schemas";

	type Props = {
		agent: Agent;
		onEdit: () => void;
		onDelete: () => void;
		class?: string;
	};

	let { agent, onEdit, onDelete, class: className }: Props = $props();

	let menuOpen = $state(false);

	const truncatedDescription = $derived.by(() => {
		const description = agent.description ?? "";
		if (description.length <= 200) return description;
		return `${description.slice(0, 200)}…`;
	});
</script>

<Card.Root class={cn("group", className)}>
	<Card.Header class="grid-cols-1">
		<div class="flex items-start justify-between gap-2">
			<div class="bg-muted flex size-10 shrink-0 items-center justify-center rounded-xl">
				<BotIcon class="text-muted-foreground" />
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
							aria-label="Agent actions"
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
		<Card.Title class="truncate">{agent.name}</Card.Title>
		{#if truncatedDescription}
			<Card.Description class="line-clamp-4">{truncatedDescription}</Card.Description>
		{/if}
	</Card.Header>
</Card.Root>
