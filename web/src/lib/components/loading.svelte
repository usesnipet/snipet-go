<script lang="ts" module>
	import { type VariantProps, tv } from "tailwind-variants";

	export const loadingVariants = tv({
		slots: {
			root: "flex h-full flex-col items-center justify-center",
			logo: "animate-spin",
			message: "text-muted-foreground text-sm"
		},
		variants: {
			size: {
				sm: {
					root: "gap-2",
					logo: "size-10",
          message: "text-sm",
				},
				md: {
					root: "gap-2.5",
					logo: "size-15",
          message: "text-base",
				},
				lg: {
					root: "gap-3",
					logo: "size-20",
          message: "text-lg",
				}
			}
		},
		defaultVariants: {
			size: "lg"
		}
	});

	export type LoadingSize = VariantProps<typeof loadingVariants>["size"];
</script>

<script lang="ts">
	import Logo from "./logo.svelte";
	import { cn } from "$lib/utils";

	type Props = {
		class?: string;
		message?: string;
		size?: LoadingSize;
	};

	let { class: className, message = "Loading…", size = "lg" }: Props = $props();

	const styles = $derived(loadingVariants({ size }));
</script>

<div class={cn(styles.root(), className)} role="status" aria-live="polite">
	<Logo {size} class={styles.logo()} />
	<p class={styles.message()}>{message}</p>
</div>
