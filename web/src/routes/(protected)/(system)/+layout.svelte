<script lang="ts">
	import { SidebarInset, SidebarProvider } from "$lib/components/ui/sidebar";
	import AppSidebar from "$lib/components/app-sidebar.svelte";
	import { systemSidebarSections } from "$lib/navigation/system-sidebar";
	let { children } = $props();
</script>

<SidebarProvider>
	<AppSidebar sections={systemSidebarSections} />
	<SidebarInset class="min-h-0 overflow-hidden">
		{@render children()}
	</SidebarInset>
</SidebarProvider>