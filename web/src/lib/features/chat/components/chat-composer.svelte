<script lang="ts">
	import { goto } from "$app/navigation";
	import { Button } from "$lib/components/ui/button";
	import { Textarea } from "$lib/components/ui/textarea";
	import ArrowUpIcon from "@lucide/svelte/icons/arrow-up";
	import { mockSessions } from "../stores/mock-sessions.svelte";
	import { resolve } from "$app/paths";

	type Props = {
		title?: string;
	};

	let { title = "Snipet" }: Props = $props();

	let message = $state("");

	function sessionNameFromPrompt(prompt: string): string {
		const trimmed = prompt.trim().replace(/\s+/g, " ");
		if (trimmed.length <= 40) return trimmed;
		return `${trimmed.slice(0, 40).trimEnd()}…`;
	}

	function handleSubmit() {
		const prompt = message.trim();
		if (!prompt) return;
		const session = mockSessions.create(sessionNameFromPrompt(prompt));
		message = "";
		goto(resolve("/(chat)/session/[session_id]", { session_id: session.id }));
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === "Enter" && !event.shiftKey) {
			event.preventDefault();
			handleSubmit();
		}
	}
</script>

<div class="flex h-full flex-col items-center justify-center gap-8 px-4">
	<div class="max-w-xl text-center">
		<h1 class="text-3xl font-semibold tracking-tight">{title}</h1>
		<p class="text-muted-foreground mt-2 text-sm">
			Start a conversation to try agents, tools, and knowledge.
		</p>
	</div>

	<form
		class="bg-background relative w-full max-w-2xl rounded-2xl border shadow-sm"
		onsubmit={(e) => {
			e.preventDefault();
			handleSubmit();
		}}
	>
		<Textarea
			class="min-h-28 resize-none border-0 bg-transparent p-4 pr-14 shadow-none focus-visible:ring-0"
			placeholder="Message…"
			bind:value={message}
			onkeydown={handleKeydown}
		/>
		<div class="absolute right-3 bottom-3">
			<Button
				type="submit"
				size="icon"
				disabled={!message.trim()}
				aria-label="Send message"
			>
				<ArrowUpIcon />
			</Button>
		</div>
	</form>
</div>
