<script lang="ts">
	import { Button, buttonVariants } from "$lib/components/ui/button";
	import * as Dialog from "$lib/components/ui/dialog";
	import { clientService } from "../service";
	import { defaults, superForm } from "sveltekit-superforms";
	import { zod4, zod4Client } from "sveltekit-superforms/adapters";
	import { updateClientSchema, type Client } from "../schemas";
	import { FieldGroup } from "$lib/components/ui/field";
	import FormInput from "$lib/components/form/form-input.svelte";
	import FormError from "$lib/components/form/form-error.svelte";
	import ClientConfigFields from "./client-config-fields.svelte";
	import type { Snippet } from "svelte";

	type Props = {
		open: boolean;
		client?: Client;
		trigger?: Snippet;
	};

	let { open = $bindable(false), client, trigger }: Props = $props();

	const updateClientMutation = clientService.update();
	const form = superForm(defaults(zod4(updateClientSchema)), {
		SPA: true,
		dataType: "json",
		validators: zod4Client(updateClientSchema),
		async onUpdate({ form: formState }) {
			if (!formState.valid) return;
			if (!client) return;
			updateClientMutation.mutate(
				{
					data: formState.data,
					code: client.code,
				},
				{
					onSuccess: () => {
						open = false;
					},
				},
			);
		},
	});

	$effect(() => {
		if (client) {
			form.form.set({
				name: client.name,
				config: {
					oidc: {
						enabled: client.config.oidc.enabled,
						issuer: client.config.oidc.issuer,
						audience: client.config.oidc.audience,
					},
					webhook: {
						enabled: client.config.webhook.enabled,
						url: client.config.webhook.url,
					},
					anonymous: {
						enabled: client.config.anonymous.enabled,
					},
				},
			});
		}
	});
	const { enhance } = form;

	function handleCancel() {
		form.reset();
		open = false;
	}
</script>

<Dialog.Root bind:open>
	{#if trigger}
		<Dialog.Trigger class={buttonVariants({ variant: "outline" })}>
			{@render trigger()}
		</Dialog.Trigger>
	{/if}

	<Dialog.Content class="sm:max-w-lg">
		<Dialog.Header>
			<Dialog.Title>Update client</Dialog.Title>
			<Dialog.Description>Fill the fields below to update the client.</Dialog.Description>
		</Dialog.Header>

		<form use:enhance class="flex max-h-[70vh] flex-col gap-6">
			<FieldGroup>
				<FormInput {form} field="name" label="Name" placeholder="e.g. Acme Corp" />
				<ClientConfigFields {form} />
				<FormError {form} />
				<Dialog.Footer>
					<Button variant="outline" type="button" onclick={handleCancel}>Cancel</Button>
					<Button disabled={updateClientMutation.isPending} type="submit">Update</Button>
				</Dialog.Footer>
			</FieldGroup>
		</form>
	</Dialog.Content>
</Dialog.Root>
