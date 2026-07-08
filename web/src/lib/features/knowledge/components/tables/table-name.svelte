<script lang="ts">
	import TableField from "$lib/components/flex-table/table-field.svelte";
	import type { ComponentProps } from "svelte";
	import type { Knowledge } from "../../schemas";
	import CircleXIcon from "@lucide/svelte/icons/circle-x";
	import { Tooltip, TooltipContent, TooltipTrigger } from "$lib/components/ui/tooltip";

  type Props = Omit<ComponentProps<typeof TableField>, "value"> & {
    knowledge: Knowledge;
  };
  let { knowledge, truncate }: Props = $props();
</script>
<div class="flex items-center gap-2">
  <TableField value={knowledge.name} {truncate} />
  {#if knowledge.sync_error}
    <Tooltip>
      <TooltipTrigger>
        <CircleXIcon class="size-4 text-destructive" />
      </TooltipTrigger>
      <TooltipContent>
        {knowledge.sync_error}
      </TooltipContent>
    </Tooltip>
  {/if}
</div>