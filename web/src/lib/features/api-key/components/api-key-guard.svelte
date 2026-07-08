<script lang="ts">
  import { goto } from "$app/navigation";
  import { resolve } from "$app/paths";
	import { apiKeyAuth } from "../stores/api-key-auth.svelte";

  let { children } = $props();

  $effect(() => {
    if (!apiKeyAuth.loading && !apiKeyAuth.isAuthenticated) {
      goto(resolve("/auth/api-key"));
    }
  });
</script>

{#if apiKeyAuth.isAuthenticated}
  {@render children()}
{/if}