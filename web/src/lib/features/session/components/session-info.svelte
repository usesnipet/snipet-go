<script lang="ts">
	import * as Card from "$lib/components/ui/card";
	import { Skeleton } from "$lib/components/ui/skeleton";
	import { JsonConfigCard } from "$lib/components/json-config";
	import Button from "$lib/components/ui/button/button.svelte";
	import { XIcon } from "@lucide/svelte";
	import { agentService } from "$lib/features/agent/service";
	import type { Session } from "../schemas";

	type Props = {
		session?: Session;
		isLoading?: boolean;
		onClose?: () => void;
	};

	let { session, isLoading = false, onClose }: Props = $props();

	const agentQuery = $derived(
		session?.agent_id ? agentService.findById(session.agent_id) : undefined,
	);
	const agent = $derived(agentQuery?.data);
	const title = $derived(session?.metadata.name ? session.metadata.name : "Session");
</script>

<Card.Root class="h-full w-full min-w-0">
	<Card.Header>
		<div class="flex items-center justify-between gap-2">
			<Card.Title>
				{#if isLoading}
					<Skeleton class="h-6 w-3/4" />
				{:else}
					{title}
				{/if}
			</Card.Title>
			{#if onClose}
				<Button variant="ghost" size="icon" onclick={onClose} aria-label="Close session details">
					<XIcon />
				</Button>
			{/if}
		</div>
	</Card.Header>

	<Card.Content class="space-y-4">
		<div class="space-y-1">
			<p class="text-muted-foreground text-sm">ID</p>
			{#if isLoading}
				<Skeleton class="h-4 w-2/3" />
			{:else}
				<p class="text-sm font-medium font-mono break-all">{session?.id ?? "—"}</p>
			{/if}
		</div>

		<div class="space-y-1">
			<p class="text-muted-foreground text-sm">Client ID</p>
			{#if isLoading}
				<Skeleton class="h-4 w-2/3" />
			{:else}
				<p class="text-sm font-medium font-mono break-all">{session?.client_id ?? "—"}</p>
			{/if}
		</div>

		<div class="space-y-1">
			<p class="text-muted-foreground text-sm">Agent</p>
			{#if isLoading}
				<Skeleton class="h-4 w-1/2" />
			{:else if agent?.name}
				<p class="text-sm font-medium">{agent.name}</p>
				<p class="text-muted-foreground text-xs font-mono break-all">{session?.agent_id}</p>
			{:else}
				<p class="text-sm font-medium font-mono break-all">{session?.agent_id ?? "—"}</p>
			{/if}
		</div>

		<div class="flex flex-col gap-1">
			<p class="text-muted-foreground text-sm">Metadata</p>
			{#if isLoading}
				<Skeleton class="h-16 w-full" />
			{:else}
				<JsonConfigCard data={session?.metadata} emptyMessage="No metadata" />
			{/if}
		</div>
	</Card.Content>
</Card.Root>
