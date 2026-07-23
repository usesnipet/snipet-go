<script lang="ts">
	import AgentCreateDialog from "$lib/features/agent/components/agent-create-dialog.svelte";
	import AgentCardGrid from "$lib/features/agent/components/agent-card-grid.svelte";
	import PageLayout from "$lib/components/page-layout.svelte";
	import { agentService } from "$lib/features/agent/service";
	import PlusIcon from "@lucide/svelte/icons/plus";

	const listQuery = agentService.list();
	const agents = $derived(listQuery.data ?? []);
</script>

<PageLayout title="Agents" description="Manage your agents.">
	{#snippet actionsRight()}
		<AgentCreateDialog>
			{#snippet trigger()}
				<PlusIcon />
				Create agent
			{/snippet}
		</AgentCreateDialog>
	{/snippet}
	<AgentCardGrid {agents} isLoading={listQuery.isLoading} />
</PageLayout>
