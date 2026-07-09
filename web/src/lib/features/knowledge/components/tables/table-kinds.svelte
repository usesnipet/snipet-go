<script lang="ts">
	import { Badge } from "$lib/components/ui/badge";
	import * as Tooltip from "$lib/components/ui/tooltip/index.js";
	import type { KnowledgeItemKind } from "../../schemas";

  type Props = { kinds: KnowledgeItemKind[] };

  let props: Props = $props();
  const kinds = $derived(props.kinds.map(kind => kind.toUpperCase()));
</script>

<div class="flex items-center gap-2">
  {#if kinds.length === 1}
  <Badge variant="outline">{kinds[0]}</Badge>
  {:else}
  <Tooltip.Root>
    <Tooltip.Trigger>
      <Badge variant="outline">{kinds[0]}</Badge>
      <Badge variant="outline">+ {kinds.length - 1} more</Badge>
    </Tooltip.Trigger>
    <Tooltip.Content>
      {kinds.join(", ")}
    </Tooltip.Content>
  </Tooltip.Root>
  {/if}
</div>