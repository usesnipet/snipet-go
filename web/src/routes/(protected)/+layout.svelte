<script lang="ts">
	import { SidebarInset, SidebarProvider } from "$lib/components/ui/sidebar";
	import AppSidebar from "$lib/components/app-sidebar.svelte";
	import ApiKeyGuard from "$lib/features/api-key/components/api-key-guard.svelte";

  let { children } = $props();
</script>

<ApiKeyGuard>
  <SidebarProvider>
    <AppSidebar />
    <SidebarInset class="min-h-0 overflow-hidden">
      {@render children()}
    </SidebarInset>
  </SidebarProvider>
</ApiKeyGuard>