<script lang="ts">
	import { Button } from "$lib/components/ui/button";
	import { resolve } from "$app/paths";
</script>

<div class="flex h-svh flex-col items-center justify-center gap-6 px-6 text-center">
	<div class="max-w-md space-y-3">
		<h1 class="text-2xl font-semibold tracking-tight">Inherit client is disabled</h1>
		<p class="text-muted-foreground text-sm leading-relaxed">
			The built-in chat client is turned off. Set the environment variable
			<code class="bg-muted rounded px-1.5 py-0.5 text-xs">APP_INHERIT_CLIENT=true</code>
			and restart the server to enable it.
		</p>
		<p class="text-muted-foreground text-sm leading-relaxed">
			You can still manage agents, knowledge, and clients from the management area.
		</p>
	</div>

	<Button href={resolve("/auth/api-key")}>Go to API key login</Button>
</div>
