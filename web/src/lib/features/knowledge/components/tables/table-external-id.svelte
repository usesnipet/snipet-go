<script lang="ts">
	import TableField from "$lib/components/flex-table/table-field.svelte";
	import { buttonVariants } from "$lib/components/ui/button";
	import type { ComponentProps } from "svelte";
	import type { Attachment } from "svelte/attachments";
	import CopyIcon from "@lucide/svelte/icons/copy";
	import { on } from "svelte/events";
	import { toast } from "svelte-sonner";

	type Props = ComponentProps<typeof TableField> & {
		attach?: Attachment;
	};

	let { value, truncate, attach }: Props = $props();

	async function copyExternalId() {
		if (!value) return;

		try {
			if (navigator.clipboard?.writeText) {
				await navigator.clipboard.writeText(value);
			} else {
				const textarea = document.createElement("textarea");
				textarea.value = value;
				textarea.setAttribute("readonly", "");
				textarea.style.position = "absolute";
				textarea.style.left = "-9999px";
				document.body.appendChild(textarea);
				textarea.select();
				document.execCommand("copy");
				document.body.removeChild(textarea);
			}

			toast.success("External ID copied to clipboard.");
		} catch {
			toast.error("Failed to copy external ID.");
		}
	}

	const copyButtonAttachment: Attachment<HTMLButtonElement> = (element) => {
		return on(element, "click", (event) => {
			event.stopPropagation();
			void copyExternalId();
		});
	};
</script>

<div class="flex items-center gap-2" {@attach attach}>
	{#if value}
		<button
			type="button"
			class={buttonVariants({ variant: "outline", size: "icon", class: "relative z-10 shrink-0" })}
			aria-label="Copy external ID"
			{@attach copyButtonAttachment}
		>
			<CopyIcon class="size-4" />
		</button>
		<span class="min-w-0 flex-1 truncate">
			<TableField {value} {truncate} />
		</span>
	{:else}
		<TableField value="N/A" />
	{/if}
</div>
