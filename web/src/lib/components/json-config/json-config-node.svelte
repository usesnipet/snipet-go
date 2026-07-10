<script lang="ts">
	import ChevronRightIcon from "@lucide/svelte/icons/chevron-right";
	import { Badge } from "$lib/components/ui/badge";
	import * as Collapsible from "$lib/components/ui/collapsible";
	import * as Tooltip from "$lib/components/ui/tooltip";
	import { cn } from "$lib/utils";
	import JsonConfigNode from "./json-config-node.svelte";
	import {
		formatKey,
		isComplexArray,
		isPlainObject,
		shouldDefaultOpen,
		summarizeValue,
		type JsonValue,
	} from "./utils";

	type Props = {
		label?: string;
		value: JsonValue;
		depth?: number;
	};

	let { label, value, depth = 0 }: Props = $props();

	const stringTruncateLimit = 80;
	const defaultOpen = $derived(shouldDefaultOpen(value, depth));
	let objectOpen = $state(false);
	let arrayOpen = $state(false);
	let openInitialized = $state(false);

	$effect.pre(() => {
		if (!openInitialized) {
			objectOpen = defaultOpen;
			arrayOpen = defaultOpen;
			openInitialized = true;
		}
	});
</script>

{#if isPlainObject(value)}
	{@const entries = Object.entries(value)}
	{#if entries.length === 0}
		<div class="flex min-w-0 flex-col gap-0.5">
			{#if label}
				<span class="text-muted-foreground text-xs">{formatKey(label)}</span>
			{/if}
			<span class="text-muted-foreground text-sm">—</span>
		</div>
	{:else if entries.length > 6 || depth > 0}
		<Collapsible.Root bind:open={objectOpen} class="group/collapsible min-w-0">
			<div class="flex min-w-0 flex-col gap-0.5">
				<Collapsible.Trigger
					class="hover:bg-muted/60 -ms-1 flex w-fit max-w-full items-center gap-1 rounded-md px-1 py-0.5 text-start transition-colors"
				>
					<ChevronRightIcon
						class="text-muted-foreground size-3.5 shrink-0 transition-transform group-data-[state=open]/collapsible:rotate-90"
					/>
					{#if label}
						<span class="text-muted-foreground text-xs">{formatKey(label)}</span>
					{/if}
					<Badge variant="outline" class="max-w-full truncate font-normal">
						{summarizeValue(value)}
					</Badge>
				</Collapsible.Trigger>
				<Collapsible.Content>
					<div
						class={cn(
							"flex flex-col gap-2 pt-1",
							depth > 0 && "border-border/60 ms-1.5 border-s ps-2.5",
						)}
					>
						{#each entries as [key, childValue] (key)}
							<JsonConfigNode label={key} value={childValue} depth={depth + 1} />
						{/each}
					</div>
				</Collapsible.Content>
			</div>
		</Collapsible.Root>
	{:else}
		<div class="flex min-w-0 flex-col gap-2">
			{#if label}
				<span class="text-muted-foreground text-xs">{formatKey(label)}</span>
			{/if}
			<div
				class={cn(
					"flex flex-col gap-2",
					depth > 0 && "border-border/60 ms-1.5 border-s ps-2.5",
				)}
			>
				{#each entries as [key, childValue] (key)}
					<JsonConfigNode label={key} value={childValue} depth={depth + 1} />
				{/each}
			</div>
		</div>
	{/if}
{:else if Array.isArray(value)}
	{#if value.length === 0}
		<div class="flex min-w-0 flex-col gap-0.5">
			{#if label}
				<span class="text-muted-foreground text-xs">{formatKey(label)}</span>
			{/if}
			<span class="text-muted-foreground text-sm">—</span>
		</div>
	{:else if isComplexArray(value)}
		<Collapsible.Root bind:open={arrayOpen} class="group/collapsible min-w-0">
			<div class="flex min-w-0 flex-col gap-0.5">
				<Collapsible.Trigger
					class="hover:bg-muted/60 -ms-1 flex w-fit max-w-full items-center gap-1 rounded-md px-1 py-0.5 text-start transition-colors"
				>
					<ChevronRightIcon
						class="text-muted-foreground size-3.5 shrink-0 transition-transform group-data-[state=open]/collapsible:rotate-90"
					/>
					{#if label}
						<span class="text-muted-foreground text-xs">{formatKey(label)}</span>
					{/if}
					<Badge variant="outline" class="font-normal">
						{summarizeValue(value)}
					</Badge>
				</Collapsible.Trigger>
				<Collapsible.Content>
					<div
						class={cn(
							"flex flex-col gap-2 pt-1",
							depth > 0 && "border-border/60 ms-1.5 border-s ps-2.5",
						)}
					>
						{#each value as item, index (index)}
							<JsonConfigNode label={`Item ${index + 1}`} value={item} depth={depth + 1} />
						{/each}
					</div>
				</Collapsible.Content>
			</div>
		</Collapsible.Root>
	{:else}
		<div class="flex min-w-0 flex-col gap-1">
			{#if label}
				<span class="text-muted-foreground text-xs">{formatKey(label)}</span>
			{/if}
			<div class="flex flex-wrap gap-1">
				{#each value as item, index (index)}
					{#if typeof item === "boolean"}
						<Badge variant={item ? "secondary" : "outline"}>{item ? "Yes" : "No"}</Badge>
					{:else}
						<Badge variant="outline" class="max-w-full truncate font-normal">
							{String(item)}
						</Badge>
					{/if}
				{/each}
			</div>
		</div>
	{/if}
{:else if value === null}
	<div class="flex min-w-0 flex-col gap-0.5">
		{#if label}
			<span class="text-muted-foreground text-xs">{formatKey(label)}</span>
		{/if}
		<span class="text-muted-foreground text-sm">—</span>
	</div>
{:else if typeof value === "boolean"}
	<div class="flex min-w-0 flex-col gap-0.5">
		{#if label}
			<span class="text-muted-foreground text-xs">{formatKey(label)}</span>
		{/if}
		<Badge variant={value ? "secondary" : "outline"} class="w-fit">
			{value ? "Yes" : "No"}
		</Badge>
	</div>
{:else if typeof value === "number"}
	<div class="flex min-w-0 flex-col gap-0.5">
		{#if label}
			<span class="text-muted-foreground text-xs">{formatKey(label)}</span>
		{/if}
		<span class="font-mono text-sm font-medium">{value}</span>
	</div>
{:else}
	{@const truncated = value.length > stringTruncateLimit}
	<div class="flex min-w-0 flex-col gap-0.5">
		{#if label}
			<span class="text-muted-foreground text-xs">{formatKey(label)}</span>
		{/if}
		{#if truncated}
			<Tooltip.Root>
				<Tooltip.Trigger class="w-fit max-w-full text-start">
					<span class="text-sm font-medium break-all">
						{value.slice(0, stringTruncateLimit)}…
					</span>
				</Tooltip.Trigger>
				<Tooltip.Content class="max-w-sm break-all">
					{value}
				</Tooltip.Content>
			</Tooltip.Root>
		{:else}
			<span class="text-sm font-medium break-all">{value}</span>
		{/if}
	</div>
{/if}
