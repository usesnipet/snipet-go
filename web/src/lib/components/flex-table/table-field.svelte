<script lang="ts">
	import { truncate } from "$lib/utils";
	import * as Tooltip from "../ui/tooltip";

  type Props = {
    value: string | null | undefined;
    truncate?: number;
  }

  let props: Props = $props();

  const truncated = $derived(props.truncate && props.value && props.value?.length > props.truncate);
  const truncatedValue = $derived(props.truncate ? truncate(props.value, props.truncate) : props.value);
</script>

{#if truncatedValue}
  {#if truncated}
    <Tooltip.Root>
      <Tooltip.Trigger>
        {truncatedValue}
      </Tooltip.Trigger>
      <Tooltip.Content>
        {props.value}
      </Tooltip.Content>
    </Tooltip.Root>
  {:else}
	  {truncatedValue}
  {/if}
{:else}
  -
{/if}