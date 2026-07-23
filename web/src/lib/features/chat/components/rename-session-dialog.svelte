<script lang="ts">
	import { untrack } from "svelte";
	import * as Dialog from "$lib/components/ui/dialog";
	import { Button } from "$lib/components/ui/button";
	import { Input } from "$lib/components/ui/input";
	import { Label } from "$lib/components/ui/label";
	import { mockSessions } from "../stores/mock-sessions.svelte";

	type Props = {
		open?: boolean;
		sessionId: string;
		initialName?: string;
		onRename?: (id: string, name: string) => void;
	};

	let {
		open = $bindable(false),
		sessionId,
		initialName = "",
		onRename,
	}: Props = $props();

	let name = $state(untrack(() => initialName));

	function handleCancel() {
		open = false;
	}

	function handleSave() {
		const trimmed = name.trim();
		if (!trimmed) return;
		mockSessions.rename(sessionId, trimmed);
		onRename?.(sessionId, trimmed);
		open = false;
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>Rename session</Dialog.Title>
			<Dialog.Description>Choose a new name for this chat session.</Dialog.Description>
		</Dialog.Header>

		<div class="grid gap-2 py-2">
			<Label for="session-name">Name</Label>
			<Input
				id="session-name"
				bind:value={name}
				placeholder="Session name"
				onkeydown={(e) => {
					if (e.key === "Enter") {
						e.preventDefault();
						handleSave();
					}
				}}
			/>
		</div>

		<Dialog.Footer>
			<Button variant="outline" onclick={handleCancel}>Cancel</Button>
			<Button onclick={handleSave} disabled={!name.trim()}>Save</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
