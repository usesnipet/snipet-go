<script lang="ts">
	import TableField from "$lib/components/flex-table/table-field.svelte";
	import type { ComponentProps } from "svelte";
	import CopyIcon from "@lucide/svelte/icons/copy";

  type Props = ComponentProps<typeof TableField>;

  let { value, truncate }: Props = $props();
</script>
<div class="flex items-center gap-2">
  <TableField {value} {truncate} />
  <CopyIcon />
</div>