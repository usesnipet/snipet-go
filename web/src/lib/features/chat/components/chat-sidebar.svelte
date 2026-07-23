<script lang="ts">
	import { goto } from "$app/navigation";
	import { Button } from "$lib/components/ui/button";
	import PlusIcon from "@lucide/svelte/icons/plus";
	import SessionList from "./session-list.svelte";
	import { resolve } from "$app/paths";
	import Logo from "$lib/components/logo.svelte";

	type Props = {
		clientName: string;
	};

	let { clientName }: Props = $props();

	function handleNewChat() {
		goto(resolve("/(chat)"));
	}
</script>

<aside class="bg-background flex h-full w-65 shrink-0 flex-col border-r">
	<div class="flex flex-col gap-3 p-3">
		<div class="flex items-center gap-2 px-1">
		  <Logo />
			<div>
				<p class="truncate text-sm font-semibold">{clientName}</p>
				<p class="text-muted-foreground text-xs">Chat</p>
			</div>
		</div>
		<Button class="w-full justify-start" variant="outline" onclick={handleNewChat}>
			<PlusIcon />
			New chat
		</Button>
	</div>

	<SessionList />
</aside>
