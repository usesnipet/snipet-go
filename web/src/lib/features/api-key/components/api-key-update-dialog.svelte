<script lang="ts">
	import { Button, buttonVariants } from "$lib/components/ui/button";
	import * as Dialog from "$lib/components/ui/dialog";
	import { FieldGroup } from "$lib/components/ui/field";
	import FormError from "$lib/components/form/form-error.svelte";
	import FormExpiration from "$lib/components/form/form-expiration.svelte";
	import { DEFAULT_EXPIRATION_OPTIONS, fromExpirationFormValue, toExpirationFormValue } from "$lib/components/form/expiration";
	import { apiKeyService } from "../service";
	import {
		updateApiKeyExpirationFormSchema,
		type APIKey,
	} from "../schemas";
	import { defaults, superForm } from "sveltekit-superforms";
	import { zod4, zod4Client } from "sveltekit-superforms/adapters";
	import type { Snippet } from "svelte";

	type Props = {
		open: boolean;
		apiKey?: APIKey;
		trigger?: Snippet;
	};

	let { open = $bindable(false), apiKey, trigger }: Props = $props();

	const updateExpirationMutation = apiKeyService.updateExpiration();
	const form = superForm(defaults(zod4(updateApiKeyExpirationFormSchema)), {
		SPA: true,
		dataType: "json",
		validators: zod4Client(updateApiKeyExpirationFormSchema),
		async onUpdate({ form: formState }) {
			if (!formState.valid) return;
			if (!apiKey) return;
			updateExpirationMutation.mutate(
				{
					id: apiKey.id,
					data: {
						expires_at: fromExpirationFormValue(formState.data.expires_at),
					},
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
		if (apiKey) {
			form.form.set({
				expires_at: toExpirationFormValue(apiKey.expires_at),
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

	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Update expiration</Dialog.Title>
			<Dialog.Description>
				{#if apiKey}
					Update the expiration date for “{apiKey.name}”.
				{:else}
					Update the API key expiration date.
				{/if}
			</Dialog.Description>
		</Dialog.Header>

		<form use:enhance>
			<FieldGroup>
				<FormExpiration
					{form}
					field="expires_at"
					label="Expires at"
					options={DEFAULT_EXPIRATION_OPTIONS}
					allowCustom
					description="Choose a preset or pick a custom date."
				/>
				<FormError {form} />
				<Dialog.Footer>
					<Button variant="outline" type="button" onclick={handleCancel}>Cancel</Button>
					<Button disabled={updateExpirationMutation.isPending} type="submit">Update</Button>
				</Dialog.Footer>
			</FieldGroup>
		</form>
	</Dialog.Content>
</Dialog.Root>
