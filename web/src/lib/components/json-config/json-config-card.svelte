<script lang="ts">
	import { cn } from "$lib/utils";
	import JsonConfigNode from "./json-config-node.svelte";
	import { isPlainObject, type JsonValue } from "./utils";

	type Props = {
		data: JsonValue | null | undefined;
		class?: string;
		emptyMessage?: string;
	};

	let { data, class: className, emptyMessage = "No configuration" }: Props = $props();

	const hasContent = $derived(
		data !== null &&
			data !== undefined &&
			!(isPlainObject(data) && Object.keys(data).length === 0) &&
			!(Array.isArray(data) && data.length === 0),
	);
</script>

<div
	class={cn(
		"bg-muted/30 border-border/60 min-w-0 rounded-lg border px-3 py-2.5",
		className,
	)}
>
	{#if hasContent}
		<JsonConfigNode value={data as JsonValue} />
	{:else}
		<p class="text-muted-foreground text-sm">{emptyMessage}</p>
	{/if}
</div>
