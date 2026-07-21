<script lang="ts">
	import ClientCreateDialog from "$lib/features/client/components/client-create-dialog.svelte";
	import ClientCardGrid from "$lib/features/client/components/client-card-grid.svelte";
	import PageLayout from "$lib/components/page-layout.svelte";
	import { clientService } from "$lib/features/client/service";
	import PlusIcon from "@lucide/svelte/icons/plus";

	const listQuery = clientService.list();
	const clients = $derived(listQuery.data ?? []);
</script>

<PageLayout title="Clients" description="Manage your clients.">
	{#snippet actionsRight()}
		<ClientCreateDialog>
			{#snippet trigger()}
				<PlusIcon />
				Create client
			{/snippet}
		</ClientCreateDialog>
	{/snippet}
	<ClientCardGrid {clients} isLoading={listQuery.isLoading} />
</PageLayout>
