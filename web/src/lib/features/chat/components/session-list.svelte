<script lang="ts">
	import { page } from "$app/state";
	import { Input } from "$lib/components/ui/input";
	import SearchIcon from "@lucide/svelte/icons/search";
	import { mockSessions } from "../stores/mock-sessions.svelte";
	import SessionListItem from "./session-list-item.svelte";

	let search = $state("");

	const sessions = $derived(mockSessions.filtered(search));
	const activeId = $derived(page.params.session_id);
</script>

<div class="flex min-h-0 flex-1 flex-col gap-3 px-3 pb-3">
	<div class="relative">
		<SearchIcon
			class="text-muted-foreground pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2"
		/>
		<Input class="pl-8" placeholder="Search sessions..." bind:value={search} />
	</div>

	<div class="min-h-0 flex-1 space-y-0.5 overflow-y-auto">
		{#each sessions as session (session.id)}
			<SessionListItem {session} active={session.id === activeId} />
		{:else}
			<p class="text-muted-foreground px-2 py-4 text-center text-sm">No sessions found.</p>
		{/each}
	</div>
</div>
