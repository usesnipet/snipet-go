<script lang="ts">
	import { truncate } from "$lib/utils";

  type Props = {
    value: string | null | undefined;
    truncate?: number;
  }

  let props: Props = $props();

  const truncatedValue = $derived(props.truncate ? truncate(props.value, props.truncate) : props.value);
</script>

{#if truncatedValue}
	{truncatedValue}
{:else}
  -
{/if}