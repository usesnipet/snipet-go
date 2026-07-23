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
	import { goto } from "$app/navigation";
	import { resolve } from "$app/paths";
	import ClientCode from "./client-code.svelte";

	type Props = {
		client: Client;
		onEdit: () => void;
		onDelete: () => void;
		class?: string;
	};

	let { client, onEdit, onDelete, class: className }: Props = $props();

	let menuOpen = $state(false);

	function handleClick() {
		goto(resolve("/manage/(protected)/(client)/c/[code]", { code: client.code }));
	}
</script>

<Card.Root class={cn("group cursor-pointer", className)} onclick={handleClick}>
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
			<ClientCode code={client.code} class="opacity-0 group-hover:opacity-100 transition-opacity"/>
		</Card.Description>
	</Card.Header>
</Card.Root>
