<script lang="ts">
	import * as Sidebar from "$lib/components/ui/sidebar";
	import AppSidebar from "$lib/components/app-sidebar.svelte";
	import { createClientSidebarSections } from "$lib/navigation/client-sidebar";
	import { ArrowLeftIcon } from "@lucide/svelte";
	import ThemeToggle from "$lib/components/theme-toggle.svelte";
	import { resolve } from "$app/paths";
	import { clientService } from "$lib/features/client/service.js";
	import ClientCode from "$lib/features/client/components/client-code.svelte";

	let { children, params } = $props();

	const sections = $derived(params.code ? createClientSidebarSections(params.code) : []);

	const clientQuery = $derived.by(() => {
		if (!params.code) {
			return null;
		}
		const query = clientService.findByCode(params.code)
		return query;
	});
	const client = $derived(clientQuery?.data);
</script>

<Sidebar.Provider>
	<AppSidebar {sections}>
		{#snippet header()}
			<Sidebar.Menu>
				<Sidebar.MenuItem>
					<Sidebar.MenuButton size="lg">
						{#snippet child()}
							<div class="flex items-center gap-2">
								<div class="flex flex-1 items-center gap-2">
									<a href={resolve("/")}>
										<div
											class="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg"
										>
											<ArrowLeftIcon class="size-4" />
										</div>
									</a>
									<div class="flex-1 text-start text-sm flex flex-col leading-tight">
										<span class="truncate font-medium">{client?.name}</span>
										<div>

											<ClientCode code={client?.code} />
										</div>
									</div>
								</div>
								<ThemeToggle />
							</div>
						{/snippet}
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
			</Sidebar.Menu>

		{/snippet}
	</AppSidebar>
	<Sidebar.Inset class="min-h-0 overflow-hidden">
		{@render children()}
	</Sidebar.Inset>
</Sidebar.Provider>