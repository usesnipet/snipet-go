<script lang="ts">
  import { auth } from "$lib/store/auth.svelte";
  import { goto } from "$app/navigation";
  import { resolve } from "$app/paths";
	import { SidebarInset, SidebarProvider } from "$lib/components/ui/sidebar";
	import AppSidebar from "$lib/components/app-sidebar.svelte";

  let { children } = $props();

  $effect(() => {
    if (!auth.loading && !auth.isAuthenticated) {
      goto(resolve("/auth/api-key"));
    }
  });
</script>

{#if auth.isAuthenticated}
  <SidebarProvider>
    <AppSidebar />
    <SidebarInset>
      {@render children()}
    </SidebarInset>
  </SidebarProvider>
{/if}