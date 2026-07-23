<script lang="ts" module>
	import { type VariantProps, tv } from "tailwind-variants";

	export const logoVariants = tv({
		base: "shrink-0",
		variants: {
			size: {
				sm: "size-4",
				md: "size-6",
				lg: "size-8"
			}
		},
		defaultVariants: {
			size: "md"
		}
	});

	export type LogoSize = VariantProps<typeof logoVariants>["size"];
</script>

<script lang="ts">
	import logo from "$lib/assets/logo.svg";
	import { cn } from "$lib/utils";

	type Props = {
		class?: string;
		size?: LogoSize;
		alt?: string;
	};

	let { class: className, size = "md", alt = "Snipet logo" }: Props = $props();
</script>

<img {alt} src={logo} class={cn(logoVariants({ size }), className)} />
