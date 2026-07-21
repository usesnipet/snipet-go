<script lang="ts">
	import PageLayout from "$lib/components/page-layout.svelte";
	import SessionCreateDialog from "$lib/features/session/components/session-create-dialog.svelte";
	import SessionInfo from "$lib/features/session/components/session-info.svelte";
	import SessionTable from "$lib/features/session/components/session-table.svelte";
	import { sessionService } from "$lib/features/session/service";
	import type { Session } from "$lib/features/session/schemas";
	import PlusIcon from "@lucide/svelte/icons/plus";

	let { params } = $props();

	const listQuery = $derived(sessionService.list(params.code, { include: [ "agent" ] }));
	const sessions = $derived(listQuery.data ?? []);

	let selectedId = $state<string | null>(null);
	const selectedSession = $derived(
		selectedId ? (sessions.find((session) => session.id === selectedId) ?? null) : null,
	);

	function handleSelect(session: Session) {
		selectedId = session.id;
	}
</script>

<PageLayout title="Sessions" description="Client session history.">
	{#snippet actionsRight()}
		<SessionCreateDialog clientCode={params.code}>
			{#snippet trigger()}
				<PlusIcon />
				Create session
			{/snippet}
		</SessionCreateDialog>
	{/snippet}

	<div class="flex h-full min-h-0 w-full gap-4">
		<div class="h-full min-w-0 flex-1 basis-0">
			<SessionTable
				{sessions}
				clientCode={params.code}
				isLoading={listQuery.isLoading}
				onRowClick={handleSelect}
			/>
		</div>
		{#if selectedSession}
			<div class="h-full w-[30%] shrink-0 basis-[30%]">
				<SessionInfo session={selectedSession} onClose={() => (selectedId = null)} />
			</div>
		{/if}
	</div>
</PageLayout>
