<script lang="ts">
	import { truncate } from "$lib/utils";
	import { Badge, type BadgeVariant } from "../ui/badge";

  type Props = {
    value: string | null | undefined;
    truncate?: number;
    variant?: BadgeVariant;
  }

  let props: Props = $props();

  const truncatedValue = $derived(props.truncate ? truncate(props.value, props.truncate) : props.value);
</script>

<Badge variant={props.variant}>{truncatedValue}</Badge>