<script lang="ts">
  import { Button } from "$lib/components/ui/button/index.js";
  import { CopyIcon } from "@lucide/svelte";
	import { toast } from "svelte-sonner";

  let { code, class: className }: { code: string | null | undefined, class?: string } = $props();

  function handleCopyCode(e: MouseEvent) {
    if (!code) return;
    e.preventDefault();
    e.stopPropagation();
    navigator.clipboard.writeText(code);
    toast.success("Client code copied.");
  }
</script>
{#if code}
  {code}
  <Button
    variant="outline"
    size="icon-xs"
    class={className}
    onclick={handleCopyCode}
    aria-label="Copy client code"
  >
    <CopyIcon class="size-2.5" />
  </Button>
{/if}
