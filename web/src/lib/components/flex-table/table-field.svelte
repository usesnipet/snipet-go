<script lang="ts">
	import { truncate } from "$lib/utils";

  type Props = {
    value: string;
    truncate?: number;
  }

  let props: Props = $props();

  const truncatedValue = $derived(props.truncate ? truncate(props.value, props.truncate) : props.value);
</script>

{truncatedValue}