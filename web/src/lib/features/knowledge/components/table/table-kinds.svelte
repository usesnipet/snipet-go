<script lang="ts">
	import { Badge } from "$lib/components/ui/badge";
	import type { KnowledgeItemKind } from "../../schemas";

  type Props = { kinds: KnowledgeItemKind[] };

  let props: Props = $props();
  const kinds = $derived(props.kinds.map(kind => kind.toUpperCase()));
</script>

<div class="flex items-center gap-2">
  {#each kinds as kind (kind)}
    <Badge variant="outline">{kind}</Badge>
  {/each}
</div>