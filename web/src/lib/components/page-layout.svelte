<script lang="ts">
	import { Separator } from "$lib/components/ui/separator";
	import { cn } from "$lib/utils";
	import type { Snippet } from "svelte";
	import type { HTMLAttributes } from "svelte/elements";
  import { toast } from "svelte-sonner";

	let {
		title,
		description,
		actionsLeft,
		actionsRight,
		class: className,
		children,
		...restProps
	}: HTMLAttributes<HTMLDivElement> & {
		title: string;
		description?: string;
		actionsLeft?: Snippet;
		actionsRight?: Snippet;
		children: Snippet;
	} = $props();
</script>
<svelte:head>
	<title>{title}</title>
	{#if description}
		<meta name="description" content={description} />
	{/if}
</svelte:head>
<div class={cn("flex h-full min-h-0 flex-1 flex-col gap-6 p-4 md:p-6", className)} {...restProps}>
	<div class="flex shrink-0 flex-col gap-4">
		<div class="flex items-start justify-between gap-4">
			<div class="flex min-w-0 items-start gap-3">
				{#if actionsLeft}
					<div class="flex shrink-0 items-center gap-2">
						{@render actionsLeft()}
					</div>
				{/if}

				<div class="flex min-w-0 flex-col gap-1">
					<h1 class="text-2xl font-semibold tracking-tight">{title}</h1>
					{#if description}
						<p class="text-muted-foreground text-sm">{description}</p>
					{/if}
				</div>
			</div>

			{#if actionsRight}
				<div class="flex shrink-0 items-center gap-2">
					{@render actionsRight()}
				</div>
			{/if}
		</div>

		<Separator />
	</div>

	<div class="flex min-h-0 flex-1 flex-col">
		<svelte:boundary onerror={(e) => toast.error((e as Error).message)}>
			{@render children()}
		</svelte:boundary>
	</div>
</div>
